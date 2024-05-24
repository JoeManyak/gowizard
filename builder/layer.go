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
	layerTypes := make(map[string]struct{})
	for _, layer := range lc.Layers {
		layerTypes[layer.Type] = struct{}{}
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

	if _, ok := layerTypes[consts.HTTPLayerType]; ok {
		err = lc.generateRouter()
		if err != nil {
			return err
		}
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

	importsToAdd := []string{util.MakeString(filepath.Join(lc.Builder.ProjectName, consts.DefaultModelsFolder))}
	privateMdl := mdl.GetLayer()
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

func (lc *LayerController) generateRouter() error {
	err := lc.generateRouterFile()
	if err != nil {
		return err
	}

	for i := range lc.Models {
		err = lc.generateSubRouterFile(lc.Models[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (lc *LayerController) generateRouterFile() error {
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

	err = g.AddImport([]string{util.MakeString(consts.GinURL)})
	if err != nil {
		return fmt.Errorf("unable to add imports")
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

func (lc *LayerController) generateSubRouterFile(mdl *system.Model) error {
	g, err := gen.NewGen(filepath.Join(lc.Builder.Path, consts.DefaultRouterFolder, mdl.GetFilename()))
	if err != nil {
		return fmt.Errorf("unable to create new generator: %w", err)
	}

	err = g.AddPackage(
		consts.DefaultRouterFolder,
	)
	if err != nil {
		return fmt.Errorf("unable to create add packages")
	}

	err = g.AddImport([]string{util.MakeString(consts.GinURL)})
	if err != nil {
		return fmt.Errorf("unable to add imports")
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
