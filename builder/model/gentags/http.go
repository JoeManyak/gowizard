package gentags

import (
	"fmt"
	"gowizard/builder/model/system"
	"gowizard/consts"
	"strings"
)

type HTTP struct {
	modelInstance *system.Model
	layer         *system.Layer
}

const (
	baseBody = `var req %s
err := ctx.ShouldBindBodyWithJSON(&req)
if err != nil {
	ctx.JSON(422, gin.H{"error": err.Error()})
	return
}

res, err := %s.%s.%s%s(&req)
if err != nil {
	ctx.JSON(500, gin.H{"error": err.Error()})
	return
}

ctx.JSON(200, gin.H{"data": res})
`
)

func NewHTTP(layer *system.Layer, modelInstance *system.Model) *HTTP {
	return &HTTP{
		layer:         layer,
		modelInstance: modelInstance,
	}
}

func (h *HTTP) Create() string {
	modelType := fmt.Sprintf("%s.%s", consts.DefaultModelsFolder, h.modelInstance.Name)
	return fmt.Sprintf(baseBody,
		modelType,
		strings.ToLower(string([]rune(h.modelInstance.Name)[0])),
		h.layer.NextLayer.Name,
		"Create",
		h.modelInstance.Name,
	)
}

func (h *HTTP) Read() string {
	modelType := fmt.Sprintf("%s.%s", consts.DefaultModelsFolder, h.modelInstance.Name)
	return fmt.Sprintf(baseBody,
		modelType,
		strings.ToLower(string([]rune(h.modelInstance.Name)[0])),
		h.layer.NextLayer.Name,
		"Read",
		h.modelInstance.Name,
	)
}

func (h *HTTP) Update() string {
	modelType := fmt.Sprintf("%s.%s", consts.DefaultModelsFolder, h.modelInstance.Name)
	return fmt.Sprintf(baseBody,
		modelType,
		strings.ToLower(string([]rune(h.modelInstance.Name)[0])),
		h.layer.NextLayer.Name,
		"Update",
		h.modelInstance.Name,
	)
}

func (h *HTTP) Delete() string {
	modelType := fmt.Sprintf("%s.%s", consts.DefaultModelsFolder, h.modelInstance.Name)
	return fmt.Sprintf(baseBody,
		modelType,
		strings.ToLower(string([]rune(h.modelInstance.Name)[0])),
		h.layer.NextLayer.Name,
		"Delete",
		h.modelInstance.Name,
	)
}

func (h *HTTP) Custom() string {
	return defaultCustom
}
