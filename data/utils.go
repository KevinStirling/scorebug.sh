package data

func setRunnerState(s []int, i int) string {
	if i >= 0 && i < len(s) && s[i] != 0 {
		return "●"
	}
	return " "
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
