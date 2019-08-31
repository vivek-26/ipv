package reporter

import (
	"fmt"
	"os"

	"github.com/briandowns/spinner"
	au "github.com/logrusorgru/aurora"
)

// Spinner is a wrapper for logging messages with a nice spinner
type Spinner struct {
	Spin *spinner.Spinner
}

// Info prints general messages to console and starts spinner
func (s *Spinner) Info(msg interface{}) {
	fmt.Printf("%v  ", au.Cyan(msg).Bold())
	s.Spin.Start()
}

// Success stops spinner and prints `✓` to console to indicate
// completion of a task.
func (s *Spinner) Success() {
	s.Spin.Stop()
	fmt.Printf("%v\n", au.BrightGreen("✓").Bold())
}

// Error stops spinner and prints `✗` to console to indicate
// failure of a task and terminates program.
func (s *Spinner) Error() {
	s.Spin.Stop()
	fmt.Printf("%v\n", au.BrightRed("✗").Bold())
	os.Exit(1)
}

// Info prints general messages to console
func Info(msg interface{}) {
	fmt.Printf("%v\n", au.BrightCyan(msg).Bold())
}

// Warn prints warning messages to console
func Warn(msg interface{}) {
	fmt.Printf(
		"%v %v\n",
		au.BrightYellow("Warn:").Bold(),
		au.BrightYellow(msg).Bold(),
	)
}

// Success prints success messages to console
func Success(msg interface{}) {
	fmt.Printf("%v\n", au.BrightGreen(msg).Bold())
}

// Error prints error messages to console and terminates program
func Error(msg interface{}) {
	fmt.Printf(
		"%v %v\n",
		au.BrightRed("Error:").Bold(),
		au.BrightRed(msg).Bold(),
	)
	os.Exit(1)
}
