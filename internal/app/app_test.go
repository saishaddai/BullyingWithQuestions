package app

import "testing"

func TestRunReturnsNilAtBaseline(t *testing.T) {
	if err := Run(); err != nil {
		t.Fatalf("Run() returned an unexpected error: %v", err)
	}
}
