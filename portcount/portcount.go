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

package portcount

import (
	"fmt"
	"time"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	"github.com/IrekRomaniuk/snap-plugin-collector-portcount/portcount/targets"
	//"github.com/intelsdi-x/snap-plugin-utilities/config"
	"net"
	"sync"
)
const (
	vendor        = "niuk"
	fs            = "portcount"
	pluginName    = "portcount"
	pluginVersion = 1
	pluginType    = plugin.CollectorPluginType
)
type PortcountCollector struct {
}

func New() *PortcountCollector {
	portcount := &PortcountCollector{}
	return portcount
}

/*  CollectMetrics collects metrics for testing.

CollectMetrics() will be called by Snap when a task that collects one of the metrics returned from this plugins
GetMetricTypes() is started. The input will include a slice of all the metric types being collected.

The output is the collected metrics as plugin.Metric and an error.
*/
func (portcount *PortcountCollector) CollectMetrics(mts []plugin.MetricType) (metrics []plugin.MetricType, err error) {
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
		Namespace_: core.NewNamespace("niuk", "portcount", port), //ns
		Data_:      count,
		Timestamp_: time.Now(),
	}
	//fmt.Print(metric.Namespace())
	metrics = append(metrics, metric)

	return metrics, nil
}

func scan(hosts []string, port string, timeout time.Duration) (int, error) {
	d := net.Dialer{Timeout: timeout}
	p := make(chan struct{}, 500) // make 500 parallel connection
	wg := sync.WaitGroup{}
        var result int

	c := func(host string) {
		conn, err := d.Dial(`tcp`, fmt.Sprintf(`%s:%s`, host, port))
		if err == nil {
			conn.Close()
			//fmt.Printf("%d passed\n", port)
			result ++
		}
		<-p
		wg.Done()
	}

	wg.Add(len(hosts))
	for _, host := range hosts {
		p <- struct{}{}
		go c(host)
	}
	wg.Wait()


       return result, nil
}

/*
	GetMetricTypes returns metric types for testing.
	GetMetricTypes() will be called when your plugin is loaded in order to populate the metric catalog(where snaps stores all
	available metrics).

	Config info is passed in. This config information would come from global config snap settings.

	The metrics returned will be advertised to users who list all the metrics and will become targetable by tasks.
*/
func (portcount *PortcountCollector) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}

	//for _, metricName := range metricNames {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace("niuk", "portcount").AddDynamicElement("Port","Port to scan").
				AddStaticElement("total_up"),//?!
			//Description_: "Name_Description: " ,
		})
	//}
	return mts, nil
}

// GetConfigPolicy returns plugin configuration
func (portcount *PortcountCollector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule0, _ := cpolicy.NewStringRule("target", true)
	rule1, _ := cpolicy.NewStringRule("port", true)
	cp := cpolicy.NewPolicyNode()
	cp.Add(rule0)
	cp.Add(rule1)
	c.Add([]string{"niuk", "portcount"}, cp)
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

