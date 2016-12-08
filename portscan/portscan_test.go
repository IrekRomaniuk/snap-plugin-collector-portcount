/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package portscan

import (
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/IrekRomaniuk/snap-plugin-collector-portscan/portscan/targets"
	"time"
)

func TestPortscanPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, pluginName )
		So(meta.Version, ShouldResemble, pluginVersion)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})
	Convey("Create Portscan Collector", t, func() {
		collector := New()
		Convey("So Portscan collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})
		Convey("So Portscan collector should be of Portscan type", func() {
			So(collector, ShouldHaveSameTypeAs, &PortscanCollector{})
		})
		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := collector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
				t.Log(configPolicy)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			Convey("So config policy namespace should be /niuk/portscan", func() {
				conf := configPolicy.Get([]string{"niuk", "portscan"})
				So(conf, ShouldNotBeNil)
				So(conf.HasRules(), ShouldBeTrue)
				tables := conf.RulesAsTable()
				So(len(tables), ShouldEqual, 2)
				for _, rule := range tables {
					So(rule.Name, ShouldBeIn, "target", "port")
					switch rule.Name {
					case "target":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					case "port":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					}
				}
			})
		})
	})
}

func TestReadScanTargets(t *testing.T) {
	Convey("Read iplist.txt from examples ", t, func() {
	target := "../examples/iplist.txt"
	hosts, _ := ReadTargets(target)
		Convey("So iplist.txt should contain 4 hosts", func() {
			So(len(hosts), ShouldEqual,4)
		})
		Convey("So 2 hosts should have port 53 opened", func() {
			count, _ := scan(hosts, "53", time.Duration(1) * time.Second)
			So(count, ShouldEqual, 2)
		})
	})
}


func TestPortscanCollector_CollectMetrics(t *testing.T) {
	cfg := setupCfg("../examples/iplist.txt", "53")
	Convey("Portscan collector", t, func() {
		p := New()
		mt, err := p.GetMetricTypes(cfg)
		if err != nil {
			t.Fatal(err)
		}
		So(len(mt), ShouldEqual, 1)
		Convey("collect metrics", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"niuk", "portscan"),
					Config_: cfg.ConfigDataNode,
				},
			}
			metrics, err := p.CollectMetrics(mts)
			So(err, ShouldBeNil)
			So(metrics, ShouldNotBeNil)
			So(len(metrics), ShouldEqual, 1)
			So(metrics[0].Namespace()[0].Value, ShouldEqual, "niuk")
			So(metrics[0].Namespace()[1].Value, ShouldEqual, "portscan")
			for _, m := range metrics {
				//fmt.Println(m.Namespace()[2].Value,m.Data())
				So(m.Namespace()[2].Value, ShouldEqual, "p53")
				So(m.Data(), ShouldEqual, 2) //Assuming 8.8.8.8:53 and 4.2.2.2:53 respond
				t.Log(m.Namespace()[2].Value, m.Data())
			}
		})
	})
}


func setupCfg(target string, port string) plugin.ConfigType {
	node := cdata.NewNode()
	node.AddItem("target", ctypes.ConfigValueStr{Value: target})
	node.AddItem("port", ctypes.ConfigValueStr{Value: port})
	return plugin.ConfigType{ConfigDataNode: node}
}

