package model

import (
	"gowizard/util"
	"strings"
)

type Model struct {
	Name    string       `yaml:"name"`
	Fields  []Field      `yaml:"fields"`
	Methods []MethodType `yaml:"methods"`
}

type MethodType string

func (m MethodType) String() string {
	if parsed, ok := DefaultMethodNamings[m.Lower()]; ok {
		return parsed
	}

	// if not found in default namings, just return the method name with first letter capitalized
	runes := []rune(string(m))
	naming := strings.ToUpper(string(runes[0])) + string(runes[1:])

	return naming
}

func (m MethodType) Lower() MethodType {
	return MethodType(strings.ToLower(string(m)))
}

func (m MethodType) GenerateNaming(methodModelName string) string {
	switch m.Lower() {
	case MethodCreate:
		return "Create" + methodModelName
	case MethodRead:
		return "Read" + methodModelName
	case MethodUpdate:
		return "Update" + methodModelName
	case MethodDelete:
		return "Delete" + methodModelName
	default:
		return m.String() + methodModelName
	}
}

const (
	MethodCreate MethodType = "create"
	MethodRead   MethodType = "read"
	MethodUpdate MethodType = "update"
	MethodDelete MethodType = "delete"
)

var DefaultMethodNamings = map[MethodType]string{
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
	FieldFloat   FieldType = "float64"
	FieldString  FieldType = "string"
	FieldBool    FieldType = "bool"
	FieldObject  FieldType = "object" // better not to use it outside the mongodb
)

func (m *Model) GetFilename() string {
	return util.PascalToSnakeCase(m.Name) + ".go"
}

func (m *Model) GetPointerName() string {
	return strings.ToLower(string([]rune(m.Name)[0])) + " *" + m.Name
}

func (m *Model) GetLayer() Model {
	return Model{
		Name: util.MakePrivateName(m.Name),
	}
}
