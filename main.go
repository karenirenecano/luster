package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/zladovan/luster/fb"
	"github.com/zladovan/luster/render"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	// define application cli interface
	app := &cli.App{
		Name:  "luster",
		Usage: "fetch some data from facebook",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "User `EMAIL` used for login",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "User `PASSWORD` used for login",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "fans",
				Usage:     "fetch all users who likes or follows your page",
				ArgsUsage: "page_name",
				Action:    fetchFans,
			},
		},
		Writer: os.Stderr,
	}

	// run application
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// resolveSession gets user email and password from flags or ask for them via interactive input if not set
// and uses them login to Facebook
func resolveSession(c *cli.Context) (*fb.Session, error) {
	email := c.String("user")
	pass := c.String("password")

	// if email is not set ask user for not empty input
	if email == "" {
		email = readInput("Enter your email: ", false)
	}
	if email == "" {
		cli.ShowSubcommandHelp(c)
		return nil, fmt.Errorf("Email cannot be empty")
	}

	// if password is not set ask user for not empty input
	if pass == "" {
		pass = readInput("Enter password: ", true)
	}
	if email == "" {
		cli.ShowSubcommandHelp(c)
		return nil, fmt.Errorf("Password cannot be empty")
	}

	// login
	session, err := fb.Login(email, pass)
	if err != nil {
		return nil, fmt.Errorf("Unable to login", err)
	}

	return session, nil
}

// fetchFans scrapes all fans of Facebook page and prints them to stdout.
func fetchFans(c *cli.Context) error {
	pgname := c.Args().Get(0)
	if pgname == "" {
		cli.ShowSubcommandHelp(c)
		return fmt.Errorf("Required argument \"page_name\" is missing")
	}

	session, err := resolveSession(c)
	if err != nil {
		return err
	}

	// open page
	page, err := session.OpenPage(pgname)
	if err != nil {
		return fmt.Errorf("Unable to open page "+pgname, err)
	}

	// fetch users from page
	fans, err := page.FetchFans()
	if err != nil {
		return fmt.Errorf("Unable to fetch fans of page"+page.Name, err)
	}

	// print render output
	fmt.Print(render.Csv(fans))

	// all good
	return nil
}

// readInput prints given label and wait for user input which is then returned.
// Input can be hidden which means characters typed by user are not shown in console.
// After input is confirmed with [Enter] input line is deleted.
func readInput(label string, hidden bool) string {
	// print label to stderr to do not affect stdout which can be then easily stored e.g. to file
	fmt.Fprintf(os.Stderr, label)

	var input string

	if hidden {
		// do not dispaly what user is typing
		passBytes, _ := terminal.ReadPassword(int(syscall.Stderr))
		input = string(passBytes)
	} else {
		// display what user is typing
		scanner := bufio.NewScanner(os.Stderr)
		scanner.Scan()
		input = scanner.Text()
		// move cursor up as new line is started after scanning
		fmt.Fprintf(os.Stderr, "\033[A")
	}

	// clear line
	fmt.Fprintf(os.Stderr, "\r")
	fmt.Fprintf(os.Stderr, strings.Repeat(" ", len(label)+len(input)))
	fmt.Fprintf(os.Stderr, "\r")

	return input
}
