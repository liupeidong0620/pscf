package main

import (
	"fmt"

	iniset "github.com/liupeidong0620/pscf/ini"
	jsonset "github.com/liupeidong0620/pscf/json"
	yamlset "github.com/liupeidong0620/pscf/yaml"
)

type yamlAttribute struct {
	ignoreField string
}

/*
   interface
*/
type modules_st interface {
	Set(params ...interface{}) error

	Get(params ...interface{}) (interface{}, error)
	GetKeys(params ...interface{}) (interface{}, error)
	Del(params ...interface{}) error
	Save() error
}

func modules_init(mod, file string) (modules_st, error) {

	var module modules_st
	var err error

	switch mod {
	case YAMLTOOL:
		module, err = yamlset.Module_init(file)
		if err != nil {
			return nil, err
		}
	case JSONTOOL:
		module, err = jsonset.Module_init(file)
		if err != nil {
			return nil, err
		}
	case INITOOL:
		module, err = iniset.Module_init(file)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported configure file type: %s.", mod)
	}

	return module, err
}

func setYamlCommentEnable(enable bool) {
	yamlset.SetCommentEnable(enable)
}

func set_module_attribute(mod string, ctx interface{}) error {
	switch mod {
	case YAMLTOOL:
		if yamlAttr, ok := ctx.(yamlAttribute); ok {
			return yamlset.SetIgnoreField(yamlAttr.ignoreField)
		}
	case JSONTOOL:
	case INITOOL:
	}
	return nil
}
