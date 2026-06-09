package data

import (
	"testing"
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
			got := runnerOn(tc.on)
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
			got := setOut(tc.outs, tc.pos)
			if got != tc.desired {
				t.Errorf("SetOut(%d, %d) = %s, want %s", tc.outs, tc.pos, got, tc.desired)
			}
		})
	}
}

func TestSetInningArrow(t *testing.T) {
	cases := []struct {
		name     string
		state    string
		position string
		desired  string
	}{
		{"Top inning", "Top", "Top", INNING_TOP},
		{"No match", "Top", "Bottom", ""},
		{"Bottom inning", "Bottom", "Bottom", INNING_BOTTOM},
		{"Invalid state", "invalid", "Top", ""},
		{"Empty state", "", "Bottom", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := setInningArrow(tc.state, tc.position)
			if got != tc.desired {
				t.Errorf("setInningArrow(%s, %s) = %s, want %s", tc.state, tc.position, got, tc.desired)
			}
		})
	}
}
