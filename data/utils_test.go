package data_test

import (
	"testing"

	"github.com/KevinStirling/scorebug.sh/data"
)

func TestSetRunnerState(t *testing.T) {
	cases := []struct {
		Name    string
		Runners []int
		Base    int
		Desired string
	}{
		{"RunnerOnFirst", []int{1}, 1, "◆"},
		{"RunnerOnSecond", []int{2}, 2, "◆"},
		{"RunnerOnThird", []int{3}, 3, "◆"},
		{"NoRunner", []int{}, 1, "◇"},
		{"WrongBase", []int{2}, 1, "◇"},
	}

	for i := range cases {
		s := data.SetRunnerState(cases[i].Runners, cases[i].Base)
		if s != cases[i].Desired {
			t.Logf("\"%s\" case expected %s, got %s", cases[i].Name, cases[i].Desired, s)
			t.Fail()
		}
	}
}
