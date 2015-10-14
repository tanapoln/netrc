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
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "netrc-path",
			Value: defaultNetrc(),
			Usage: "Path to the netrc file",
		},
		cli.BoolFlag{
			Name:  "no-machine, n",
			Usage: "disable display of machine values",
		},
		cli.BoolFlag{
			Name:  "login, l",
			Usage: "toggle display of login values",
		},
		cli.BoolFlag{
			Name:  "password, p",
			Usage: "toggle display of password values",
		},
	}
	app.CommandNotFound = func(c *cli.Context, cmd string) {
		exit(c, ErrInvalidCommand)
	}
	app.Action = listCommand
	app.BashComplete = machineCompletion

	app.EnableBashCompletion = true
	app.Run(os.Args)
}

type formattableMachine struct {
	*netrc.Machine
	*machineFormat
}

type machineFormat struct {
	showMachine  bool
	showLogin    bool
	showPassword bool
}

func (m formattableMachine) Print() {
	var v []string

	if m.showMachine {
		v = append(v, m.Machine.Name)
	}
	if m.showLogin {
		v = append(v, m.Machine.Login)
	}
	if m.showPassword {
		v = append(v, m.Machine.Password)
	}

	fmt.Print(strings.Join(v, " "))
}

func listMachines(c *cli.Context, mf *machineFormat) {
	filter, netrcFile := cmdSetup(c)
	machines, _, err := netrc.ParseFile(netrcFile)
	if err != nil {
		exit(c, ErrInvalidNetrc)
	}

	printMachines(filterMachines(machines, filter), mf)
}

func listCommand(c *cli.Context) {
	mf := &machineFormat{
		showMachine:  !c.GlobalBool("no-machine"),
		showLogin:    c.GlobalBool("login"),
		showPassword: c.GlobalBool("password"),
	}
	listMachines(c, mf)
}

func machineCompletion(c *cli.Context) {
	mf := &machineFormat{
		showMachine:  true,
		showLogin:    false,
		showPassword: false,
	}
	listMachines(c, mf)
}

func cmdSetup(c *cli.Context) (filter, netrcPath string) {
	return c.Args().First(), c.GlobalString("netrc-path")
}

func printMachines(machines []*netrc.Machine, mf *machineFormat) {
	printNewLine := false

	for _, m := range machines {
		if printNewLine {
			fmt.Println()
		}
		printNewLine = true

		formattableMachine{m, mf}.Print()
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
