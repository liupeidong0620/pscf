package inilib

func Module_init(file string) (*Ini, error) {
	ini := New()
	ini.file = file
	return ini, nil
}
