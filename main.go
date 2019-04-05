package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli"
)

var configFile string
var environList string
var setType string

const (
	/* Operation profile type */
	YAMLTOOL string = "yaml"
	JSONTOOL string = "json"
	INITOOL  string = "ini"

	/*tool name*/
	TOOLNAME string = "pscf"

	/* commands */
	CommandDel     string = "del"
	CommandGet     string = "get"
	CommandSet     string = "set"
	CommandGetKeys string = "getkeys"

	CommandIgnore string = "tostr"

	CommandFile string = "config"

	CommandType    string = "type"
	CommandVerbose string = "verbose,v"
	CommandComment string = "comment"

	Version   string = "1.0.0"
	CurSystem string = runtime.GOOS + " (" + runtime.GOARCH + ")"
)

/* command line flags */
var commandBaseFlags []cli.Flag = []cli.Flag{
	cli.StringFlag{
		Name:        "config, c",
		Usage:       "Load configuration from `FILE`",
		Destination: &configFile,
	},
	cli.StringFlag{
		Name:        CommandDel,
		Usage:       "Delete a node in the configuration file.(--del node)",
		Destination: &environList,
	},
	cli.StringFlag{
		Name:        CommandGet,
		Usage:       "Get the value of the node in the configuration file.(--get node)",
		Destination: &environList,
	},
	cli.StringFlag{
		Name:        CommandSet,
		Usage:       "Set the value of the node and add a new node if the node does not exist.(--set node=value)",
		Destination: &environList,
	},
	cli.BoolFlag{
		Name:  CommandVerbose,
		Usage: "Debug mode.(--verbose, -v)",
	},
}

var commandYamlFlags []cli.Flag = []cli.Flag{
	cli.StringFlag{
		Name:  CommandIgnore,
		Usage: "Force some tags into strings.(--tostr \"yes,no\")",
	},
	cli.StringFlag{
		Name:        CommandGetKeys,
		Usage:       "Get the keys value of the node in the configuration file.(--getkeys node)",
		Destination: &environList,
	},
	cli.BoolFlag{
		Name:  CommandComment,
		Usage: "Load comment enable.(--comment)",
	},
}

var commandIniFlags []cli.Flag = []cli.Flag{
	cli.StringFlag{
		Name:        CommandGetKeys,
		Usage:       "Get the keys value of the node in the configuration file.(--getkeys node)",
		Destination: &environList,
	},
}

var subCommands []cli.Command = []cli.Command{
	{
		Name:  YAMLTOOL,
		Usage: "Provides tools for set, delete, and get operations for a node in YAML format.",
		UsageText: "pscf yaml [-c config.yaml] [--set node.node1.node2=value]\n" +
			"   pscf yaml [-c config.yaml] [--del node.node1.node2]\n" +
			"   pscf yaml [-c config.yaml] [--get node.node1.node2]\n" +
			"   pscf yaml [-c config.yaml] [--getkeys node.node1.node2]",
		Flags: append(commandBaseFlags, commandYamlFlags...),
		Action: func(c *cli.Context) error {
			if len(os.Args) < 3 || c.NArg() != 0 {
				cli.ShowCommandHelpAndExit(c, YAMLTOOL, -1)
			}

			err := run(c, YAMLTOOL, os.Args[2:])
			if err != nil {
				return err
			}

			return nil
		},
	},
	{
		Name:  JSONTOOL,
		Flags: commandBaseFlags,
		Usage: "Provides tools for set, delete, and get operations for a node in JSON format.",
		UsageText: "pscf json [-c config.yaml] [--set node.node1.node2=value]\n" +
			"   pscf json [-c config.yaml] [--del node.node1.node2]\n" +
			"   pscf json [-c config.yaml] [--get node.node1.node2]",
		Action: func(c *cli.Context) error {
			if len(os.Args) < 3 || c.NArg() != 0 {
				cli.ShowCommandHelpAndExit(c, JSONTOOL, -1)
			}

			err := run(c, JSONTOOL, os.Args[2:])
			if err != nil {
				return err
			}

			return nil
		},
	},
	{
		Name:  INITOOL,
		Usage: "Provides tools for set, delete, and get operations for section or key in INI format.",
		Flags: append(commandBaseFlags, commandIniFlags...),
		UsageText: "pscf ini [-c config.yaml] [--set section.key=value]\n" +
			"   pscf ini [-c config.yaml] [--del section.key]\n" +
			"   pscf ini [-c config.yaml] [--get section.key]\n" +
			"   pscf ini [-c config.yaml] [--getkeys section]",
		Action: func(c *cli.Context) error {
			if len(os.Args) < 3 || c.NArg() != 0 {
				cli.ShowCommandHelpAndExit(c, INITOOL, -1)
			}

			err := run(c, INITOOL, os.Args[2:])
			if err != nil {
				return err
			}
			return nil
		},
	},
}

func main() {

	/* app settings */
	app := cli.NewApp()
	app.Name = TOOLNAME
	app.Usage = "Processing standardized config file."
	app.UsageText = app.Name + " [yaml|json|ini]" + " -h"
	app.Copyright = "(c) 2018 LPD"
	app.Authors = append(app.Authors, cli.Author{
		Name:  "liupeidong",
		Email: "liupeidong0620@163.com",
	})
	app.Version = Version + " " + CurSystem + " Build Date " + app.Compiled.String()
	app.Commands = subCommands

	/* app exec func*/
	app.Action = func(c *cli.Context) error {
		if len(os.Args) < 2 || c.NArg() != 0 {
			cli.ShowAppHelpAndExit(c, -1)
		}
		/*
			switch c.Args().Get(1) {
			case YAMLTOOL:
			case JSONTOOL:
			case INITOOL:
			default:
				cli.ShowAppHelpAndExit(c, -1)
			}*/
		return nil
	}
	/* app start */
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Err: ", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
