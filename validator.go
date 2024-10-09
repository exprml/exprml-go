package exprml

import (
	"bytes"
	_ "embed"
	"fmt"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
	"github.com/goccy/go-yaml"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

type Validator interface {
	Validate(input *pb.ValidateInput) *pb.ValidateOutput
}

func NewValidator() Validator {
	v, err := jsonschema.UnmarshalJSON(bytes.NewBuffer(schemaJSON))
	if err != nil {
		panic(fmt.Sprintf("fail to decode json: %+v", err))
	}

	url := "https://github.com/exprml/exprml-language/schema.json"
	c := jsonschema.NewCompiler()
	if err := c.AddResource(url, v); err != nil {
		panic(fmt.Sprintf("fail to add schema: %+v", err))
	}

	return &validator{schema: c.MustCompile(url)}
}

type validator struct {
	schema *jsonschema.Schema
}

//go:embed schema.json
var schemaJSON []byte

func (v validator) Validate(input *pb.ValidateInput) *pb.ValidateOutput {
	b, err := yaml.YAMLToJSON([]byte(input.Source))
	if err != nil {
		return &pb.ValidateOutput{
			Status:       pb.ValidateOutput_YAML_ERROR,
			ErrorMessage: fmt.Sprintf("fail to convert yaml to json: %+v", err),
		}
	}

	sourceJSON, err := jsonschema.UnmarshalJSON(bytes.NewBuffer(b))
	if err != nil {
		return &pb.ValidateOutput{
			Status:       pb.ValidateOutput_YAML_ERROR,
			ErrorMessage: fmt.Sprintf("fail to decode json: %+v", err),
		}
	}

	if err := v.schema.Validate(sourceJSON); err != nil {
		vo := &pb.ValidateOutput{
			Status:       pb.ValidateOutput_VALIDATION_ERROR,
			ErrorMessage: fmt.Sprintf("validation error: %#+v", err),
		}
		return vo
	}

	return &pb.ValidateOutput{Status: pb.ValidateOutput_OK}
}
