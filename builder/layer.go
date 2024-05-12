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
	for _, layer := range lc.Layers {
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

	err = g.AddImport([]string{util.MakeString(filepath.Join(lc.Builder.ProjectName, consts.DefaultModelsFolder))})
	if err != nil {
		return fmt.Errorf("unable to add imports %s: %w", layer.Name, err)
	}

	// Generate layer general file
	for j, mdl := range *layer.Models {
		methods := make([]model.InterfaceMethodInstance, 0, len(mdl.Methods))
		for _, method := range (*layer.Models)[j].Methods {
			methods = append(methods, model.InterfaceMethodInstance{
				Name:    method.String() + mdl.Name,
				Args:    []string{util.MakePrivateName(mdl.Name), "*" + consts.DefaultModelsFolder + "." + mdl.Name},
				Returns: method.GetDefaultReturns(mdl),
			})
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
