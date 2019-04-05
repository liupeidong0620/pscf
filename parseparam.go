package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

/*
   Is the command line parameter combination correct ?
*/
func isParamSet(c *cli.Context, args []string) (string, error) {
	//fmt.Println("args: ", args)
	var extraLen int = 0
	if c.IsSet(CommandFile) == false {
		return "", fmt.Errorf("--config no set")
	} else {
		extraLen += 2
	}
	if c.IsSet("v") {
		extraLen += 1
	}
	if c.IsSet(CommandComment) {
		extraLen += 1
	}
	if c.IsSet(CommandIgnore) {
		extraLen += 2
	}

	argsLen := len(args)
	extraLen += 2
	if c.IsSet(CommandSet) {
		if argsLen == extraLen {
			return CommandSet, nil
		} else {
			return "", fmt.Errorf("Command line format error: pscf -c test --set key=value.")
		}
	} else if c.IsSet(CommandDel) {
		if argsLen == extraLen {
			return CommandDel, nil
		} else {
			return "", fmt.Errorf("Command line format error: pscf -c test --del key.")
		}
	} else if c.IsSet(CommandGet) {
		if argsLen == extraLen {
			return CommandGet, nil
		} else {
			return "", fmt.Errorf("Command line format error: pscf -c test --get key.")
		}
	} else if c.IsSet(CommandGetKeys) {
		if argsLen == extraLen {
			return CommandGetKeys, nil
		} else {
			return "", fmt.Errorf("Command line format error: pscf -c test --getkeys key.")
		}
	} else {
		//fmt.Println(args)
		return "", fmt.Errorf("The command line parameter combination is error: %v", args)
	}
	return "", nil
}

/*
   set node value type
*/
func parseParamType(valType string, value string) (interface{}, error) {
	value = strings.TrimSpace(value)
	flag := false
	var params []string
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		tmp := strings.Trim(value, "[]")
		params = strings.Split(tmp, ",")
		flag = true
	}
	switch valType {
	case "int":
		if flag {
			intParams := []int64{}
			for _, param := range params {
				val, err := strconv.ParseInt(param, 10, 64)
				if err != nil {
					return nil, err
				}
				intParams = append(intParams, val)
			}
			return intParams, nil
		}
		return strconv.ParseInt(value, 10, 64)
	case "string":
		if flag {
			return params, nil
		}
		return value, nil
	case "float":
		if flag {
			floatParams := []float64{}
			for _, param := range params {
				val, err := strconv.ParseFloat(param, 64)
				if err != nil {
					return nil, err
				}
				floatParams = append(floatParams, val)
			}
			return floatParams, nil
		}
		return strconv.ParseFloat(value, 64)
	case "uint":
		if flag {
			uintParams := []uint64{}
			for _, param := range params {
				val, err := strconv.ParseUint(param, 10, 64)
				if err != nil {
					return nil, err
				}
				uintParams = append(uintParams, val)
			}
			return uintParams, nil
		}
		return strconv.ParseUint(value, 10, 64)
	case "bool":
		if flag {
			boolParams := []bool{}
			for _, param := range params {
				val, err := strconv.ParseBool(param)
				if err != nil {
					return nil, err
				}
				boolParams = append(boolParams, val)
			}
			return boolParams, nil
		}
		return strconv.ParseBool(value)
	default:
		return nil, fmt.Errorf("Unsupported type: %s", valType)
	}

	return nil, nil
}

/*
   Parsing command line arguments.
*/
func parseParamNode(c *cli.Context, command string, mod string) (string, interface{}, error) {

	/* get node key-value */
	params := c.String(command)
	if params == "" {
		return "", nil, fmt.Errorf("%s: No value.", command)
	}

	switch command {
	case CommandSet:
		/*
			var valType string
			if mod != YAMLTOOL && mod != INITOOL && mod != JSONTOOL {
				if !c.IsSet(CommandType) {
					return "", nil, fmt.Errorf("--type: no set.")
				}
				// get value type
				valType = c.String(CommandType)
				if valType == "" {
					return "", nil, fmt.Errorf("--type: no Value.")
				}
			} else {
				valType = "string"
			}*/
		/* is value ? */
		if strings.Contains(params, "=") == false {
			return "", nil, fmt.Errorf("--%s %s=value: No value.", command, params)
		}
		/* parse key and value */
		param := strings.Split(params, "=")
		if len(param) != 2 {
			return "", nil, fmt.Errorf("--%s %s=value: No value.", command, params)
		}
		/*paser value type */
		/*
			nodeVal, err := parseParamType(valType, param[1])
			if err != nil {
				return "", nil, err
			}*/
		return param[0], param[1], nil
	case CommandGet, CommandDel, CommandGetKeys:
		// get node
		return params, nil, nil
	default:
		return "", nil, fmt.Errorf("No such command line argument: %s", command)

	}

	return "", nil, nil
}
