package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "upamune"
	app.Email = "jajkeqos@gmail.com"
	app.Usage = ""
	app.Action = doBlock

	app.Flags = GlobalFlags

	app.Run(os.Args)
}
