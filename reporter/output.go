package reporter

import (
	"fmt"

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

// Error prints error messages to console
func Error(msg interface{}) {
	fmt.Println(au.BrightRed(msg).Bold())
}
