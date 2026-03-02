package data

func RunnerOn(on bool) string {
	if on {
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

func setInningArrow(state, position string) string {
	if state == position {
		switch state {
		case "Top":
			return "↑"
		case "Bottom":
			return "↓"
		default:
			return ""
		}
	}
	return ""
}
