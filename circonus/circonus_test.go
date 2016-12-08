package circonus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	. "github.com/smartystreets/goconvey/convey"
)

func testServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			fallthrough
		case "PUT":
			defer r.Body.Close()
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			var ret []byte
			var r interface{}
			err = json.Unmarshal(b, &r)
			if err != nil {
				ret, err = json.Marshal(err)
				if err != nil {
					panic(err)
				}
			} else {
				ret, err = json.Marshal(r)
				if err != nil {
					panic(err)
				}
			}
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(ret))
		default:
			w.WriteHeader(500)
			fmt.Fprintln(w, "unsupported method")
		}
	}
	return httptest.NewServer(http.HandlerFunc(f))
}

func TestCirconusPlugin(t *testing.T) {
	Convey("Create Circonus publisher", t, func() {
		cp := &Publisher{}
		Convey("So publisher should not be nil", func() {
			So(cp, ShouldNotBeNil)
		})

		Convey("Publisher should be of type Publisher", func() {
			So(cp, ShouldHaveSameTypeAs, &Publisher{})
		})

		configPolicy, err := cp.GetConfigPolicy()
		Convey("Should return a config policy", func() {
			Convey("configPolicy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)

				Convey("and retrieving config policy should not error", func() {
					So(err, ShouldBeNil)

					Convey("config policy should be a cpolicy.ConfigPolicy", func() {
						So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
					})

					testConfig := make(plugin.Config)
					testConfig["api_token"] = "foo"
					testConfig["check_submission_url"] = "http://127.0.0.1:43191/"

					token, err := testConfig.GetString("api_token")
					Convey("So testConfig should return the right api_token config", func() {
						So(err, ShouldBeNil)
						So(token, ShouldEqual, "foo")
					})

					url, err := testConfig.GetString("check_submission_url")
					Convey("So testConfig should return the right check_submission_url config", func() {
						So(err, ShouldBeNil)
						So(url, ShouldEqual, "http://127.0.0.1:43191/")
					})
				})
			})
		})
	})

	Convey("toText should return a string", t, func() {
		metric := plugin.Metric{
			Namespace: plugin.NewNamespace("test"),
			Data:      int32(1),
			Tags:      map[string]string{"circonus_type": "numeric"},
		}

		str := toText(metric.Data)
		So(str, ShouldNotBeNil)
	})

	Convey("toFloat64 should return a float", t, func() {

		var (
			f   float64
			err error
		)

		metric := plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: int(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: int8(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: int16(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: int32(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: int64(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: uint(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: uint8(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: uint16(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: uint32(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: uint64(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: float32(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: float64(1)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldBeNil)
		So(f, ShouldBeGreaterThan, 0)

		metric = plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: bool(true)}
		f, err = toFloat64(metric.Data)
		So(err, ShouldNotBeNil)
		So(f, ShouldBeZeroValue)
	})
}

func TestCirconusPublish(t *testing.T) {
	server := testServer()
	defer server.Close()

	config := make(plugin.Config)
	config["api_token"] = ""
	config["api_app"] = ""
	config["api_url"] = ""
	config["check_id"] = ""
	config["check_submission_url"] = ""
	config["log_level"] = "debug"
	cp := &Publisher{}

	Convey("snap publishing plugin Circonus testing", t, func() {

		Convey("Invalid config", func() {
			metrics := []plugin.Metric{
				plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: bool(true)},
			}
			err := cp.Publish(metrics, config)
			So(err, ShouldNotBeNil)
		})

		config["check_submission_url"] = server.URL + "/metrics_endpoint"

		Convey("Publish numeric metric", func() {
			metrics := []plugin.Metric{
				plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: int(1), Tags: map[string]string{"circonus_type": "numeric"}},
			}
			err := cp.Publish(metrics, config)
			So(err, ShouldBeNil)
		})

		Convey("Publish histogram metric", func() {
			metrics := []plugin.Metric{
				plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: int(1), Tags: map[string]string{"circonus_type": "histogram"}},
			}
			err := cp.Publish(metrics, config)
			So(err, ShouldBeNil)
		})

		Convey("Publish histogram metric [bad type]", func() {
			metrics := []plugin.Metric{
				plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: string("foo"), Tags: map[string]string{"circonus_type": "histogram"}},
			}
			err := cp.Publish(metrics, config)
			So(err, ShouldBeNil)
		})

		Convey("Publish text metric", func() {
			metrics := []plugin.Metric{
				plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: string("foo"), Tags: map[string]string{"circonus_type": "text"}},
			}
			err := cp.Publish(metrics, config)
			So(err, ShouldBeNil)
		})

		Convey("Publish bad metric type", func() {
			metrics := []plugin.Metric{
				plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: int(1), Tags: map[string]string{"circonus_type": "badtype"}},
			}
			err := cp.Publish(metrics, config)
			So(err, ShouldBeNil)
		})

		Convey("Publish default metric type", func() {
			metrics := []plugin.Metric{
				plugin.Metric{Namespace: plugin.NewNamespace("test"), Data: int(1)},
			}
			err := cp.Publish(metrics, config)
			So(err, ShouldBeNil)
		})
	})
}
