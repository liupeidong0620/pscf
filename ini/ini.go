package inilib

import ()

type Ini struct {
	file string
}

func New() *Ini {
	self := new(Ini)
	return self
}

func (self *Ini) SetFile(file string) {
	self.file = file
}

func (self *Ini) generateArgs(params ...interface{}) ([]string, error) {
	var args []string
	args = append(args, self.file)
	for i := 0; i < len(params); i++ {
		arg := params[i].(string)
		args = append(args, arg)
	}
	return args, nil
}

func (self *Ini) Set(params ...interface{}) error {
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

func (self *Ini) Get(params ...interface{}) (interface{}, error) {
	args, err := self.generateArgs(params...)
	if err != nil {
		return nil, err
	}
	data, err := ReadProperty(args, false)
	if err != nil {
		return nil, err
	}
	printReadContext(data)

	return data, nil
}

func (self *Ini) Del(params ...interface{}) error {
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

func (self *Ini) GetKeys(params ...interface{}) (interface{}, error) {
	args, err := self.generateArgs(params...)
	if err != nil {
		return nil, err
	}
	data, err := ReadProperty(args, true)
	if err != nil {
		return nil, err
	}
	printReadContext(data)

	return data, nil
}

func (self *Ini) Save() error {
	return nil
}
