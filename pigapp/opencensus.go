package main

import (
	"fmt"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/trace"
	"golang.org/x/net/context"
)

var sd *stackdriver.Exporter

func init() {
	var err error
	sd, err = stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
		// MetricPrefix helps uniquely identify your metrics.
		MetricPrefix: "demo-prefix",
	})
	if err != nil {
		log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
	}

	// Configure 100% sample rate, otherwise, few traces will be sampled.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Register it as a trace exporter
	trace.RegisterExporter(sd)
}

func startSpan(c context.Context, caption string) (c2 context.Context, endSpan func()) {
	infof(c, "start of span %q", caption)
	var span *trace.Span
	c2, span = trace.StartSpan(c, caption)
	// c2, span = trace.StartSpan(context.Background(), caption)
	endSpan = func() {
		span.End()
		infof(c, "end of span %q", caption)
		sd.Flush()
	}
	return c2, endSpan
}

func startSpanf(c context.Context, msg string, args ...interface{}) (c2 context.Context, endSpan func()) {
	caption := fmt.Sprintf(msg, args...)
	return startSpan(c, caption)
}

// "With Remote Parent" ...??
func startSpanfWRT(r *http.Request, msg string, args ...interface{}) (c2 context.Context, endSpan func()) {
	caption := fmt.Sprintf(msg, args...)

	c := r.Context()
	infof(c, "start of spanWRT %q", caption)
	spanContext, ok := (&propagation.HTTPFormat{}).SpanContextFromRequest(r)
	if !ok {
		return c, func() {}
	}
	var span *trace.Span
	c2, span = trace.StartSpanWithRemoteParent(c, caption, spanContext)
	// c2, span = trace.StartSpan(context.Background(), caption)
	endSpan = func() {
		span.End()
		infof(c, "end of spanWRT %q", caption)
		// sd.Flush()  // Should not be necessary *every time* ...?
	}
	return c2, endSpan
}
