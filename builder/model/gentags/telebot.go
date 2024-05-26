package gentags

import (
	"fmt"
	"gowizard/builder/model/system"
	"gowizard/consts"
	"gowizard/util"
	"strings"
)

type Telebot struct {
	modelInstance *system.Model
	layer         *system.Layer
}

const (
	baseTelebotBody = `args := strings.Split(m.Payload, " ")
if len(args) > 0 {
	var req %s.%s
	body := strings.Join(args[1:], " ")
	json.Unmarshal([]byte(body), &req)
	res, err := %s.%s.%s%s(&req)
if err != nil {
	%s.bot.Send(m.Sender, err.Error())
}

b, _ := json.Marshal(res)
%s.bot.Send(m.Sender, string(b)) 
} else {
	%s.bot.Send(m.Sender, "Unable to %s %s")
}
`

	deleteTelebotBody = `args := strings.Split(m.Payload, " ")
if len(args) > 0 {
	var req %s.%s
	body := strings.Join(args[1:], " ")
	json.Unmarshal([]byte(body), &req)
	err := %s.%s.%s%s(&req)
if err != nil {
	%s.bot.Send(m.Sender, err.Error())
}

%s.bot.Send(m.Sender, "Success") 
} else {
	%s.bot.Send(m.Sender, "Unable to %s %s")
}
`
)

func NewTelebot(layer *system.Layer, modelInstance *system.Model) *Telebot {
	return &Telebot{
		layer:         layer,
		modelInstance: modelInstance,
	}
}

func (t *Telebot) Create() string {
	return fmt.Sprintf(baseTelebotBody,
		consts.DefaultModelsFolder,
		t.modelInstance.Name,
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		t.layer.NextLayer.Name,
		"Create",
		t.modelInstance.Name,
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		"create",
		util.MakePrivateName(t.modelInstance.Name),
	)
}

func (t *Telebot) Read() string {
	return fmt.Sprintf(baseTelebotBody,
		consts.DefaultModelsFolder,
		t.modelInstance.Name,
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		t.layer.NextLayer.Name,
		"Read",
		t.modelInstance.Name,
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		"read",
		util.MakePrivateName(t.modelInstance.Name),
	)
}

func (t *Telebot) Update() string {
	return fmt.Sprintf(baseTelebotBody,
		consts.DefaultModelsFolder,
		t.modelInstance.Name,
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		t.layer.NextLayer.Name,
		"Update",
		t.modelInstance.Name,
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		"update",
		util.MakePrivateName(t.modelInstance.Name),
	)
}

func (t *Telebot) Delete() string {
	return fmt.Sprintf(deleteTelebotBody,
		consts.DefaultModelsFolder,
		t.modelInstance.Name,
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		t.layer.NextLayer.Name,
		"Delete",
		t.modelInstance.Name,
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		strings.ToLower(string([]rune(t.modelInstance.Name)[0])),
		"delete",
		util.MakePrivateName(t.modelInstance.Name),
	)
}

func (t *Telebot) Custom() string {
	return defaultCustom
}
