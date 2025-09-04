package data

func SetRunnerState(s []int, base int) string {
	for i := range s {
		if s[i] == base {
			return "●"
		}
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
