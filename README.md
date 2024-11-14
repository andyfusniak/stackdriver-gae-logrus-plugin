# Stackdriver GAE Logrus Plugin

[Logrus is a structured logger for Go](https://github.com/sirupsen/logrus) compatible with the standard library logger.

Stackdriver GAE Logrus Plugin is a lightweight plugin providing threaded message entries on Google Cloud Run and Google App Engine (GAE) Standard Environment.

Cloud Run and GAE adds a request header `X-Cloud-Trace-Context` to HTTP requests. `middleware.XCloudTraceContext` provides an HTTP Middleware component that parses this header and make it available in the context.


``` go
contextLogger := log.WithContext(r.Context())

contextLogger.WithFields(log.Fields{
	"battery": "50",
}).Debug("Flux capacitor low")

contextLogger.WithFields(log.Fields{
	"status": "busted",
}).Info("Warp speed activated")

contextLogger.WithFields(log.Fields{
	"status": "hmmm",
}).Warn("You have been warning")

contextLogger.WithFields(log.Fields{
	"status": "busted",
}).Error("These are not the drones you are looking for")
```

In the Stack Driver Logging Cloud Console you can filter results by request.

![Filter results by request](./docs/screenshots/filter-show-requests-only.png)

Opening an individual request will show each log entry nested within. Each entry will have a color coded icon showing its severity corrsponding to its nearest equivilent Logrus level.

![](./docs/screenshots/threaded-log-entries.png)


Whilst debugging select any individual log message's `trace` value and click "Show matching entries".

![](./docs/screenshots/show-matching-entries.png)


This will automatically set the filters to reveal all log entries associated to that HTTP request.

![](./docs/screenshots/log_entries_list.png)


Clicking the discover icon expands the log entry to reveal the `jsonPayload` section. `jsonPayload.data` contains the field data set using the standard Logrus `WithFields` method.

``` go
contextLogger.WithFields(log.Fields{
	"status": "busted",
}).Error("These are not the drones you are looking for")
```

![](./docs/screenshots/json_payload.png)


To run the example locally use:
``` bash
make run
```

## Windows

If you do not have Make installed.

``` bash
go run ./examples/http-service/main.go
```

## Limitations

The Trace ID and Span ID can be attained using either the `Traceparent` header (v2), or the `X-Cloud-Trace-Context` header (v1). Google App Engine and Google Cloud Run use hexadecimal trace and span IDs and my understanding is that they intend to do so in the future. As of v0.3.0 the formatter will assume that the trace and span IDs are hexadecimal in the request headers sent by Google App Engine and Google Cloud Run and will format them as such.


## Awcknowledgements

Thanks to [csilvers](https://github.com/csilvers) for his contributions that proces the Span ID, as well as his help with understanding how logging works in Google App Engine.

## Todo

- Implement a solution based on the `Traceparent` header and only fallback to `X-Cloud-Trace-Context` if the `Traceparent` header is not present.
