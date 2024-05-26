package gentags

import (
	"fmt"
	"gowizard/builder/model/system"
	"gowizard/consts"
	"gowizard/util"
	"strings"
)

type Postgres struct {
	modelInstance *system.Model
	layer         *system.Layer
}

func NewPostgres(layer *system.Layer, modelInstance *system.Model) *Postgres {
	return &Postgres{
		layer:         layer,
		modelInstance: modelInstance,
	}
}

func (p *Postgres) Create() string {
	return fmt.Sprintf(`result := %s.db.Create(%sModel)
return %sModel, result.Error
`, strings.ToLower(string([]rune(p.modelInstance.Name)[0])), util.MakePrivateName(p.modelInstance.Name), util.MakePrivateName(p.modelInstance.Name))
}

func (p *Postgres) Read() string {
	return fmt.Sprintf(`var %sModelList []%s.%s
result := %s.db.Where(%sModel).Find(&%sModelList)
return %sModelList, result.Error
`, util.MakePrivateName(p.modelInstance.Name),
		consts.DefaultModelsFolder,
		p.modelInstance.Name,
		strings.ToLower(string([]rune(p.modelInstance.Name)[0])),
		util.MakePrivateName(p.modelInstance.Name),
		util.MakePrivateName(p.modelInstance.Name),
		util.MakePrivateName(p.modelInstance.Name))
}

func (p *Postgres) Update() string {
	return fmt.Sprintf(`result := %s.db.Save(%sModel)
return %sModel, result.Error
`, strings.ToLower(string([]rune(p.modelInstance.Name)[0])), util.MakePrivateName(p.modelInstance.Name), util.MakePrivateName(p.modelInstance.Name))
}

func (p *Postgres) Delete() string {
	return fmt.Sprintf(`result := %s.db.Delete(%sModel)
return result.Error
`, strings.ToLower(string([]rune(p.modelInstance.Name)[0])), util.MakePrivateName(p.modelInstance.Name))
}

func (p *Postgres) Custom() string {
	return defaultCustom
}
