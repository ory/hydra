package handler

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"
	"github.com/pborman/uuid"
)

type Account struct {
	Ctx *Context
}

func getPassword() (password string) {
	fmt.Printf("Password: ")
	password = string(gopass.GetPasswd())
	if password == "" {
		fmt.Println("You did not provide a password. Please try again.")
		return getPassword()
	}

	fmt.Printf("Confirm password: ")
	if password != string(gopass.GetPasswd()) {
		fmt.Println("Password and confirmation do not match. Please try again.")
		return getPassword()
	}
	return
}

func (c *Account) Create(ctx *cli.Context) error {
	email := ctx.Args().First()
	if email == "" {
		return fmt.Errorf("Please provide an email address.")
	}
	password := ctx.String("password")
	if password == "" {
		password = getPassword()
	}

	c.Ctx.Start()
	user, err := c.Ctx.Accounts.Create(uuid.New(), email, password, "{}")
	if err != nil {
		return fmt.Errorf("Could not create account because %s", err)
	}

	fmt.Printf(`Created account as "%s".`+"\n", user.GetID())
	if ctx.Bool("as-superuser") {
		if err := c.Ctx.Policies.Create(superUserPolicy(user.GetID())); err != nil {
			return fmt.Errorf("Could not create policy for account because %s", err)
		}
		fmt.Printf(`Granted superuser privileges to account "%s".`+"\n", user.GetID())
	}
	return nil
}
