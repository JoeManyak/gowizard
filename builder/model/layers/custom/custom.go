package custom

import (
	"fmt"
	"gowizard/builder/model"
	"gowizard/builder/storage"
)

type Custom struct {
	modelInstance *model.Model
	layer         *storage.Layer
}

func NewCustom(layer *storage.Layer, modelInstance *model.Model) *Custom {
	return &Custom{
		layer:         layer,
		modelInstance: modelInstance,
	}
}

func (c *Custom) Create() string {
	if c.layer.NextLayer == nil {
		return "panic(\"implement me\")"
	}
	return fmt.Sprintf("return %s.%s.Create%s(/*should be dto*/)\n",
		c.modelInstance.GetPointerName(),
		c.layer.NextLayer.Name,
		c.modelInstance.Name,
	)
}
