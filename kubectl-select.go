/*
Copyright 2020 Scott Nichols <author@n3wscott.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

	app := tview.NewApplication()
	list := tview.NewList()
	list.SetBorder(true).SetTitle("Select a Context")

	doSelect := func(i int) {
		_, err := cmd(fmt.Sprintf("kubectl config use-context %s", cfg.Contexts[i].Name))
		if err != nil {
			panic(err)
		}
		app.Stop()
		fmt.Printf("selected %s\n", cfg.Contexts[i].Name)
	}

	shortcut := 'a'
	for i, c := range cfg.Contexts {
		selected := ""
		if c.Name == cfg.CurrentContext {
			selected = "[current]"
			// TODO: can we select this ui element?
		}
		// TODO: if someone runs shortcut into q, not sure which shortcut item wins.
		i := i
		list.AddItem(fmt.Sprintf("%s %s", selected, c.Name), fmt.Sprintf("%s@%s", c.Context.User, c.Context.Cluster), shortcut, func() {
			doSelect(i)
		})
		shortcut++
	}

	list.AddItem("Quit", "Press `q` or `ESC` to exit", 'q', func() {
		app.Stop()
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			fmt.Printf("no selection; context unchanged\n")
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
