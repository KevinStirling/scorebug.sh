package screen_test

import (
	"os"
	"testing"

	"github.com/KevinStirling/scorebug.sh/ui/screen"
	"github.com/charmbracelet/x/term"
)

// TestGetTerminalDimensions tests that GetTerminalDimensions returns valid dimensions when running in a terminal.
func TestGetTerminalDimensions(t *testing.T) {
	// Skip the test if not running in an interactive terminal
	if !term.IsTerminal(os.Stdout.Fd()) {
		t.Skip("Skipping test; not running in a terminal")
	}

	width, height, err := screen.GetTerminalDimensions()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if width <= 0 || height <= 0 {
		t.Fatalf("Invalid terminal dimensions returned: width=%d, height=%d", width, height)
	}
}

// TestCalculateScorebugsValid tests CalculateScorebugs with valid scorebug dimensions.
func TestCalculateScorebugsValid(t *testing.T) {
	// Skip the test if not running in an interactive terminal
	if !term.IsTerminal(os.Stdout.Fd()) {
		t.Skip("Skipping test; not running in a terminal")
	}

	// Obtain terminal dimensions for calculation reference
	termWidth, termHeight, err := screen.GetTerminalDimensions()
	if err != nil {
		t.Fatalf("Failed to get terminal dimensions: %v", err)
	}

	scorebugWidth := 10
	scorebugHeight := 5

	expectedCols := termWidth / scorebugWidth
	expectedRows := termHeight / scorebugHeight
	expectedTotal := expectedCols * expectedRows

	cols, rows, total, err := screen.CalculateScorebugs(scorebugWidth, scorebugHeight)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cols != expectedCols || rows != expectedRows || total != expectedTotal {
		t.Errorf("Expected cols=%d, rows=%d, total=%d; got cols=%d, rows=%d, total=%d", expectedCols, expectedRows, expectedTotal, cols, rows, total)
	}
}

// TestCalculateScorebugsInvalid tests CalculateScorebugs with invalid scorebug dimensions.
func TestCalculateScorebugsInvalid(t *testing.T) {
	// Test invalid width
	if _, _, _, err := screen.CalculateScorebugs(0, 10); err == nil {
		t.Error("Expected error for zero scorebug width, but got nil")
	}

	// Test invalid height
	if _, _, _, err := screen.CalculateScorebugs(10, -1); err == nil {
		t.Error("Expected error for negative scorebug height, but got nil")
	}
}
