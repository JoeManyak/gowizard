package gen

import (
	"fmt"
	"gowizard/builder/model"
	"gowizard/util"
	"os"
)

type Gen struct {
	File *os.File
}

func NewGen(path string) (*Gen, error) {
	var newGen = Gen{}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	newGen.File = file
	return &newGen, nil
}

func (g *Gen) Close() error {
	return g.File.Close()
}

func (g *Gen) AddPackage(packageName string) error {
	_, err := g.File.WriteString("package " + packageName + "\n\n")
	return err
}

func (g *Gen) AddInterface(name string, methods []model.InterfaceMethodInstance) error {
	if _, err := g.File.WriteString("type " + name + " interface {\n"); err != nil {
		return err
	}

	for i := range methods {
		if _, err := g.File.WriteString(methods[i].Name + "("); err != nil {
			return err
		}

		for j := 0; j < len(methods[i].Args)-1; j++ {
			comma := ""
			if j < len(methods[i].Args)-2 {
				comma = ", "
			}

			if _, err := g.File.WriteString(methods[i].Args[j] + " " + methods[i].Args[j+1] + comma); err != nil {
				return err
			}
		}

		if _, err := g.File.WriteString(")"); err != nil {
			return err
		}

		switch len(methods[i].Returns) {
		case 0:
			if _, err := g.File.WriteString("\n"); err != nil {
				return err
			}
			break
		case 1:
			if _, err := g.File.WriteString(" " + methods[i].Returns[0] + "\n"); err != nil {
				return err
			}
			break
		default:
			if _, err := g.File.WriteString(" ("); err != nil {
				return err
			}

			for j := range methods[i].Returns {
				comma := ""
				if j < len(methods[i].Returns)-1 {
					comma = ", "
				}
				if _, err := g.File.WriteString(methods[i].Returns[j] + comma); err != nil {
					return err
				}
			}

			if _, err := g.File.WriteString(")\n"); err != nil {
				return err
			}
		}
	}
	if _, err := g.File.WriteString("}\n"); err != nil {
		return err
	}

	return nil
}

func (g *Gen) AddStruct(model *model.Model) error {
	if _, err := g.File.WriteString("type " + model.Name + " struct {\n"); err != nil {
		return err
	}

	for i := range model.Fields {
		if _, err := g.File.WriteString(
			fmt.Sprintf("%s %s `json:\"%s\"`\n",
				model.Fields[i].Name,
				string(model.Fields[i].Type),
				util.PascalToSnakeCase(model.Fields[i].Name))); err != nil {
			return err
		}
	}

	if _, err := g.File.WriteString("}\n"); err != nil {
		return err
	}

	return nil
}

func (g *Gen) AddMethod(mdl *model.Model, method *model.MethodInstance) error {
	_, err := g.File.WriteString("func (" + mdl.GetPointerName() + ") " + method.Name + "(")
	if err != nil {
		return err
	}

	for j := 0; j < len(method.Args)-1; j++ {
		comma := ""
		if j < len(method.Args)-2 {
			comma = ", "
		}

		if _, err := g.File.WriteString(method.Args[j] + " " + method.Args[j+1] + comma); err != nil {
			return err
		}
	}

	if _, err := g.File.WriteString(")"); err != nil {
		return err
	}

	switch len(method.Returns) {
	case 0:
		if _, err := g.File.WriteString("{ \n"); err != nil {
			return err
		}
		break
	case 1:
		if _, err := g.File.WriteString(" " + method.Returns[0] + "{\n"); err != nil {
			return err
		}
		break
	default:
		if _, err := g.File.WriteString(" ("); err != nil {
			return err
		}

		for j := range method.Returns {
			comma := ""
			if j < len(method.Returns)-1 {
				comma = ", "
			}
			if _, err := g.File.WriteString(method.Returns[j] + comma); err != nil {
				return err
			}
		}

		if _, err := g.File.WriteString(") {\n"); err != nil {
			return err
		}
	}

	_, err = g.File.WriteString(method.GetMethodBody())
	if err != nil {
		return err
	}

	if _, err := g.File.WriteString("}\n"); err != nil {
		return err
	}

	return nil
}
