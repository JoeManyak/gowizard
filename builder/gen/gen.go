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
		case 1:
			if _, err := g.File.WriteString(" " + methods[i].Returns[0] + "\n"); err != nil {
				return err
			}
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

func (g *Gen) NewLayerFunc(layer *system.Layer, mdl *system.Model) error {
	args := fmt.Sprintf("%s *%s.%s", consts.DefaultConfigFolder, consts.DefaultConfigFolder, util.MakePublicName(consts.DefaultConfigFolder))
	if layer.NextLayer != nil {
		args = fmt.Sprintf("%s, %s %s.%s", args, layer.NextLayer.Name, layer.NextLayer.Name, mdl.Name)
	}

	if layer.Type == consts.RepoLayerType {
		args = fmt.Sprintf("%s, db *gorm.DB", args)
	}
	if layer.Type == consts.TelebotLayerType {
		args = fmt.Sprintf("%s, bot *telebot.Bot", args)
	}

	_, err := g.File.WriteString(fmt.Sprintf("func New%s%s(%s) %s {\nreturn &%s{\n",
		mdl.Name, util.MakePublicName(layer.Name), args, mdl.Name, util.MakePrivateName(mdl.Name)))
	if err != nil {
		return err
	}

	_, err = g.File.WriteString(fmt.Sprintf("%s: %s,\n", consts.DefaultConfigFolder, consts.DefaultConfigFolder))
	if err != nil {
		return err
	}

	if layer.NextLayer != nil {
		_, err = g.File.WriteString(fmt.Sprintf("%s: %s,\n", layer.NextLayer.Name, layer.NextLayer.Name))
		if err != nil {
			return err
		}
	}
	if layer.Type == consts.RepoLayerType {
		_, err = g.File.WriteString("db: db,\n")
		if err != nil {
			return err
		}
	}
	if layer.Type == consts.TelebotLayerType {
		_, err = g.File.WriteString("bot: bot,\n")
		if err != nil {
			return err
		}
	}

	_, err = g.File.WriteString("}\n}\n\n")
	return err
}

func (g *Gen) AddMethodWithSwagger(mdl *system.Model, method *model.MethodInstance) error {
	publicName := util.MakePublicName(mdl.Name)
	route := method.Type.GetRoute()
	if route != "" {
		route = "/" + route
	}
	_, err := g.File.WriteString(fmt.Sprintf(
		`// @Summary %s %s
// @Tags %s
// @Accept json
// @Produce json
// @Param message body %s.%s true "%s"
// @Success 200 {object} %s.%s
// @Router /%s%s [%s]
`, method.Type.String(), publicName, publicName, consts.DefaultModelsFolder, publicName, mdl.Name, consts.DefaultModelsFolder, publicName, mdl.GetLayer().Name, route, strings.ToLower(method.Type.GetHTTPType())))
	if err != nil {
		return err
	}

	return g.AddMethod(mdl, method)
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
	case 1:
		if _, err := g.File.WriteString(" " + method.Returns[0] + "{\n"); err != nil {
			return err
		}
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

func (g *Gen) AddMainTelerouterNewFunc(mdl *system.Model) error {
	args := fmt.Sprintf("%s %s", mdl.Fields[0].Name, mdl.Fields[0].Type)
	for _, f := range mdl.Fields[1:] {
		args = fmt.Sprintf("%s, %s %s", args, f.Name, f.Type)
	}

	_, err := g.File.WriteString(fmt.Sprintf("func NewTeleRouter(%s) *TeleRouter {\nreturn &TeleRouter{\n", args))
	if err != nil {
		return err
	}

	for _, f := range mdl.Fields {
		_, err = g.File.WriteString(fmt.Sprintf("%s: %s,\n", f.Name, f.Name))
		if err != nil {
			return err
		}
	}

	_, err = g.File.WriteString("}\n}\n\n")
	if err != nil {
		return err
	}
	return nil
}

func (g *Gen) AddMainRouterNewFunc(mdl *system.Model) error {
	args := fmt.Sprintf("%s %s", mdl.Fields[0].Name, mdl.Fields[0].Type)
	for _, f := range mdl.Fields[1:] {
		args = fmt.Sprintf("%s, %s %s", args, f.Name, f.Type)
	}

	_, err := g.File.WriteString(fmt.Sprintf("func NewRouter(%s) *Router {\nreturn &Router{\n", args))
	if err != nil {
		return err
	}

	for _, f := range mdl.Fields {
		_, err = g.File.WriteString(fmt.Sprintf("%s: %s,\n", f.Name, f.Name))
		if err != nil {
			return err
		}
	}

	_, err = g.File.WriteString("}\n}\n\n")
	if err != nil {
		return err
	}
	return nil
}

func (g *Gen) AddMainTeleRouterFunc(mdls []*system.Model) error {
	_, err := g.File.WriteString(`func (r *TeleRouter) Run() {
`)
	if err != nil {
		return err
	}

	for i := range mdls {
		_, err = g.File.WriteString(fmt.Sprintf("// Generated router for %s use cases\n", mdls[i].Name))
		if err != nil {
			return err
		}

		for _, method := range mdls[i].Methods {
			_, err = g.File.WriteString(fmt.Sprintf("r.Bot.Handle(\"/%s%s\", r.%s.%s)\n",
				strings.ToLower(mdls[i].Name),
				strings.ToLower(string(method)),
				mdls[i].Name,
				method.GenerateNaming(mdls[i].Name),
			))
			if err != nil {
				return err
			}
		}

		_, err = g.File.WriteString("\n")
		if err != nil {
			return err
		}
	}

	_, err = g.File.WriteString(`r.Bot.Start()
}`)
	return err
}

func (g *Gen) AddMainRouterFunc(mdls []*system.Model) error {
	configFile := util.MakePublicName(consts.DefaultConfigFolder)

	_, err := g.File.WriteString(`func (r *Router) Run() {
g := gin.New()

g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
`)
	if err != nil {
		return err
	}

	for i := range mdls {
		routerName := fmt.Sprintf("%sRouter", mdls[i].Name)
		_, err = g.File.WriteString(fmt.Sprintf("// Generated router for %s use cases\n%s := g.Group(\"/%s\")\n", mdls[i].Name, routerName, mdls[i].GetLayer().Name))
		if err != nil {
			return err
		}

		for _, method := range mdls[i].Methods {
			_, err = g.File.WriteString(fmt.Sprintf("%s.%s(\"/%s\", r.%s.%s%s)\n", routerName, method.GetHTTPType(), method.GetRoute(), mdls[i].Name, method.String(), mdls[i].Name))
			if err != nil {
				return err
			}
		}

		_, err = g.File.WriteString("\n")
		if err != nil {
			return err
		}
	}

	_, err = g.File.WriteString(fmt.Sprintf(`err := g.Run(r.%s.HttpHost + ":" + r.%s.HttpPort)
if err != nil {
panic(err.Error())
}
}
`, configFile, configFile))
	return err
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
		"HttpHost": "",
		"HttpPort": "8080",

		"PostgresHost":     "localhost",
		"PostgresPort":     "5432",
		"PostgresDb":       "default",
		"PostgresUser":     "postgres",
		"PostgresPassword": "postgres",

		"TelebotToken": "put-your-token-here",
	}
}
