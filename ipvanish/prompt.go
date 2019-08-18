package ipvanish

import (
	"github.com/manifoldco/promptui"

	"github.com/vivek-26/ipv/reporter"
)

// SelectServerPrompt asks user to choose a server to connect to
// and returns the hostname of the chosen server.
func SelectServerPrompt(servers *[]IPVServer, n int) string {
	// Select template
	selectTemplates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "> {{ . | green }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "Selected Server: {{ .Hostname | green }}",
	}

	// Items for prompt
	items := func() []IPVServer {
		if len(*servers) <= n {
			return *servers
		}

		return (*servers)[:n]
	}()

	prompt := promptui.Select{
		Label:     "Choose a server",
		Items:     items,
		Templates: selectTemplates,
		Size:      n,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		reporter.Error(err)
	}

	return (*servers)[idx].Hostname
}
