package main

import (
	"github.com/urfave/cli"

	pscfLog "pscf/log"
)

func run(c *cli.Context, mod string, args []string) error {
	var commandErr error
	var verbose bool
	if c.IsSet("v") {
		verbose = true
	} else {
		verbose = false
	}
	// parse command line param
	pscfLog.LogInit(mod, verbose)
	command, err := isParamSet(c, args)
	if err != nil {
		return err
	}
	route, value, err := parseParamNode(c, command, mod)
	if err != nil {
		return err
	}
	// module init

	module, err := modules_init(mod, c.String(CommandFile))
	if err != nil {
		return err
	}
	// ignore some field
	if c.IsSet(CommandIgnore) && mod == YAMLTOOL {
		ignoreData := yamlAttribute{
			ignoreField: c.String(CommandIgnore),
		}
		err := set_module_attribute(mod, ignoreData)
		if err != nil {
			return err
		}
	}
	// load yaml comment
	if c.IsSet(CommandComment) &&
		mod == YAMLTOOL && command != CommandGetKeys {
		setYamlCommentEnable(true)
	}
	// add del set mod get operation
	switch command {
	case CommandDel:
		commandErr = module.Del(route)
	case CommandSet:
		commandErr = module.Set(route, value)
	case CommandGet:
		_, err := module.Get(route)
		if err != nil {
			return err
		}
		return nil
	case CommandGetKeys:
		_, err := module.GetKeys(route)
		if err != nil {
		}
		return nil
	default:
		// error
	}
	if commandErr != nil {
		return commandErr
	}
	module.Save()

	return nil
}
