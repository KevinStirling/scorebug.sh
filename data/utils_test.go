package data_test

import (
	"testing"

	"github.com/KevinStirling/scorebug.sh/data"
)

func TestRunnerOn(t *testing.T) {
	cases := []struct {
		name    string
		on      bool
		desired string
	}{
		{"runner on base", true, "◆"},
		{"no runner", false, "◇"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := data.RunnerOn(tc.on)
			if got != tc.desired {
				t.Errorf("RunnerOn(%v) = %s, want %s", tc.on, got, tc.desired)
			}
		})
	}
}

func TestSetOut(t *testing.T) {
	cases := []struct {
		name    string
		outs    int
		pos     int
		desired string
	}{
		{"0 outs, pos 1", 0, 1, "◯"},
		{"1 out, pos 1", 1, 1, "◉"},
		{"1 out, pos 2", 1, 2, "◯"},
		{"2 outs, pos 2", 2, 2, "◉"},
		{"3 outs, pos 3", 3, 3, "◉"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := data.SetOut(tc.outs, tc.pos)
			if got != tc.desired {
				t.Errorf("SetOut(%d, %d) = %s, want %s", tc.outs, tc.pos, got, tc.desired)
			}
		})
	}
}
