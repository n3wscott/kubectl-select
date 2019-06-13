package main

import (
	"encoding/json"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"log"
	"os/exec"
	"strings"
)

func cmd(cmdLine string) ([]byte, error) {
	cmdSplit := strings.Split(cmdLine, " ")
	cmd := cmdSplit[0]
	args := cmdSplit[1:]

	return exec.Command(cmd, args...).Output()
}

type K8sContext struct {
	Cluster string `json:"cluster"`
	User    string `json:"user"`
}

type K8sNamedContext struct {
	Name    string     `json:"name"`
	Context K8sContext `json:"context"`
}

type K8sConfig struct {
	Contexts       []K8sNamedContext `json:"contexts"`
	CurrentContext string            `json:"current-context"`
}

func getConfig() *K8sConfig {
	bytes, err := cmd("kubectl config view -o json")
	if err != nil {
		panic(err)
	}
	cfg := &K8sConfig{}
	if err := json.Unmarshal(bytes, cfg); err != nil {
		panic(err)
	}
	return cfg
}

func main() {

	cfg := getConfig()

	table := tui.NewTable(0, 0)
	table.SetColumnStretch(0, 1)
	table.SetColumnStretch(1, 4)
	table.SetColumnStretch(2, 4)
	table.SetColumnStretch(3, 4)
	table.SetFocused(true)

	table.AppendRow(
		tui.NewLabel("SELECTED"),
		tui.NewLabel("NAME"),
		tui.NewLabel("CLUSTER"),
		tui.NewLabel("USER"),
	)

	for i, c := range cfg.Contexts {
		selected := ""
		if c.Name == cfg.CurrentContext {
			selected = "*"
			table.Select(i + 1)
		}
		table.AppendRow(
			tui.NewLabel(selected),
			tui.NewLabel(c.Name),
			tui.NewLabel(c.Context.Cluster),
			tui.NewLabel(c.Context.User),
		)
	}

	status := tui.NewStatusBar("")
	status.SetPermanentText(`ESC or 'q' to QUIT`)

	root := tui.NewVBox(
		table,
		tui.NewSpacer(),
		status,
	)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	table.OnItemActivated(func(t *tui.Table) {
		_, err := cmd(fmt.Sprintf("kubectl config use-context %s", cfg.Contexts[t.Selected()-1].Name))
		if err != nil {
			panic(err)
		}
		ui.Quit()
		fmt.Printf("selected %s\n", cfg.Contexts[t.Selected()-1].Name)
	})

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("q", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
