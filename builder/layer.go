package builder

import (
	"fmt"
	"gowizard/builder/gen"
	"gowizard/builder/model"
	"gowizard/builder/model/system"
	"gowizard/consts"
	"gowizard/util"
	"path/filepath"
	"strings"
)

type LayerController struct {
	Layers  []*system.Layer
	Builder *Builder
	Models  []*system.Model
}

func NewLayerController(
	b *Builder,
	layers []LayerDTO,
	models []*system.Model,
) *LayerController {
	lc := LayerController{
		Builder: b,
		Models:  models,
		Layers:  make([]*system.Layer, 0, len(layers)),
	}

	for _, l := range layers {
		lc.Layers = append(lc.Layers, &system.Layer{
			Name:   l.Layer,
			Type:   l.Tag,
			Models: &models,
		})
	}

	for i := 0; i < len(lc.Layers)-1; i++ {
		lc.Layers[i].NextLayer = lc.Layers[i+1]
	}

	return &lc
}

func (lc *LayerController) Generate() error {
	layerTypes := make(map[string]*system.Layer)
	for _, layer := range lc.Layers {
		if layer.Type != "" {
			layerTypes[layer.Type] = layer
		}
		// generate general file
		err := lc.generateMainLayerFile(layer)
		if err != nil {
			return err
		}
		for i := range *layer.Models {
			// generate model file
			err = lc.generateModelLayerFile(layer, (*layer.Models)[i])
			if err != nil {
				return err
			}
		}
	}

	err := lc.generateModelStorageFile()
	if err != nil {
		return err
	}

	err = lc.generateConfigStorageFile()
	if err != nil {
		return err
	}

	if l, ok := layerTypes[consts.HTTPLayerType]; ok {
		err = lc.generateRouter(l)
		if err != nil {
			return err
		}
	}

	err = lc.generateMainFile()
	if err != nil {
		return err
	}

	return nil
}

func (lc *LayerController) generateMainFile() error {
	g, err := gen.NewGen(filepath.Join(lc.Builder.Path, "main.go"))
	if err != nil {
		return fmt.Errorf("unable to create new main generator: %w", err)
	}

	err = g.AddPackage("main")
	if err != nil {
		return fmt.Errorf("unable to add package main: %w", err)
	}

	var httpLayer *system.Layer
	for i := range lc.Layers {
		if lc.Layers[i].Type == consts.HTTPLayerType {
			httpLayer = lc.Layers[i]
		}
	}

	importsToAdd := make([]string, 0, len(lc.Layers)+3)
	importsToAdd = append(importsToAdd,
		util.MakeString(filepath.Join(lc.Builder.ProjectName, consts.DefaultConfigFolder)),
	)
	if httpLayer != nil {
		importsToAdd = append(importsToAdd,
			util.MakeString(filepath.Join(lc.Builder.ProjectName, consts.DefaultRouterFolder)),
		)
	}
	for i := range lc.Layers {
		importsToAdd = append(importsToAdd, util.MakeString(filepath.Join(lc.Builder.ProjectName, lc.Layers[i].Name)))
	}

	err = g.AddImport(importsToAdd)
	if err != nil {
		return fmt.Errorf("unable to add imports: %w", err)
	}

	_, err = g.File.WriteString("func main() {\n")
	if err != nil {
		return err
	}

	_, err = g.File.WriteString(fmt.Sprintf("%s, err := %s.New%s()\nif err != nil {\npanic(err.Error())\n}\n\n",
		consts.DefaultConfigFolder, consts.DefaultConfigFolder, util.MakePublicName(consts.DefaultConfigFolder)))
	if err != nil {
		return err
	}

	for i := len(lc.Layers) - 1; i >= 0; i-- {
		for _, mdl := range lc.Models {
			args := "config"
			if i < len(lc.Layers)-1 {
				args = fmt.Sprintf("%s, %s%s", args, mdl.Name, util.MakePublicName(lc.Layers[i+1].Name))
			}

			_, err = g.File.WriteString(fmt.Sprintf("%s%s := %s.New%s%s(%s)\n", mdl.Name, util.MakePublicName(lc.Layers[i].Name), lc.Layers[i].Name, mdl.Name, util.MakePublicName(lc.Layers[i].Name), args))
		}

		_, err = g.File.WriteString("\n")
		if err != nil {
			return err
		}
	}

	if httpLayer != nil {
		args := lc.Models[0].Name + util.MakePublicName(httpLayer.Name)
		for i := range lc.Models[1:] {
			args += ", " + lc.Models[i+1].Name + util.MakePublicName(httpLayer.Name)
		}
		args += ", " + consts.DefaultConfigFolder

		_, err = g.File.WriteString(fmt.Sprintf("r := %s.New%s(%s)\n", consts.DefaultRouterFolder, util.MakePublicName(consts.DefaultRouterFolder), args))
		if err != nil {
			return err
		}

		_, err = g.File.WriteString("r.Run()")
		if err != nil {
			return err
		}
	}

	_, err = g.File.WriteString("}\n")
	if err != nil {
		return err
	}

	return nil
}

func (lc *LayerController) generateMainLayerFile(layer *system.Layer) error {
	g, err := gen.NewGen(layer.Path + layer.Name + ".go")
	if err != nil {
		return fmt.Errorf("unable to create new generator: %w", err)
	}

	defer g.Close()

	err = g.AddPackage(layer.Name)
	if err != nil {
		return fmt.Errorf("unable to add package %s: %w", layer.Name, err)
	}

	imports := make([]string, 0, 2)
	if layer.Type == consts.HTTPLayerType {
		imports = append(imports, util.MakeString("github.com/gin-gonic/gin"))
	} else {
		imports = append(imports, util.MakeString(filepath.Join(lc.Builder.ProjectName, consts.DefaultModelsFolder)))
	}
	err = g.AddImport(imports)
	if err != nil {
		return fmt.Errorf("unable to add imports %s: %w", layer.Name, err)
	}

	// Generate layer general file
	for j, mdl := range *layer.Models {
		methods := make([]model.InterfaceMethodInstance, 0, len(mdl.Methods))
		if layer.Type == consts.HTTPLayerType {
			for _, method := range (*layer.Models)[j].Methods {
				methods = append(methods, model.InterfaceMethodInstance{
					Name: method.String() + mdl.Name,
					Args: []string{"ctx", "*gin.Context"},
				})
			}
		} else {
			for _, method := range (*layer.Models)[j].Methods {
				methods = append(methods, model.InterfaceMethodInstance{
					Name:    method.String() + mdl.Name,
					Args:    []string{util.MakePrivateName(mdl.Name), "*" + consts.DefaultModelsFolder + "." + mdl.Name},
					Returns: method.GetDefaultReturns(mdl),
				})
			}
		}

		err = g.AddInterface((*layer.Models)[j].Name, methods)
		if err != nil {
			return fmt.Errorf("unable to add interface %s to file %s: %w",
				mdl.Name, layer.Name, err)
		}
	}
	return nil
}

func (lc *LayerController) generateModelLayerFile(layer *system.Layer, mdl *system.Model) error {
	g, err := gen.NewGen(layer.Path + mdl.GetFilename())
	if err != nil {
		return fmt.Errorf("unable to create new generator: %w", err)
	}
	defer g.Close()

	err = g.AddPackage(layer.Name)
	if err != nil {
		return fmt.Errorf("unable to add package %s: %w", layer.Name, err)
	}

	importsToAdd := []string{
		util.MakeString(filepath.Join(lc.Builder.ProjectName, consts.DefaultModelsFolder)),
		util.MakeString(filepath.Join(lc.Builder.ProjectName, consts.DefaultConfigFolder)),
	}
	privateMdl := mdl.GetLayer()
	privateMdl.Fields = append(privateMdl.Fields, system.Field{
		Name: consts.DefaultConfigFolder,
		Type: system.FieldType("*" + consts.DefaultConfigFolder + "." + util.MakePublicName(consts.DefaultConfigFolder)),
	})

	if layer.NextLayer != nil {
		importsToAdd = append(importsToAdd, util.MakeString(filepath.Join(lc.Builder.ProjectName, layer.NextLayer.Name)))

		privateMdl.Fields = append(privateMdl.Fields, system.Field{
			Name: util.MakePrivateName(layer.NextLayer.Name),
			Type: system.FieldType(strings.ToLower(layer.NextLayer.Name) + "." + mdl.Name),
		})
	}

	if layer.Type == consts.HTTPLayerType {
		importsToAdd = append(importsToAdd, util.MakeString(consts.GinURL))
	}

	err = g.AddImport(importsToAdd)
	if err != nil {
		return fmt.Errorf("unable to add imports %s: %w", mdl.Name, err)
	}

	err = g.AddStruct(&privateMdl)
	if err != nil {
		return fmt.Errorf("unable to add struct %s: %w", mdl.Name, err)
	}

	err = g.NewLayerFunc(layer, mdl)
	if err != nil {
		return fmt.Errorf("unable to add new layer func %s: %w", mdl.Name, err)
	}

	for _, iMdl := range mdl.Methods {
		genMethod := model.MethodInstance{
			Layer: layer,
			Args:  iMdl.GetDefaultArgs(mdl, layer),
			Type:  iMdl,
			Model: mdl,
		}
		genMethod.UpdateByMethodType()

		if layer.Type == consts.HTTPLayerType {
			genMethod.Returns = []string{""}
			err = g.AddMethodWithSwagger(&privateMdl, &genMethod)
			if err != nil {
				return fmt.Errorf("unable to add method %s with swagger: %w", mdl.Name, err)
			}

			continue
		}

		err = g.AddMethod(&privateMdl, &genMethod)
		if err != nil {
			return fmt.Errorf("unable to add method %s: %w", mdl.Name, err)
		}
	}

	return nil
}

func (lc *LayerController) generateModelStorageFile() error {
	for _, mdl := range lc.Models {
		g, err := gen.NewGen(filepath.Join(lc.Builder.Path, consts.DefaultModelsFolder, mdl.GetFilename()))
		if err != nil {
			return fmt.Errorf("unable to create new generator: %w", err)
		}

		err = g.AddPackage(consts.DefaultModelsFolder)
		if err != nil {
			return fmt.Errorf("unable to add package %s: %w", consts.DefaultModelsFolder, err)
		}

		err = g.AddStruct(mdl)
		if err != nil {
			return fmt.Errorf("unable to add struct %s: %w", mdl.Name, err)
		}

		err = g.Close()
		if err != nil {
			return fmt.Errorf("unable to close file %s: %w", mdl.Name, err)
		}
	}

	return nil
}

func (lc *LayerController) generateRouter(httpLayer *system.Layer) error {
	err := lc.generateRouterFile(httpLayer, lc.Models)
	if err != nil {
		return err
	}
	/*
		for i := range lc.Models {
			err = lc.generateSubRouterFile(httpLayer, lc.Models[i])
			if err != nil {
				return err
			}
		}*/

	return nil
}

func (lc *LayerController) generateRouterFile(layer *system.Layer, mdls []*system.Model) error {
	g, err := gen.NewGen(filepath.Join(lc.Builder.Path, consts.DefaultRouterFolder, consts.DefaultRouterFolder+".go"))
	if err != nil {
		return fmt.Errorf("unable to create new generator: %w", err)
	}

	err = g.AddPackage(
		consts.DefaultRouterFolder,
	)
	if err != nil {
		return fmt.Errorf("unable to create add packages")
	}

	err = g.AddImport([]string{
		util.MakeString(consts.GinURL),
		util.MakeString(filepath.Join(lc.Builder.ProjectName, layer.Name)),
		util.MakeString(filepath.Join(lc.Builder.ProjectName, consts.DefaultConfigFolder)),
	})
	if err != nil {
		return fmt.Errorf("unable to add imports")
	}

	routerFields := make([]system.Field, len(mdls))
	routerFields = append(routerFields, system.Field{
		Name: "Config",
		Type: "*" + system.FieldType(consts.DefaultConfigFolder+"."+util.MakePublicName("Config")),
	})
	for i := range mdls {
		routerFields[i] = system.Field{
			Name: mdls[i].Name,
			Type: system.FieldType(layer.Name + "." + mdls[i].Name),
		}
	}
	routerModel := &system.Model{
		Name:   "Router",
		Fields: routerFields,
	}
	err = g.AddStruct(routerModel)
	if err != nil {
		return fmt.Errorf("unable to add struct Router")
	}

	err = g.AddMainRouterNewFunc(routerModel)
	if err != nil {
		return err
	}

	err = g.AddMainRouterFunc(lc.Models)
	if err != nil {
		return fmt.Errorf("unable to add main router file")
	}

	err = g.Close()
	if err != nil {
		return fmt.Errorf("unable to close router generator")
	}

	return nil
}

func (lc *LayerController) generateConfigStorageFile() error {
	g, err := gen.NewGen(filepath.Join(lc.Builder.Path, consts.DefaultConfigFolder, consts.DefaultConfigFolder+".go"))
	if err != nil {
		return fmt.Errorf("unable to create new generator: %w", err)
	}

	err = g.AddPackage(consts.DefaultConfigFolder)
	if err != nil {
		return fmt.Errorf("unable to add package %s: %w", consts.DefaultConfigFolder, err)
	}

	err = g.AddImport([]string{
		util.MakeString("encoding/json"),
		util.MakeString("io"),
		util.MakeString("os"),
	})
	if err != nil {
		return fmt.Errorf("unable to add imports %s: %w", consts.DefaultConfigFolder, err)
	}

	mdlToCreate := system.Model{
		Name: util.MakePublicName(consts.DefaultConfigFolder),
	}
	httpFields := addHTTPConfig(lc.Layers)
	if len(httpFields) > 0 {
		mdlToCreate.Fields = append(mdlToCreate.Fields, httpFields...)
	}

	postgresFields := addPostgresConfig(lc.Layers)
	if len(postgresFields) > 0 {
		mdlToCreate.Fields = append(mdlToCreate.Fields, postgresFields...)
	}

	err = g.AddStruct(&mdlToCreate)
	if err != nil {
		return fmt.Errorf("unable to add struct %s: %w", mdlToCreate.Name, err)
	}

	err = g.AddParseConfigMethod(&mdlToCreate)
	if err != nil {
		return fmt.Errorf("unable to add parse config method %s: %w", mdlToCreate.Name, err)
	}

	err = g.Close()
	if err != nil {
		return fmt.Errorf("unable to close file %s: %w", mdlToCreate.Name, err)
	}

	g, err = gen.NewGen(filepath.Join(lc.Builder.Path, consts.DefaultConfigFolder+".json"))
	if err != nil {
		return fmt.Errorf("unable to create new generator: %w", err)
	}

	err = g.WriteJSON(&mdlToCreate)
	if err != nil {
		return fmt.Errorf("unable to write json %s: %w", mdlToCreate.Name, err)
	}

	return nil
}

func addHTTPConfig(layers []*system.Layer) []system.Field {
	for _, layer := range layers {
		if layer.Type == consts.HTTPLayerType {
			return []system.Field{
				{
					Name: "HttpHost",
					Type: "string",
				},
				{
					Name: "HttpPort",
					Type: "string",
				},
			}
		}
	}
	return nil
}

func addPostgresConfig(layers []*system.Layer) []system.Field {
	for _, layer := range layers {
		if layer.Type == consts.RepoLayerType {
			return []system.Field{
				{
					Name: "PostgresHost",
					Type: "string",
				},
				{
					Name: "PostgresPort",
					Type: "string",
				},
				{
					Name: "PostgresUser",
					Type: "string",
				},
				{
					Name: "PostgresPassword",
					Type: "string",
				}}
		}
	}
	return nil
}
