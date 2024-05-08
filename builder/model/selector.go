package model

type MethodBodyGenerator interface {
	Create() string
	// todo do for all methods
}

type BaseMethodBodyGenerator struct {
	Model         *Model
	NextLayerName string
}

func (b *BaseMethodBodyGenerator) Create() string {
	return ""
}
