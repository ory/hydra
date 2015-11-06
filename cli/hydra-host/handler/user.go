package handler

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"
	"github.com/pborman/uuid"
	"log"
)

type User struct {
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

func (c *User) Create(ctx *cli.Context) {
	email := ctx.Args().First()
	if email == "" {
		log.Fatalf("Please provide an email address.")
	}
	password := ctx.String("password")
	if password == "" {
		password = getPassword()
	}

	c.Ctx.Start()
	user, err := c.Ctx.Accounts.Create(uuid.New(), email, password, "{}")
	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Printf(`Created user as "%s".`+"\n", user.GetID())

	if ctx.Bool("as-superuser") {
		if err := c.Ctx.Policies.Create(superUserPolicy(user.GetID())); err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Printf(`Granted superuser privileges to user "%s".`+"\n", user.GetID())
	}
}
