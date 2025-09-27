package data

import (
	"slices"
)

func SetRunnerState(s []int, base int) string {
	if slices.Contains(s, base) {
		return "◆"
	}
	return "◇"
}

func mark[T any](ptr *T) string {
	if ptr != nil {
		return "◆"
	}
	return "◇"
}

func SetBaseRunner(feed Feed) (on1, on2, on3 string) {
	ri := feed.LiveData.Plays.CurrentPlay.RunnerIndex

	if feed.LiveData.Linescore.Offense.First != nil {
		on1 = mark(feed.LiveData.Linescore.Offense.First)
	} else {
		on1 = SetRunnerState(ri, 0)
	}

	if feed.LiveData.Linescore.Offense.Second != nil {
		on2 = mark(feed.LiveData.Linescore.Offense.Second)
	} else {
		on2 = SetRunnerState(ri, 1)
	}

	if feed.LiveData.Linescore.Offense.Third != nil {
		on3 = mark(feed.LiveData.Linescore.Offense.Third)
	} else {
		on3 = SetRunnerState(ri, 2)
	}

	return on1, on2, on3
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
