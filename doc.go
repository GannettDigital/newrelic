/*
package newrelic is an client for NewRelic's plugin API.

Simple usage:
	import (
		"runtime"
		"time"

		"github.com/neocortical/newrelic"
	)

	func main() {
		client := newrelic.New("abc123") // license key goes here
		myplugin := &Plugin{
			Name: "My Plugin",
			GUID: "com.example.newrelic.myplugin",
		}
		client.AddPlugin(myplugin)

		// the easiest way to add a metric is to use a closure
		metric := newrelic.NewMetric("MyApp/Total CGO Calls",
			func() (float64, error) { return float64(runtime.NumCgoCall()), nil })
		plugin.AddMetric()

		// call run after doing all plugin config
		client.Run()

		// main application code here..
		for {
			time.Sleep(time.Minute)
		}
	}
*/
package newrelic
