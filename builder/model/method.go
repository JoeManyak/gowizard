package model

import (
	"gowizard/builder/model/layers/custom"
	"gowizard/builder/storage"
)

type InterfaceMethodInstance struct {
	Name string
	// Args should be a pairs of type and name
	Args    []string
	Returns []string `json:"returns"`

	Type  MethodType
	Model *Model
}

type MethodInstance struct {
	Name string
	// Args should be a pairs of type and name
	Args    []string
	Returns []string

	Layer *storage.Layer
	Type  MethodType
	Model *Model
}

func (mi *MethodInstance) UpdateByMethodType() {
	// todo: should consider to change to GenerateNaming()
	switch mi.Type.Lower() {
	case MethodCreate:
		mi.Name = "Create" + mi.Model.Name
	case MethodRead:
		mi.Name = "Read" + mi.Model.Name
	case MethodUpdate:
		mi.Name = "Update" + mi.Model.Name
	case MethodDelete:
		mi.Name = "Delete" + mi.Model.Name
	default:
		mi.Name = mi.Type.String() + mi.Model.Name
	}
}

const (
	httpLayerType = "http"
	repoLayerType = "postgres"
)

func (mi *MethodInstance) GetMethodBody() string {
	selector := SelectMethods(mi.Model, mi.Layer)

	switch mi.Type {
	case MethodCreate:
		return selector.Create()
	default:
		return "panic(\"not implemented\")\n"
	}
}

func SelectMethods(mdl *Model, layer *storage.Layer) GenerateMethodBody {
	switch layer.Type {
	/*case httpLayerType:
		return &HttpMethodBody{}
	case repoLayerType:
		return &RepoMethodBody{}*/
	default:
		return custom.NewCustom(layer, mdl)
	}
}

type GenerateMethodBody interface {
	Create() string
	//Read() string
	//Update() string
	//Delete() string
	//Custom() string
}
