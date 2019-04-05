package yamllib

func Module_init(file string) (*Yaml, error) {
	yaml := New()
	yaml.file = file
	return yaml, nil
}
