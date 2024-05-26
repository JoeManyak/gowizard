package system

import (
	"gowizard/consts"
	"gowizard/util"
	"net/http"
	"strings"
)

type Model struct {
	Name    string       `yaml:"name"`
	Fields  []Field      `yaml:"fields"`
	Methods []MethodType `yaml:"methods"`
}

type MethodType string

func (mt MethodType) String() string {
	if parsed, ok := DefaultMethodNamings[mt.Lower()]; ok {
		return parsed
	}

	// if not found in default namings, just return the method name with first letter capitalized
	runes := []rune(string(mt))
	naming := strings.ToUpper(string(runes[0])) + string(runes[1:])

	return naming
}

func (mt MethodType) Lower() MethodType {
	return MethodType(strings.ToLower(string(mt)))
}

func (mt MethodType) GenerateNaming(methodModelName string) string {
	switch mt.Lower() {
	case MethodCreate:
		return "Create" + methodModelName
	case MethodRead:
		return "Read" + methodModelName
	case MethodUpdate:
		return "Update" + methodModelName
	case MethodDelete:
		return "Delete" + methodModelName
	default:
		return mt.String() + methodModelName
	}
}

func (mt MethodType) GetHTTPType() string {
	switch mt.Lower() {
	case MethodCreate:
		return http.MethodPost
	case MethodRead:
		return http.MethodGet
	case MethodUpdate:
		return http.MethodPatch
	case MethodDelete:
		return http.MethodDelete
	default:
		return http.MethodPost
	}
}

func (mt MethodType) GetRoute() string {
	switch mt.Lower() {
	case MethodCreate, MethodRead, MethodUpdate, MethodDelete:
		return ""
	default:
		return string(mt.Lower())
	}
}

const (
	MethodCreate MethodType = "create"
	MethodRead   MethodType = "read"
	MethodUpdate MethodType = "update"
	MethodDelete MethodType = "delete"
)

func (mt MethodType) GetDefaultReturns(mdl *Model) []string {
	switch MethodType(strings.ToLower(string(mt))) {
	case MethodRead:
		return []string{
			"[]" + consts.DefaultModelsFolder + "." + mdl.Name,
			"error"}
	case MethodDelete:
		return []string{"error"}
	default:
		return []string{
			"*" + consts.DefaultModelsFolder + "." + mdl.Name,
			"error"}
	}
}

func (mt MethodType) GetDefaultArgs(mdl *Model, layer *Layer) []string {
	if layer.Type == consts.HTTPLayerType {
		return []string{"ctx", "*gin.Context"}
	}

	return []string{util.MakePrivateName(mdl.Name + "Model"), " *" + consts.DefaultModelsFolder + "." + mdl.Name}
}

var DefaultMethodNamings = map[MethodType]string{
	MethodCreate: "Create",
	MethodRead:   "Read",
	MethodUpdate: "Update",
	MethodDelete: "Delete",
}

type Field struct {
	Name string    `yaml:"name"`
	Type FieldType `yaml:"type"`
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

type Layer struct {
	Name      string
	Type      string
	Path      string
	Models    *[]*Model
	NextLayer *Layer
}
