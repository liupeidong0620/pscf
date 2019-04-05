package jsonlib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
)

func WriteProperty(args []string) error {
	var raw bool
	var val string
	var err error
	var input []byte
	var outb []byte
	var nodePath string

	if len(args) < 1 {
		return errors.New("Must provide filename")
	}
	if len(args) != 3 {
		return errors.New("Wrong set parameter.")
	}
	// read file
	input, err = readStream(args[0])
	if err != nil {
		return err
	}
	ss := parsePath(args[1])
	nodePath = strings.Join(ss, ".")

	val = args[2]
	if val[0] == '"' && val[len(val)-1] == '"' {
		val = val[1 : len(val)-1]
		raw = false
	} else {
		switch val {
		default:
			if len(val) > 0 {
				if (val[0] >= '0' && val[0] <= '9') || val[0] == '-' {
					if _, err = strconv.ParseFloat(val, 64); err == nil {
						raw = true
					}
				}
			}
		case "true", "false", "null":
			raw = true
		}
	}

	// set json data
	opts := &sjson.Options{}
	opts.Optimistic = true
	opts.ReplaceInPlace = true
	if raw {
		// set as raw block
		outb, err = sjson.SetRawBytesOptions(
			input, nodePath, []byte(val), opts)
	} else {
		// set as a string
		outb, err = sjson.SetBytesOptions(input, nodePath, val, opts)
	}

	if err != nil {
		return err
	}

	// format json data
	outb = pretty.Pretty(outb)

	// save data
	return writeJsonData(args[0], outb)
}

func DeleteProperty(args []string) error {
	// args format = finename server.arg
	if len(args) < 1 {
		return errors.New("Must provide filename")
	}
	if len(args) != 2 {
		return errors.New("Wrong delete parameter.")
	}
	// read file
	input, err := readStream(args[0])
	if err != nil {
		return err
	}
	ss := parsePath(args[1])
	nodePath := strings.Join(ss, ".")
	// del json data
	outb, err := sjson.DeleteBytes(input, nodePath)
	if err != nil {
		return err
	}
	// format json data
	outb = pretty.Pretty(outb)

	// save data
	return writeJsonData(args[0], outb)
}

func ReadProperty(args []string, indent bool) (interface{}, error) {
	// args format = finename server.arg
	if len(args) < 1 {
		return nil, errors.New("Must provide filename")
	}
	if len(args) != 2 {
		return nil, errors.New("Wrong query parameter.")
	}
	// read file
	input, err := readStream(args[0])
	if err != nil {
		return nil, err
	}
	ss := parsePath(args[1])
	nodePath := strings.Join(ss, ".")
	// get json data
	res := gjson.GetBytes(input, nodePath)
	if res.Raw == "" {
		return "null", nil
	}
	if res.Type == gjson.String {
		return res.Str, nil
	} else if indent {
		return pretty.Pretty([]byte(res.Raw)), nil
	} else {
		return pretty.Ugly([]byte(res.Raw)), nil
	}

	return nil, nil
}

func printJsonData(data interface{}) {
	var str string
	switch data.(type) {
	case string:
		str = data.(string)
	case []byte:
		str = string((data.([]byte))[:])
	}
	fmt.Println(str)
}

func readStream(filename string) ([]byte, error) {
	if filename == "" {
		return nil, errors.New("Must provide filename.")
	}

	var input []byte
	var err error

	if filename == "-" {
		input, err = ioutil.ReadAll(os.Stdin)
	} else {
		input, err = ioutil.ReadFile(filename)
	}
	return input, err
}

func writeJsonData(filename string, data []byte) error {
	var err error
	var f *os.File

	if filename == "" {
		return errors.New("Must provide filename.")
	}
	if filename != "-" {
		f, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		f.Write(data)
	} else {
		fmt.Println(string(data[:]))
	}

	return nil
}
