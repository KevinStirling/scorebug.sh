package data

import (
	"slices"

	"github.com/charmbracelet/log"
)

func SetRunnerState(s []int, base int) string {
	log.Debugf("RunnerIndex value: %v, searching for %d", s, base)
	if slices.Contains(s, base) {
		return "◆"
	}
	return "◇"
}

func SetOut(outs, pos int) string {
	if pos <= outs {
		return "◉"
	}
	return "◯"
}

func setInningArrow(state string) string {
	switch state {
	case "Top":
		return "↑"
	case "Bottom":
		return "↓"
	default:
		return ""
	}
}
