// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonus

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"

	cgm "github.com/circonus-labs/circonus-gometrics"
)

var (
	// Name of plugin
	Name = "circonus"
	// Version of plugin
	Version = 1
)

// Publisher defines the Circonus publisher
type Publisher struct {
	mu      sync.Mutex
	metrics *cgm.CirconusMetrics
}

// GetConfigPolicy returns plugin configuration policy
func (p *Publisher) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	hn, _ := os.Hostname()

	policy.AddNewStringRule([]string{""}, "interval", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "reset_counters", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "reset_gauges", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "reset_histograms", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "reset_text", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "api_token", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "api_app", false, plugin.SetDefaultString("snap-cgm"))
	policy.AddNewStringRule([]string{""}, "api_url", false, plugin.SetDefaultString("https://api.circonus.com/v2"))
	policy.AddNewStringRule([]string{""}, "check_id", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "check_submission_url", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "check_instance_id", false, plugin.SetDefaultString(fmt.Sprintf("%s:snap-telemetry", hn)))
	policy.AddNewStringRule([]string{""}, "check_target_host", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "check_display_name", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "check_search_tag", false, plugin.SetDefaultString("service:snap-telemetry"))
	policy.AddNewStringRule([]string{""}, "check_secret", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "check_tags", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "check_max_url_age", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "check_force_metric_activation", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "broker_id", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "broker_select_tag", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "broker_max_response_time", false, plugin.SetDefaultString(""))
	policy.AddNewStringRule([]string{""}, "log_level", false, plugin.SetDefaultString(""))

	return *policy, nil
}

func getCGMConfig(cfg plugin.Config) (*cgm.Config, error) {
	cmc := &cgm.Config{}

	if val, err := cfg.GetString("interval"); err == nil {
		cmc.Interval = val
	}

	if val, err := cfg.GetString("reset_counters"); err == nil {
		cmc.ResetCounters = val
	}

	if val, err := cfg.GetString("reset_gauges"); err == nil {
		cmc.ResetGauges = val
	}

	if val, err := cfg.GetString("reset_histograms"); err == nil {
		cmc.ResetHistograms = val
	}

	if val, err := cfg.GetString("reset_text"); err == nil {
		cmc.ResetText = val
	}

	if val, err := cfg.GetString("api_token"); err == nil {
		cmc.CheckManager.API.TokenKey = val
	}

	if val, err := cfg.GetString("api_app"); err == nil {
		cmc.CheckManager.API.TokenApp = val
	}

	if val, err := cfg.GetString("api_url"); err == nil {
		cmc.CheckManager.API.URL = val
	}

	if val, err := cfg.GetString("check_id"); err == nil {
		cmc.CheckManager.Check.ID = val
	}

	if val, err := cfg.GetString("check_submission_url"); err == nil {
		cmc.CheckManager.Check.SubmissionURL = val
	}

	if val, err := cfg.GetString("check_instance_id"); err == nil {
		cmc.CheckManager.Check.InstanceID = val
	}

	if val, err := cfg.GetString("check_target_host"); err == nil {
		cmc.CheckManager.Check.TargetHost = val
	}

	if val, err := cfg.GetString("check_display_name"); err == nil {
		cmc.CheckManager.Check.DisplayName = val
	}

	if val, err := cfg.GetString("check_search_tag"); err == nil {
		cmc.CheckManager.Check.SearchTag = val
	}

	if val, err := cfg.GetString("check_secret"); err == nil {
		cmc.CheckManager.Check.Secret = val
	}

	if val, err := cfg.GetString("check_tags"); err == nil {
		cmc.CheckManager.Check.Tags = val
	}

	if val, err := cfg.GetString("check_max_url_age"); err == nil {
		cmc.CheckManager.Check.MaxURLAge = val
	}

	if val, err := cfg.GetString("check_force_metric_activation"); err == nil {
		cmc.CheckManager.Check.ForceMetricActivation = val
	}

	if val, err := cfg.GetString("broker_id"); err == nil {
		cmc.CheckManager.Broker.ID = val
	}

	if val, err := cfg.GetString("broker_select_tag"); err == nil {
		cmc.CheckManager.Broker.SelectTag = val
	}

	if val, err := cfg.GetString("broker_max_response_time"); err == nil {
		cmc.CheckManager.Broker.MaxResponseTime = val
	}

	if cmc.CheckManager.API.TokenKey == "" && cmc.CheckManager.Check.SubmissionURL == "" {
		return nil, fmt.Errorf("One of 'api_token' or 'check_submission_url' is *required* for Circonus publisher plugin")
	}

	return cmc, nil
}

// Publish metrics to Circonus
func (p *Publisher) Publish(metrics []plugin.Metric, cfg plugin.Config) error {
	logger := getLogger(cfg)

	if p.metrics == nil {
		cgmcfg, err := getCGMConfig(cfg)
		if err != nil {
			return err
		}

		cgmcfg.Log = log.New(logger.Logger.Writer(), "", log.LstdFlags)
		cgmcfg.Debug = logrus.GetLevel().String() == "debug"

		p.mu.Lock()
		mh, err := cgm.NewCirconusMetrics(cgmcfg)
		if err != nil {
			p.mu.Unlock()
			return err
		}
		p.metrics = mh
		p.metrics.Start()
		p.mu.Unlock()
	}

	for _, m := range metrics {
		metricName := strings.Join(m.Namespace.Strings(), "`")
		metricType := m.Tags["circonus_type"]
		if metricType == "" {
			metricType = "numeric"
		}
		switch metricType {
		case "numeric":
			// will be ignored if it is not an accpetable (int/uint/float) type
			p.metrics.SetGauge(metricName, m.Data)
		case "text":
			p.metrics.SetText(metricName, toText(m.Data))
		case "histogram":
			val, err := toFloat64(m.Data)
			if err != nil {
				logger.Errorf("Unable to convert %v to uint64 for histogram %v", m.Data, err)
			} else {
				p.metrics.RecordValue(metricName, val)
			}
		default:
			logger.Errorf("Unsupported circonus metric type %s", metricType)
		}
	}

	return nil
}

func toText(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func toFloat64(i interface{}) (float64, error) {

	switch val := i.(type) {
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	case float64:
		return float64(val), nil
	default:
		return 0, fmt.Errorf("Unable to convert [%v] to float64, unsupported type", val)
	}
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

func getLogger(cfg plugin.Config) *logrus.Entry {
	logger := logrus.WithFields(logrus.Fields{
		"plugin-name":    Name,
		"plugin-version": Version,
		"plugin-type":    "publisher",
	})

	logrus.SetLevel(logrus.WarnLevel)

	levelValue, err := cfg.GetString("log_level")
	if err == nil {
		if level, err := logrus.ParseLevel(strings.ToLower(levelValue)); err == nil {
			logrus.SetLevel(level)
		} else {
			logrus.WithFields(logrus.Fields{
				"value":             strings.ToLower(levelValue),
				"acceptable values": "warn, error, debug, info",
			}).Warn("Invalid config value")
		}
	}
	return logger
}
