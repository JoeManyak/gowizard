package gentags

import (
	"fmt"
	"gowizard/builder/model/system"
	"gowizard/util"
	"strings"
)

type Postgres struct {
	modelInstance *system.Model
	layer         *system.Layer
}

func (c *Postgres) Create() string {
	if c.layer.NextLayer == nil {
		return defaultError
	}

	return fmt.Sprintf("return %s.%s.Create%s(%s)\n",
		strings.ToLower(string([]rune(c.modelInstance.Name)[0])),
		c.layer.NextLayer.Name,
		c.modelInstance.Name,
		util.MakePrivateName(c.modelInstance.Name)+"Model")
}

func (c *Postgres) Read() string {
	if c.layer.NextLayer == nil {
		return defaultError
	}

	return fmt.Sprintf("return %s.%s.Read%s(%s)\n",
		strings.ToLower(string([]rune(c.modelInstance.Name)[0])),
		c.layer.NextLayer.Name,
		c.modelInstance.Name,
		util.MakePrivateName(c.modelInstance.Name)+"Model")
}

func (c *Postgres) Update() string {
	if c.layer.NextLayer == nil {
		return defaultError
	}

	return fmt.Sprintf("return %s.%s.Update%s(%s)\n",
		strings.ToLower(string([]rune(c.modelInstance.Name)[0])),
		c.layer.NextLayer.Name,
		c.modelInstance.Name,
		util.MakePrivateName(c.modelInstance.Name)+"Model")
}

func (c *Postgres) Delete() string {
	if c.layer.NextLayer == nil {
		return defaultError
	}

	return fmt.Sprintf("return %s.%s.Delete%s(%s)\n",
		strings.ToLower(string([]rune(c.modelInstance.Name)[0])),
		c.layer.NextLayer.Name,
		c.modelInstance.Name,
		util.MakePrivateName(c.modelInstance.Name)+"Model")
}

func (c *Postgres) Custom() string {
	return defaultCustom
}
