package jsonlib

func Module_init(file string) (*Json, error) {
	json := New()
	json.file = file
	return json, nil
}
