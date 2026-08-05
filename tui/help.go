package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

var content = `
	t               - show next sample
	T               - show preview sample
	p			    - switch to process view
	s			    - switch to system view
	c			    - switch to cgroup view
	b               - open dialog to search specific sample
	/               - show/hide filter for process/cgroup view

process view:
	s               - show/hide sort
	<Tab>           - switch sub category (genernal, cpu, memory etc)
	<Shift>+<Tab>   - reverse switch
	z               - zoom into cgroup view filtered by process
	F               - toggle process tree mode

system view:
	<Tab>           - switch sub category (cpu, memory, disk etc)
	<Shift>+<Tab>   - reverse switch

cgroup view:
	<Tab>           - switch sub category (genernal, cpu, memory etc)
	<Shift>+<Tab>   - reverse switch
	<Enter>         - collapse/expand cgroup tree
	z               - zoom into process view filtered by cgroup

	Type 'ESC' to close
`

type Help struct {
	*tview.TextView
}

func NewHelp() *Help {

	help := &Help{
		TextView: tview.NewTextView(),
	}

	help.SetTitle("Help").SetBorder(true).SetTitleAlign(tview.AlignLeft)

	fmt.Fprint(help, content)

	return help
}
