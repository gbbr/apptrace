package httptrace

import (
	"net/http"

	"sourcegraph.com/sourcegraph/apptrace"
)

const (
	// HeaderSpanID is the name of the HTTP header by which the trace
	// and span IDs are passed along.
	HeaderSpanID = "Span-ID"

	// HeaderParentSpanID is the name of the HTTP header by which the
	// parent trace and span IDs are passed along. It should only be
	// set by clients that are incapable of creating their own span
	// IDs (e.g., JavaScript API clients in a web page, which can
	// easily pass along an existing parent span ID but not create a
	// new child span ID).
	HeaderParentSpanID = "Parent-Span-ID"
)

// SetSpanIDHeader sets the Span-ID header.
func SetSpanIDHeader(h http.Header, e apptrace.SpanID) {
	h.Set(HeaderSpanID, e.String())
}

// GetSpanID returns the SpanID for the current request, based on the
// values in the HTTP headers. If a Span-ID header is provided, it is
// parsed; if a Parent-Span-ID header is provided, a new child span is
// created and it is returned; otherwise a new root SpanID is created.
func GetSpanID(h http.Header) (*apptrace.SpanID, error) {
	spanID, _, err := getSpanID(h)
	return spanID, err
}

func getSpanID(h http.Header) (spanID *apptrace.SpanID, fromHeader string, err error) {
	// Check for Span-ID.
	fromHeader = HeaderSpanID
	spanID, err = getSpanIDHeader(h, HeaderSpanID)
	if err != nil {
		return nil, fromHeader, err
	}

	// Check for Parent-Span-ID.
	if spanID == nil {
		fromHeader = HeaderParentSpanID
		spanID, err = getSpanIDHeader(h, HeaderParentSpanID)
		if err != nil {
			return nil, fromHeader, err
		}
	}

	// Create a new root span ID.
	if spanID == nil {
		fromHeader = ""
		newSpanID := apptrace.NewRootSpanID()
		spanID = &newSpanID
	}
	return spanID, fromHeader, nil
}

// getSpanIDHeader returns the SpanID in the header (specified by
// key), nil if no such header was provided, or an error if the value
// was unparseable.
func getSpanIDHeader(h http.Header, key string) (*apptrace.SpanID, error) {
	s := h.Get(key)
	if s == "" {
		return nil, nil
	}
	return apptrace.ParseSpanID(s)
}
