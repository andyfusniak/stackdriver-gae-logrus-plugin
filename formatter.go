package stackdriver

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/andyfusniak/stackdriver-gae-logrus-plugin/middleware"

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

// The SPAN_ID is a decimal number in v1 of the trace-context API,
// but a (16-digit) hexadecimal number in v2:
//    https://cloud.google.com/python/docs/reference/cloudtrace/1.2.0/google.cloud.trace_v1.types.TraceSpan
//    https://cloud.google.com/python/docs/reference/cloudtrace/1.2.0/google.cloud.trace_v2.types.Span
// Stackdriver uses only v2 these days:
//    https://issuetracker.google.com/issues/338634230?pli=1
// Unfortunately, the X-Cloud-Trace-Context may still be using the
// v1 API.  We try to detect that situation and convert it to a v2
// compatible trace-id by formatting the number in hex.
//
// We cannot perfectly detect whether we're v1 or v2, because
// there are some decimal numbers that look like 16-digit hex
// numbers.  But not very many; we'd expect the logic below to
// make a mistake about 1 in every 2,000 log entries.  In such
// cases, we will use the wrong span-id, making it harder to
// associate log-messages with the requestlog that spawned it.
func maybeFixSpanID(s string) string {
	// We say it's a v2 span-id if it's the right length and is
	// obviously a hex value.  Otherwise we assume it's a v1 span id.
	// About 0.0005 of v2 span-ids will fail this check (because the
	// only have digits 0-9 in them) and seem like v1, and thus be
	// wrongly converted.
	if len(s) == 16 && strings.IndexAny(s, "abcdefABCDEF") != -1 {
		return s
	}

	sAsInt, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return s
	}
	// While the docs don't say it, in practice the hex numbers
	// have lowercase a-f, so that's what we do.
	return strconv.FormatInt(sAsInt, 16)
}

// "X-Cloud-Trace-Context: TRACE_ID/SPAN_ID;o=TRACE_TRUE"
//
// `TRACE_ID` is a 32-character hexadecimal value representing a 128-bit
// number. It should be unique between your requests, unless you
// intentionally want to bundle the requests together. You can use UUIDs.
//
// `SPAN_ID` is the decimal representation of the (unsigned) span ID. It
// should be 0 for the first span in your trace. For subsequent requests,
// set SPAN_ID to the span ID of the parent request. See the description
// of TraceSpan (REST, RPC) for more information about nested traces.
//
// `TRACE_TRUE` must be 1 to trace this request. Specify 0 to not trace the
// request.
func parseXCloudTraceContext(t string) (traceID, spanID string) {
	if t == "" {
		return "", ""
	}

	// 32 characters plus 1 (forward slash) plus 1 (at least one decimal
	// representing the span).
	if len(t) < 34 {
		return "", ""
	}

	// The first character after the TRACE_ID should be a forward slash.
	if t[32] != '/' {
		return "", ""
	}

	// handle "TRACE_ID/SPAN_ID" missing the ";o=1" part.
	last := strings.LastIndex(t, ";")
	if last == -1 {
		return t[0:32], maybeFixSpanID(t[33:])
	}
	return t[0:32], maybeFixSpanID(t[33:last])
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
		traceID, spanID := parseXCloudTraceContext(string(xctc))
		if traceID != "" && spanID != "" {
			ee.Trace = fmt.Sprintf("projects/%s/traces/%s", f.projectID, traceID)
			ee.SpanID = spanID
		}
	}

	b, err := json.Marshal(ee)
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}
