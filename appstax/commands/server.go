package commands

import (
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/hosting"
	"appstax-cli/appstax/term"
	"github.com/codegangsta/cli"
	"strconv"
)

func DoServer(c *cli.Context) {
	useOptions(c)
	if !config.Exists() {
		term.Println("Can't find appstax.conf. Run 'appstax init' to initialize before deploying.")
		return
	}
	loginIfNeeded()
	
	args := c.Args()
	if len(args) == 0 {
		term.Println("Too few arguments. Usage: appstax server create|delete|status|start|stop|log")
		return
	}

	operation := args[0]
	switch operation {
	case "create":
		selectSubdomainIfNeeded()
		accessCode := term.GetString("Please enter early access code")
		err := hosting.CreateServer(accessCode)
		if err == nil {
			term.Println("Server created successfully!")
		} else {
			term.Println("Error creating server:")
			term.Println(err.Error())
		}
	case "delete":
		err := hosting.DeleteServer()
		if err == nil {
			term.Println("Server deleted!")
		} else {
			term.Println("Error deleting server:")
			term.Println(err.Error())
		}
	case "status":
		status, err := hosting.GetServerStatus()
		if err == nil {
			term.Println("Server status: "+status.Status)
		} else {
			term.Println("Error getting server status:")
			term.Println(err.Error())
		}
	case "start":
		err := hosting.SendServerAction(operation)
		if err == nil {
			term.Println("Server started!")
		} else {
			term.Println("Error starting server:")
			term.Println(err.Error())
		}
	case "stop":
		err := hosting.SendServerAction(operation)
		if err == nil {
			term.Println("Server stopped!")
		} else {
			term.Println("Error stopping server:")
			term.Println(err.Error())
		}
	case "log", "logs":
		nlines := int64(10)
		if len(args) >= 2 {
			n, err := strconv.ParseInt(args[1], 10, 64)
			if err == nil {
				nlines = n
			}
		}
		log, err := hosting.GetServerLog(nlines)
		if err == nil {
			term.Layout(false)
			term.Dump(log)
		} else {
			term.Println("Error getting server log:")
			term.Println(err.Error())
		}
	default:
		term.Printf("Unknown server operation '%s'\n", operation)
	}
}
