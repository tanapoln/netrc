package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"code.google.com/p/go-netrc/netrc"
	"github.com/codegangsta/cli"
)

const (
	ErrInvalidCommand = iota
	ErrInvalidNetrc
	ErrMachineNotFound
)

func main() {
	app := cli.NewApp()
	app.Name = "netrc"
	app.Usage = "Manage your netrc file."
	app.Author = "naaman@heroku.com"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "netrc-path",
			Value: defaultNetrc(),
			Usage: "Path to the netrc file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "List netrc entries.",
			Action:    listCommand,
		},
		{
			Name:      "password",
			ShortName: "p",
			Usage:     "Show password for a machine.",
			Action:    passwordCommand,
		},
	}
	app.CommandNotFound = func(c *cli.Context, cmd string) {
		exit(c, ErrInvalidCommand)
	}
	app.Action = listCommand

	app.EnableBashCompletion = true
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
	filter, netrcFile := cmdSetup(c)
	machines, _, err := netrc.ParseFile(netrcFile)
	if err != nil {
		exit(c, ErrInvalidNetrc)
	}

	printMachines(filterMachines(machines, filter), false)
}

func passwordCommand(c *cli.Context) {
	filter, netrcFile := cmdSetup(c)
	if filter == "" {
		exit(c, ErrInvalidCommand)
	}

	machine, err := netrc.FindMachine(netrcFile, filter)
	if err != nil {
		if err.Error() == "no machine found" {
			fmt.Printf("%s does not contain %s\n", netrcFile, filter)
			os.Exit(ErrMachineNotFound)
		} else {
			exit(c, ErrInvalidNetrc)
		}
	}

	printMachines([]*netrc.Machine{machine}, true)
}

func cmdSetup(c *cli.Context) (filter, netrcPath string) {
	return c.Args().First(), c.GlobalString("netrc-path")
}

func printMachines(machines []*netrc.Machine, showPw bool) {
	printNewLine := false

	for _, m := range machines {
		if printNewLine {
			fmt.Println()
		}
		printNewLine = true

		formattableMachine{m, showPw}.Print()
	}
}

func filterMachines(machines []*netrc.Machine, filter string) []*netrc.Machine {
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

func exit(c *cli.Context, e int) {
	cli.ShowAppHelp(c)
	os.Exit(e)
}
