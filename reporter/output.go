package reporter

import (
	"fmt"
	"os"

	au "github.com/logrusorgru/aurora"
)

// Info prints general messages to console
func Info(msg interface{}) {
	fmt.Println(au.BrightCyan(msg).Bold())
}

// Warn prints warning messages to console
func Warn(msg interface{}) {
	fmt.Println(au.BrightYellow(msg).Bold())
}

// Success prints success messages to console
func Success(msg interface{}) {
	fmt.Println(au.BrightGreen(msg).Bold())
}

// Error prints error messages to console
func Error(msg interface{}) {
	fmt.Println("Error: ", au.BrightRed(msg).Bold())
	os.Exit(1)
}
