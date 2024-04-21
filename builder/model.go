package builder

import "strings"

type Model struct {
	Name    string   `yaml:"name"`
	Fields  []Field  `yaml:"fields"`
	Methods []Method `yaml:"methods"`
}

type Method string

func (m Method) String() string {
	if parsed, ok := DefaultMethodNamings[Method(strings.ToLower(string(m)))]; ok {
		return parsed
	}

	// if not found in default namings, just return the method name with first letter capitalized
	runes := []rune(string(m))
	naming := strings.ToUpper(string(runes[0])) + string(runes[1:])

	return naming
}

const (
	MethodCreate Method = "create"
	MethodRead   Method = "read"
	MethodUpdate Method = "update"
	MethodDelete Method = "delete"
)

var DefaultMethodNamings = map[Method]string{
	MethodCreate: "Create",
	MethodRead:   "Read",
	MethodUpdate: "Update",
	MethodDelete: "Delete",
}

type Field struct {
	Name string    `yaml:"name"`
	Type FieldType `yaml:"type"`
	Tags []string  `yaml:"tags"`
}

type Tag struct {
	Key string `yaml:"key"`
	Val string `yaml:"val"`
}

type FieldType string

const (
	FieldTypeID  FieldType = "id"
	FieldTypeInt FieldType = "int"
	FieldFloat   FieldType = "float"
	FieldString  FieldType = "string"
	FieldBool    FieldType = "bool"
	FieldObject  FieldType = "object" // better not to use it outside the mongodb
)
