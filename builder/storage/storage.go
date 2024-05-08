package storage

import "gowizard/builder/model"

type Layer struct {
	Name      string
	Type      string
	Path      string
	Models    *[]*model.Model
	NextLayer *Layer
}
