package main

import (
	"flag"
	"code.google.com/p/go-netrc/netrc"
	"os"
	"fmt"
	"os/user"
	"strings"
	"github.com/codegangsta/cli"
)

const (
	ErrInvalidCommand = iota
	ErrInvalidNetrc
)

func main() {
	app := cli.NewApp()
	app.Name = "netrc"
	app.Usage = "Manage your netrc file."
	app.Flags = []cli.Flag {
		cli.StringFlag {
			Name: "netrc-path",
			Value: defaultNetrc(),
			Usage: "Path to the netrc file",
		},
	}
	app.Commands = []cli.Command {
		{
			Name: "list",
			ShortName: "l",
			Usage: "List netrc entries.",
			Action: listCommand,
		},
		{
			Name: "password",
			ShortName: "p",
			Usage: "Show password for a machine.",
			Action: passwordCommand,
		},
	}

	app.Run(os.Args)
}

type formattableMachine struct {
	*netrc.Machine
	showPw bool
}

func (m formattableMachine) Print() {
	var v []string
	if m.showPw {
		v = append(v, m.Machine.Password)
	} else {
		v = append(v, m.Machine.Name, m.Machine.Login)
	}

	fmt.Print(strings.Join(v, " "))
}

func listCommand(c *cli.Context) {
	filter, machines := cmdSetup(c)
	printMachines(machines, filter, false)
}

func passwordCommand(c *cli.Context) {
	filter, machines := cmdSetup(c)
	if filter == "" {
		fmt.Println("Provide a machine name.")
		os.Exit(ErrInvalidCommand)
	}
	printMachines(machines, filter, true)
}

func cmdSetup(c *cli.Context) (string, []*netrc.Machine) {
	filter := c.Args().First()
	machines, _, err := netrc.ParseFile(c.GlobalString("netrc-path"))
	if err != nil {
		flag.Usage()
		os.Exit(ErrInvalidNetrc)
	}
	return filter, machines
}

func printMachines(machines []*netrc.Machine, filter string, showPw bool) {
	printNewLine := false

	for _, m := range filterMachines(machines, filter) {
		if printNewLine {
			fmt.Println()
		}
		printNewLine = true

		formattableMachine{m, showPw}.Print()
	}
}

func filterMachines(machines []*netrc.Machine, filter string) ([]*netrc.Machine) {
	var filteredMachines []*netrc.Machine
	if filter == "" {
		return machines
	}
	for _, m := range machines {
		if m.Name == filter {
			filteredMachines = append(filteredMachines, m)
		}
	}
	return filteredMachines
}

func defaultNetrc() string {
	if u, err := user.Current(); err == nil {
		netrcPath := u.HomeDir + "/.netrc"
		if _, err := os.Stat(netrcPath); err == nil {
			return netrcPath
		}
	}
	return ""
}
