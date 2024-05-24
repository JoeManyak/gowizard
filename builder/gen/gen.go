package gen

import (
	"encoding/json"
	"fmt"
	"gowizard/builder/model"
	"gowizard/builder/model/system"
	"gowizard/consts"
	"gowizard/util"
	"os"
	"strings"
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
			if j < len(methods[i].Args)-2 && j%2 == 1 {
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
	if _, err := g.File.WriteString("}\n\n"); err != nil {
		return err
	}

	return nil
}

func (g *Gen) AddStruct(model *system.Model) error {
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

	if _, err := g.File.WriteString("}\n\n"); err != nil {
		return err
	}

	return nil
}

func (g *Gen) AddMethod(mdl *system.Model, method *model.MethodInstance) error {
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

	method.Returns = method.GetReturns()
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
			space := ""
			if j < len(method.Returns)-1 {
				space = ", "
			}
			if _, err := g.File.WriteString(method.Returns[j] + space); err != nil {
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

	if _, err := g.File.WriteString("}\n\n"); err != nil {
		return err
	}

	return nil
}

func (g *Gen) AddMainRouterFunc(mdls []*system.Model) error {
	//todo should me method, not function, in router struct initialized controller lol
	_, err := g.File.WriteString(`func Run() {
gin.New()
//todo implement here
`)

	_, err = g.File.WriteString(`}`)
	return err
}

func (g *Gen) AddSubRouterFunc(mdl *system.Model) error {
	_, err := g.File.WriteString(fmt.Sprintf("func %sRouter(r *gin.RouterGroup) {\n", mdl.Name))
	if err != nil {
		return err
	}

	for i := range mdl.Methods {
		_, err = g.File.WriteString(fmt.Sprintf("r.%s(%s, ", mdl.Name))
		if err != nil {
			return err
		}
	}
}

func (g *Gen) AddParseConfigMethod(mdl *system.Model) error {
	_, err := g.File.WriteString("func New" + util.MakePublicName(consts.DefaultConfigFolder) + "() (*" + mdl.Name + ", error) {\n")
	if err != nil {
		return err
	}

	_, err = g.File.WriteString(fmt.Sprintf(`f, err := os.Open("%s.json")
if err != nil {
	return nil, err
}
defer f.Close()

var c %s
bytes, err := io.ReadAll(f)
if err != nil {
	return nil, err
}

err = json.Unmarshal(bytes, &c)
if err != nil {
	return nil, err
}

return &c, nil
}`, consts.DefaultConfigFolder, util.MakePublicName(consts.DefaultConfigFolder)))
	if err != nil {
		return err
	}

	return err
}

func (g *Gen) AddImport(imports []string) error {
	_, err := g.File.WriteString(genImports(imports))
	return err
}

func (g *Gen) WriteJSON(mdl *system.Model) error {
	defaults := getDefaultConfigValues()
	var data = make(map[string]string, 10)
	for i := range mdl.Fields {
		data[util.PascalToSnakeCase(mdl.Fields[i].Name)] =
			defaults[mdl.Fields[i].Name]
	}

	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	_, err = g.File.Write(b)
	return err
}

func genImports(imports []string) string {
	switch len(imports) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("import %s\n", imports[0])
	default:
		var sb strings.Builder
		sb.WriteString("import (\n")
		for _, v := range imports {
			sb.WriteString(v + "\n")
		}
		sb.WriteString(")\n\n")
		return sb.String()
	}
}

func getDefaultConfigValues() map[string]string {
	return map[string]string{
		"HttpHost": "localhost",
		"HttpPort": "8080",

		"PostgresHost":     "localhost",
		"PostgresPort":     "5432",
		"PostgresUser":     "postgres",
		"PostgresPassword": "postgres",
	}
}
