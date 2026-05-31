package screen

import (
	"fmt"
	"os"

	term "github.com/charmbracelet/x/term"
)

// Deprecated: just use tea.WindowSizeMsg in Update to get this anytime the window resizes
// GetTerminalDimensions retrieves the current terminal's width and height in characters.
// It uses the file descriptor of stdout to determine the terminal size.
func GetTerminalDimensions() (width int, height int, err error) {
	fd := os.Stdout.Fd()
	width, height, err = term.GetSize(fd)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get terminal size: %w", err)
	}
	return width, height, nil
}

// Deprecated
// CalculateScorebugs calculates how many scorebugs can be displayed on the screen given the dimensions of a single scorebug.
// It takes the width and height of a scorebug (in characters) and returns the number of columns, rows, and total scorebugs that can fit on the screen in a grid pattern.
func CalculateScorebugs(scorebugWidth, scorebugHeight int) (cols int, rows int, total int, err error) {
	if scorebugWidth <= 0 || scorebugHeight <= 0 {
		return 0, 0, 0, fmt.Errorf("invalid scorebug dimensions provided")
	}

	termWidth, termHeight, err := GetTerminalDimensions()
	if err != nil {
		return 0, 0, 0, err
	}

	cols = termWidth / scorebugWidth
	rows = termHeight / scorebugHeight
	total = cols * rows

	return cols, rows, total, nil
}

// Determines amount of scorebugs can fit on a paginated list in a given windowHeight
func GetSchedulePageSize(scorebugHeight int, margin int, windowHeight int) (rows int, err error) {
	if scorebugHeight <= 0 {
		return 0, fmt.Errorf("invalid scorebug dimensions provided")
	}

	rows = (windowHeight) / (scorebugHeight + (margin * 2))

	return rows, nil
}
