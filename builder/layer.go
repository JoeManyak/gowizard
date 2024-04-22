package builder

import (
	"fmt"
	"gowizard/builder/gen"
	"os"
)

const (
	httpLayerType = "http"

	repoLayerType = "postgres"
)

type LayerController struct {
	Layers  []*Layer
	Builder *Builder
	Models  []*Model
}

func NewLayerController(
	b *Builder,
	layers map[string]string,
	models []*Model,
) *LayerController {
	lc := LayerController{
		Builder: b,
		Models:  models,
		Layers:  make([]*Layer, 0, len(layers)),
	}

	for name, layerType := range layers {
		lc.Layers = append(lc.Layers, &Layer{
			Name:   name,
			Type:   layerType,
			Models: &models,
		})
	}

	for i := 0; i < len(lc.Layers)-1; i++ {
		lc.Layers[i].NextLayer = lc.Layers[i+1]
	}

	return &lc
}

type Layer struct {
	Name      string
	Type      string
	Path      string
	Models    *[]*Model
	NextLayer *Layer
}

func (lc *LayerController) Generate() error {
	for _, layer := range lc.Layers {
		f, err := os.Create(layer.Path + layer.Name + ".go")
		if err != nil {
			return fmt.Errorf("unable to create general file %s: %w", layer.Name, err)
		}

		defer f.Close()

		// Generate layer general file
		for j, model := range *layer.Models {
			methods := make([]gen.InterfaceMethod, 0, len(model.Methods))
			for _, method := range (*layer.Models)[j].Methods {
				methods = append(methods, gen.InterfaceMethod{
					Name:    method.String() + model.Name,
					Args:    []string{"test string"},
					Returns: nil,
				})
			}

			err = gen.AddInterface(f, (*layer.Models)[j].Name, methods)
			if err != nil {
				return fmt.Errorf("unable to add interface %s to file %s: %w",
					model.Name, layer.Name, err)
			}
		}
	}

	return nil
}
