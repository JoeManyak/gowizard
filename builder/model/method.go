package model

import (
	"fmt"
	"gowizard/builder/model/gentags"
	"gowizard/builder/model/system"
	"gowizard/consts"
)

type InterfaceMethodInstance struct {
	Name string
	// Args should be a pairs of type and name
	Args    []string
	Returns []string `json:"returns"`

	Type  system.MethodType
	Model *system.Model
}

type MethodInstance struct {
	Name string
	// Args should be a pairs of type and name
	Args    []string
	Returns []string

	Layer *system.Layer
	Type  system.MethodType
	Model *system.Model
}

func (mi *MethodInstance) UpdateByMethodType() {
	// todo: should consider to change to GenerateNaming()
	switch mi.Type.Lower() {
	case system.MethodCreate:
		mi.Name = "Create" + mi.Model.Name
	case system.MethodRead:
		mi.Name = "Read" + mi.Model.Name
	case system.MethodUpdate:
		mi.Name = "Update" + mi.Model.Name
	case system.MethodDelete:
		mi.Name = "Delete" + mi.Model.Name
	default:
		mi.Name = mi.Type.String() + mi.Model.Name
	}
}

func (mi *MethodInstance) GetMethodBody() string {
	selector := SelectMethods(mi.Model, mi.Layer)

	switch mi.Type.Lower() {
	case system.MethodCreate:
		return selector.Create()
	case system.MethodRead:
		return selector.Read()
	case system.MethodUpdate:
		return selector.Update()
	case system.MethodDelete:
		return selector.Delete()
	default:
		return selector.Custom()
	}
}

func (mi *MethodInstance) GetReturns() []string {
	if len(mi.Returns) != 0 {
		return mi.Returns
	}

	return mi.Type.GetDefaultReturns(mi.Model)
}

func SelectMethods(mdl *system.Model, layer *system.Layer) GenerateMethodBody {
	if layer == nil {
		fmt.Println("Warn: no layer provided")
		return gentags.NewCustom(layer, mdl)
	}

	switch layer.Type {
	case consts.HTTPLayerType:
		return gentags.NewHTTP(layer, mdl)
	/*case repoLayerType:
	return &RepoMethodBody{}*/
	default:
		return gentags.NewCustom(layer, mdl)
	}
}

type GenerateMethodBody interface {
	Create() string
	Read() string
	Update() string
	Delete() string
	Custom() string
}

var _ GenerateMethodBody = &gentags.Custom{}
var _ GenerateMethodBody = &gentags.HTTP{}
