package stackdriver

import (
	"encoding/json"
	"fmt"
	"github.com/andyfusniak/stackdriver-gae-logrus-plugin/middleware"
	"strings"

	"github.com/sirupsen/logrus"
)

// The severity of the event described in a log entry.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
const (
	lsDEFAULT   = "DEFAULT"
	lsDEBUG     = "DEBUG"
	lsINFO      = "INFO"
	lsNOTICE    = "NOTICE"
	lsWARNING   = "WARNING"
	lsERROR     = "ERROR"
	lsCRITICAL  = "CRITICAL"
	lsALERT     = "ALERT"
	lsEMERGENCY = "EMERGENCY"
)

var levToSev = map[logrus.Level]string{
	logrus.TraceLevel: lsDEFAULT,
	logrus.DebugLevel: lsDEBUG,
	logrus.InfoLevel:  lsINFO,
	logrus.WarnLevel:  lsWARNING,
	logrus.ErrorLevel: lsERROR,
	logrus.FatalLevel: lsCRITICAL,
	logrus.PanicLevel: lsEMERGENCY,
}

// Formatter implements threaded Stackdriver formatting for logrus.
type Formatter struct {
	projectID string
}

type entry struct {
	// Trace string
	// Optional. Resource name of the trace associated with the log
	// entry, if any. If it contains a relative resource name, the name
	// is assumed to be relative to //tracing.googleapis.com. Example:
	// projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824
	Trace string `json:"logging.googleapis.com/trace,omitempty"`

	// Span
	// Optional. The span ID within the trace associated with the log
	// entry.
	//
	// For Trace spans, this is the same format that the Trace API
	// v2 uses: a 16-character hexadecimal encoding of an 8-byte
	// array, such as "000000000000004a"
	SpanID string `json:"logging.googleapis.com/spanId,omitempty"`

	Data     logrus.Fields `json:"data"`
	Message  string        `json:"message,omitempty"`
	Severity string        `json:"severity,omitempty"`
}

// GAEStandardFormatter returns a new Formatter.
func GAEStandardFormatter(options ...Option) *Formatter {
	fmtr := Formatter{}
	for _, option := range options {
		option(&fmtr)
	}
	return &fmtr
}

func parseXCloudTraceContext(t string) (traceID, spanID, traceTrue string) {
	slash := strings.Index(t, "/")
	semi := strings.Index(t, ";")
	equal := strings.Index(t, "=")
	return t[0:slash], t[slash+1 : semi], t[equal+1:]
}

// Option lets you configure the Formatter.
type Option func(*Formatter)

// WithProjectID lets you configure the GAE project for threaded messaging.
func WithProjectID(pid string) Option {
	return func(f *Formatter) {
		f.projectID = pid
	}
}

// XCTContext holds the X Cloud Trace Context value
type XCTContext string

// Format formats a logrus entry in Stackdriver format.
func (f *Formatter) Format(e *logrus.Entry) ([]byte, error) {
	ee := entry{
		Severity: levToSev[e.Level],
		Message:  e.Message,
		Data:     e.Data,
	}

	xctc := middleware.XCTC(e.Context)
	if xctc != "" {
		traceID, spanID, _ := parseXCloudTraceContext(string(xctc))
		ee.Trace = fmt.Sprintf("projects/%s/traces/%s", f.projectID, traceID)
		ee.SpanID = spanID
	}

	b, err := json.Marshal(ee)
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}
