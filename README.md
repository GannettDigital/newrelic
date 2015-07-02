# newrelic: Go NewRelic plugin API client

newrelic is an embeddable client for implementing custom NewRelic plugins. It is aided and inpired by  [yvasiyarov's](https://github.com/yvasiyarov) [newrelic_platform_go package](https://github.com/yvasiyarov/newrelic_platform_go).

[![Build Status](https://travis-ci.org/neocortical/newrelic.svg?branch=master)](https://travis-ci.org/neocortical/newrelic) [![Coverage](http://gocover.io/_badge/github.com/neocortical/newrelic)](http://gocover.io/github.com/neocortical/newrelic) [![GoDoc](https://godoc.org/github.com/neocortical/newrelic?status.svg)](https://godoc.org/github.com/neocortical/newrelic)

# Installation

```go
go get github.com/neocortical/newrelic
```

# Use

```go
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
```

# Advanced Features

### Set log levels and custom log destination
```go
newrelic.LogLevel = newrelic.LogAll
newrelic.Logger = myAwesomeLogger // standard library logger
```

### Use an HTTP proxy to send data to NewRelic

```go
proxy := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(myProxyUrl)}}
client := newrelic.New("abc123")
client.HTTPClient = proxy

```

# Implementation Notes

The NewRelic plugin API reference can be found [here](https://docs.newrelic.com/docs/plugins/plugin-developer-resources/planning-your-plugin/parts-plugin). There is some naming confusion in the API that can throw people off. Namely, when crafting API requests, the term `components` is used when `plugins` would be more accurate. Additionally, in the reference, the term Agent refers to both the code interacting with the API and the host/process information sent in requests.

When writing this package, the term Plugin is used exclusively to refer to plugins, and the term Client is used to refer to the "agent" that sends data to NewRelic. The objects that model JSON requests are completely separated from business objects to keep things clear.
