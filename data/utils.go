package data

import "strings"

const (
	RUNNER_TRUE  = "◆"
	RUNNER_FALSE = "◇"

	OUT_TRUE  = "◉"
	OUT_FALSE = "◯"

	INNING_TOP    = "↑"
	INNING_BOTTOM = "↓"
)

func runnerOn(on bool) string {
	if on {
		return RUNNER_TRUE
	}
	return RUNNER_FALSE
}

func setOut(outs, pos int) string {
	if pos <= outs {
		return OUT_TRUE
	}
	return OUT_FALSE
}

func setInningArrow(state, position string) string {
	if state == position {
		switch state {
		case "Top":
			return INNING_TOP
		case "Bottom":
			return INNING_BOTTOM
		default:
			return ""
		}
	}
	return ""
}

func RenderInningState(state string) string {
	state = strings.ToLower(state)
	switch state {
	case "top":
		return INNING_TOP
	case "bottom":
		return INNING_BOTTOM
	default:
		return "-"
	}
}
