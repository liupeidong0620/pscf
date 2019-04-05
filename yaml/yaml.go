package yamllib

import (
	"fmt"
	"reflect"
	"strings"

	yaml "github.com/liupeidong0620/yaml"
)

type Yaml struct {
	file    string
	versobe bool
}

/*
	Creates and returns a YAML struct.
*/
func New() *Yaml {
	self := new(Yaml)
	yaml.DefaultMapType = reflect.TypeOf(yaml.MapSlice{})
	return self
}

func SetIgnoreField(str string) error {
	if str == "" {
		return fmt.Errorf("The parameter is empty.")
	}
	fields := strings.Split(str, ",")
	for _, field := range fields {
		yaml.IgnoreResolve[field] = true
	}
	return nil
}

func SetCommentEnable(enable bool) {
	yaml.DefaultCommentsEnable = enable
}

func (self *Yaml) generateArgs(params ...interface{}) ([]string, error) {
	var args []string
	args = append(args, self.file)
	for i := 0; i < len(params); i++ {
		arg := params[i].(string)
		args = append(args, arg)
	}
	return args, nil
}

/*
	Sets a YAML setting
*/
func (self *Yaml) Set(params ...interface{}) error {
	args, err := self.generateArgs(params...)
	if err != nil {
		return err
	}
	return WriteProperty(args)
}

/*
   Delete a YAML setting
*/
func (self *Yaml) Del(params ...interface{}) error {
	args, err := self.generateArgs(params...)
	if err != nil {
		return err
	}

	return DeleteProperty(args)
}

/*
	Returns a YAML setting
*/
func (self *Yaml) Get(params ...interface{}) (interface{}, error) {
	args, err := self.generateArgs(params...)
	if err != nil {
		return nil, err
	}

	data, err := ReadProperty(args)
	if err != nil {
		return nil, err
	}
	return ConvReadContext(data, true, false)
}

func (self *Yaml) GetKeys(params ...interface{}) (interface{}, error) {
	args, err := self.generateArgs(params...)
	if err != nil {
		return nil, err
	}
	data, err := ReadProperty(args)
	if err != nil {
		return nil, err
	}
	return ConvReadContext(data, true, true)
}

/*
	Writes changes to the currently opened YAML file.
*/
func (self *Yaml) Save() error {

	return nil
}
