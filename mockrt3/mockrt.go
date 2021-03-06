/*
Package mockrt3 provides mock runtime for mockgo
*/
package mockrt3

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// P is a trait for types of request parameter.
type P interface{ P__() }

// R is a trait for types of response result.
type R interface{ R__() }

// C defines pair of request parameter (P) and response result (R) for a method
// call.
type C struct {
	P P
	R R
}

// Q is a checker of sequence of method calls
type Q struct {
	t     *testing.T
	calls []C
	opts  []cmp.Option
	index int
}

// NewQ is an alias for NewSequence, creates a sequence of calls.
// This is called by test codes.
func NewQ(t *testing.T, calls ...C) *Q {
	return &Q{
		t:     t,
		calls: calls,
	}
}

// AddCall adds call data.
// This is called by test codes.
func (q *Q) AddCall(calls ...C) *Q {
	q.calls = append(q.calls, calls...)
	return q
}

// WithOption updates compare option.
// This is called by test codes.
func (q *Q) WithOption(opts ...cmp.Option) *Q {
	q.opts = opts
	return q
}

// Call checks call parameter and returns result.
// This is called by mock code.
func (q *Q) Call(name string, param P) R {
	q.t.Helper()
	if q.index >= len(q.calls) {
		q.t.Fatalf("no calls at #%d for %s\nparam=%+v", q.index, name, param)
	}
	c := q.calls[q.index]
	if d := cmp.Diff(c.P, param, q.opts...); d != "" {
		q.t.Fatalf("call for %s (#%d) has unexpected arguments: -want +got\n%s", name, q.index, d)
	}
	q.index++
	return c.R
}

// T returns *testing.T.
// This is called by mock code.
func (q *Q) T() *testing.T {
	return q.t
}

// IsEnd checks sequence has end or not.
// This is called by test code.
func (q *Q) IsEnd() {
	q.t.Helper()
	if q.index < len(q.calls) {
		q.t.Fatalf("there are non-proceeded calles: %+v", q.calls[q.index:])
	}
}
