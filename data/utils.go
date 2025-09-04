package data

func SetRunnerState(s []int, base int) string {
	if len(s) != 0 {
		for _, v := range s {
			if v == base {
				return "●"
			}
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
