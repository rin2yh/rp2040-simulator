// Package tabletest provides the shared control flow for table-driven tests.
package tabletest

import "testing"

// Case describes one table-driven test case.
type Case[I, W any] struct {
	Name string
	In   I
	Want W
}

// Run executes each case as a named subtest.
func Run[I, W any](t *testing.T, cases []Case[I, W], run func(t *testing.T, in I, want W)) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			run(t, tc.In, tc.Want)
		})
	}
}
