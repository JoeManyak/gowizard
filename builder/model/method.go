package model

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

	LayerType string
	Type      MethodType
	Model     *Model
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

func (mi *MethodInstance) GetMethodBody() string {
	switch mi.Type {
	case MethodCreate:
		return "return nil\n"
	default:
		return "\n"
	}
}

func (mi *MethodInstance) getCreateMethodBody() string {
	return "//todo implement me\n"
}

func (mi *MethodInstance) getDefaultMethodBody() string {
	return "panic(\"not implemented\")\n"
}
