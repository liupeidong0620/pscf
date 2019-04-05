package inilib

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

func ReadProperty(args []string, getkeys bool) (interface{}, error) {
	// args format = finename server.arg
	if len(args) < 1 {
		return nil, errors.New("Must provide filename")
	}
	if len(args) != 2 {
		return nil, errors.New("Wrong query parameter.")
	}
	cfg, err := readStream(args[0])
	if err != nil {
		return nil, err
	}
	paths := parsePath(args[1])
	if len(paths) > 2 {
		return nil, errors.New("Wrong query parameter.")
	}

	res, err := readSection(cfg, paths, getkeys)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func WriteProperty(args []string) error {
	// args format = finename server.arg value
	var keyName string
	var sectionName string
	if len(args) < 1 {
		return errors.New("Must provide filename")
	}
	if len(args) != 3 {
		return fmt.Errorf("Wrong set parameter.")
	}
	cfg, err := readStream(args[0])
	if err != nil {
		return err
	}
	paths := parsePath(args[1])
	if len(paths) > 2 {
		return errors.New("Wrong set parameter.")
	}
	if paths[0] == "[]" {
		sectionName = ""
	} else {
		sectionName = paths[0]
	}
	section := cfg.Section(sectionName)
	if len(paths) > 1 {
		keyName = paths[1]
		var arrayFlag bool
		keys := section.KeyStrings()
		if len(keys) > 0 && strings.HasPrefix(keys[0], "#") {
			arrayFlag = true
		}
		// set array
		if paths[1] == "+" {
			if len(keys) > 0 {
				if !arrayFlag {
					return fmt.Errorf("not array type.")
				}
			}
			keyName = "-"
		} else if arrayFlag { // set array key
			i, err := strconv.ParseInt(paths[1], 10, 32)
			if err != nil {
				return err
			}
			if i >= int64(len(keys)) || i < 0 {
				return fmt.Errorf("Array index overflow.")
			}
			keyName = keys[i]
		}

		key := section.Key(keyName)
		key.SetValue(args[2])
	}
	// save
	if args[0] == "-" {
		cfg.WriteTo(os.Stdout)
	} else {
		cfg.SaveTo(args[0])
	}
	return nil
}

func DeleteProperty(args []string) error {
	// args format = finename server.arg
	if len(args) < 1 {
		return errors.New("Must provide filename")
	}
	if len(args) != 2 {
		return errors.New("Wrong delete parameter.")
	}
	cfg, err := readStream(args[0])
	if err != nil {
		return err
	}
	paths := parsePath(args[1])
	if len(paths) > 2 {
		return errors.New("Wrong delete parameter.")
	}
	err = delSection(cfg, paths)
	if err != nil {
		return err
	}

	if args[0] == "-" {
		cfg.WriteTo(os.Stdout)
	} else {
		cfg.SaveTo(args[0])
	}
	return nil
}

func readStream(filename string) (*ini.File, error) {
	if filename == "" {
		return nil, errors.New("Must provide filename.")
	}

	var cfg *ini.File
	var err error

	if filename == "-" {
		stream := bufio.NewReader(os.Stdin)
		cfg, err = ini.Load(ioutil.NopCloser(stream))
	} else {
		cfg, err = ini.Load(filename)
	}
	return cfg, err
}

func readSection(cfg *ini.File, paths []string, getkeys bool) ([]string, error) {
	var values []string
	var path string
	sliceFlag := false

	if paths[0] == "[]" {
		path = ""
	} else {
		path = paths[0]
	}
	section, err := cfg.GetSection(path)
	if err != nil {
		return nil, err
	}
	keys := section.KeyStrings()
	if keys == nil || len(keys) == 0 {
		return nil, nil
	}
	if getkeys == true {
		if len(paths) > 1 {
			return nil, nil
		}
		return keys, nil
	}
	// parse slice
	if strings.HasPrefix(keys[0], "#") {
		sliceFlag = true
	}
	if len(paths) == 1 {
		for _, key := range keys {
			value := section.Key(key).String()
			if sliceFlag {
				values = append(values, "-"+" = "+value)
			} else {
				values = append(values, key+" = "+value)
			}
		}
	} else {
		if sliceFlag {
			i, err := strconv.ParseInt(paths[1], 10, 32)
			if err != nil {
				return nil, err
			}
			if i >= int64(len(keys)) || i < 0 {
				return nil, fmt.Errorf("Array index overflow.")
			}
			value := section.Key(keys[i]).String()
			values = append(values, value)
		} else {
			value := section.Key(paths[1]).String()
			if len(value) != 0 {
				values = append(values, value)
			}
		}
	}
	return values, nil
}

func printReadContext(data interface{}) {
	var values []string
	switch data.(type) {
	case []string:
		values = data.([]string)
		if len(values) <= 0 {
			fmt.Println("null")
			return
		}
	default:
		fmt.Println("null")
		return
	}
	for _, value := range values {
		fmt.Println(value)
	}
}

func delSection(cfg *ini.File, paths []string) error {
	var secName string
	if paths[0] == "[]" {
		secName = ""
	} else {
		secName = paths[0]
	}

	// del section
	if len(paths) == 1 {
		if secName == "" { // del default section key
			section, err := cfg.GetSection(secName)
			if err != nil {
				return err
			}
			keys := section.KeyStrings()
			if keys == nil || len(keys) == 0 {
				return nil
			}
			for _, key := range keys {
				section.DeleteKey(key)
			}
		} else {
			cfg.DeleteSection(secName)
		}
		return nil
	} else {
		// del key
		var delKey string
		section, err := cfg.GetSection(secName)
		if err != nil {
			return err
		}
		delKey = paths[1]
		keys := section.KeyStrings()
		if keys == nil || len(keys) == 0 {
			return nil
		}
		if strings.HasPrefix(keys[0], "#") {
			if paths[1] == "-" {
				delKey = keys[len(keys)-1]
			} else {
				i, err := strconv.ParseInt(paths[1], 10, 32)
				if err != nil {
					return err
				}
				if i >= int64(len(keys)) || i < 0 {
					return fmt.Errorf("Array index overflow.")
				}
				delKey = keys[i]
			}
			// fmt.Println("delKey: ", delKey)
		}
		section.DeleteKey(delKey)
	}
	return nil
}
