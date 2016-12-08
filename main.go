// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/circonus-labs/snap-plugin-publisher-circonus/circonus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

func main() {
	plugin.StartPublisher(&circonus.Publisher{}, circonus.Name, circonus.Version, plugin.Exclusive(true))
}
