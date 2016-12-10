/*
http://www.apache.org/licenses/LICENSE-2.0.txt

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
	"fmt"
	"time"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	"github.com/IrekRomaniuk/snap-plugin-collector-portscan/portscan/targets"
	//"github.com/intelsdi-x/snap-plugin-utilities/config"
	"net"
)
const (
	vendor        = "niuk"
	fs            = "portscan"
	pluginName    = "portscan"
	pluginVersion = 1
	pluginType    = plugin.CollectorPluginType
)
var (
	metricNames = []string{
		//"53",
	}
)
type PortscanCollector struct {
}

func New() *PortscanCollector {
	portscan := &PortscanCollector{}
	return portscan
}

/*  CollectMetrics collects metrics for testing.

CollectMetrics() will be called by Snap when a task that collects one of the metrics returned from this plugins
GetMetricTypes() is started. The input will include a slice of all the metric types being collected.

The output is the collected metrics as plugin.Metric and an error.
*/
func (portscan *PortscanCollector) CollectMetrics(mts []plugin.MetricType) (metrics []plugin.MetricType, err error) {
	var (
		//err error
		target string
		port string
		timeout = time.Duration(1) * time.Second
	)
	conf := mts[0].Config().Table()
        //fmt.Println(conf)
	targetConf, ok := conf["target"]
	if !ok || targetConf.(ctypes.ConfigValueStr).Value == "" {
		return nil, fmt.Errorf("target missing from config, %v", conf)
	} else {
		target = targetConf.(ctypes.ConfigValueStr).Value
	}
	portConf, ok := conf["port"]
	if !ok || portConf.(ctypes.ConfigValueStr).Value == "" {
		return nil, fmt.Errorf("port missing from config, %v", conf)
	} else {
		port = portConf.(ctypes.ConfigValueStr).Value
	}

	hosts, err := targets.ReadTargets(target)
	if err != nil { return nil, fmt.Errorf("Error reading target: %v", err) }
	if len(hosts) == 0 { return nil, fmt.Errorf("No host defined in file %v", target)}

	count, _ := scan(hosts, port, timeout)

	metric := plugin.MetricType{
		Namespace_: core.NewNamespace("niuk", "portscan", port), //ns
		//Namespace_: core.NewNamespace("niuk", "portscan").AddDynamicElement("port_number", port),
		Data_:      count,
		Timestamp_: time.Now(),
	}
	//fmt.Print(metric.Namespace())
	metrics = append(metrics, metric)

	return metrics, nil
}

func scan(hosts []string, port string, timeout time.Duration) (int, error) {
	ports := make(map[string]int)
	for _, host := range hosts {
		//fmt.Printf("%s\n", host + ":" + port)
		conn, err := net.DialTimeout("tcp", host + ":" + port, timeout)
		if err == nil {
			conn.Close()
			ports[port]++
		} else {
			//fmt.Errorf("No connection: %v", err)
		}
	}
	return ports[port], nil
}

/*
	GetMetricTypes returns metric types for testing.
	GetMetricTypes() will be called when your plugin is loaded in order to populate the metric catalog(where snaps stores all
	available metrics).

	Config info is passed in. This config information would come from global config snap settings.

	The metrics returned will be advertised to users who list all the metrics and will become targetable by tasks.
*/
func (portscan *PortscanCollector) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}

	//for _, metricName := range metricNames {
		mts = append(mts, plugin.MetricType{
			//Namespace_: core.NewNamespace("niuk", "portscan", metricName),
			//Namespace_: createNamespace(metricName.ns),
			Namespace_: core.NewNamespace("niuk", "portscan").AddDynamicElement("Name","Description").
				AddStaticElement("blabla"),//?!
			//Description_: "Name_Description: " ,
		})
	//}
	return mts, nil
}

// GetConfigPolicy returns plugin configuration
func (portscan *PortscanCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule0, _ := cpolicy.NewStringRule("target", true)
	rule1, _ := cpolicy.NewStringRule("port", true)
	cp := cpolicy.NewPolicyNode()
	cp.Add(rule0)
	cp.Add(rule1)
	c.Add([]string{"niuk", "portscan"}, cp)
	return c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		pluginName,
		pluginVersion,
		pluginType,
		[]string{plugin.SnapGOBContentType},//[]string{},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.ConcurrencyCount(1),
	)
}

