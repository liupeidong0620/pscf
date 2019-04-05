package jsonlib

import ()

type Json struct {
	file string
}

func New() *Json {
	self := new(Json)
	return self
}

func (self *Json) SetFile(file string) {
	self.file = file
}

func (self *Json) generateArgs(params ...interface{}) ([]string, error) {
	var args []string
	args = append(args, self.file)
	for i := 0; i < len(params); i++ {
		arg := params[i].(string)
		args = append(args, arg)
	}
	return args, nil
}

func (self *Json) Set(params ...interface{}) error {
	args, err := self.generateArgs(params...)
	if err != nil {
		return err
	}
	err = WriteProperty(args)
	if err != nil {
		return err
	}
	return nil
}

func (self *Json) Get(params ...interface{}) (interface{}, error) {
	args, err := self.generateArgs(params...)
	if err != nil {
		return nil, err
	}
	data, err := ReadProperty(args, true)
	if err != nil {
		return nil, err
	}
	printJsonData(data)

	return data, nil
}

func (self *Json) Del(params ...interface{}) error {
	args, err := self.generateArgs(params...)
	if err != nil {
		return err
	}
	err = DeleteProperty(args)
	if err != nil {
		return err
	}
	return nil
}

func (self *Json) GetKeys(params ...interface{}) (interface{}, error) {
	return nil, nil
}

func (self *Json) Save() error {
	return nil
}
