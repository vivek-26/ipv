package utils

import (
	"strings"
)

// HelpText returns indented helpTxt string
func HelpText(helpTxt string) string {
	lines := strings.Split(helpTxt, "\n")
	var indentedHelpTxt string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = "  " + line + "\n" // 2 spaces
		indentedHelpTxt += line
	}

	return indentedHelpTxt
}
