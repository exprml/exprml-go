package exprml

import (
	"bytes"
	"encoding/json"
	"fmt"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
	"github.com/goccy/go-yaml"
)

type Encoder interface {
	Encode(input *pb.EncodeInput) *pb.EncodeOutput
}

func NewEncoder() Encoder {
	return &encoder{}
}

type encoder struct{}

func (d encoder) Encode(input *pb.EncodeInput) *pb.EncodeOutput {
	g := convertToGo(input.Value)
	switch input.Format {
	default:
		return &pb.EncodeOutput{IsError: true, ErrorMessage: fmt.Sprintf("unexpected format %v", input.Format)}
	case pb.EncodeInput_JSON:
		b := bytes.NewBuffer(nil)
		e := json.NewEncoder(b)
		if err := e.Encode(g); err != nil {
			return &pb.EncodeOutput{
				IsError:      true,
				ErrorMessage: fmt.Sprintf("fail to encode json: %+v", err),
			}
		}
		return &pb.EncodeOutput{Result: b.String()}
	case pb.EncodeInput_YAML:
		b := bytes.NewBuffer(nil)
		e := yaml.NewEncoder(b)
		if err := e.Encode(g); err != nil {
			return &pb.EncodeOutput{
				IsError:      true,
				ErrorMessage: fmt.Sprintf("fail to encode yaml: %+v", err),
			}
		}
		return &pb.EncodeOutput{Result: b.String()}
	}
}

func convertToGo(v *pb.Value) any {
	switch v.Type {
	default:
		panic(fmt.Sprintf("unexpected type %v", v.Type))
	case pb.Value_NULL:
		return nil
	case pb.Value_BOOL:
		return v.Bool
	case pb.Value_NUM:
		return v.Num
	case pb.Value_STR:
		return v.Str
	case pb.Value_ARR:
		arr := []any{}
		for _, elem := range v.Arr {
			arr = append(arr, convertToGo(elem))
		}
		return arr
	case pb.Value_OBJ:
		obj := map[string]any{}
		for key, value := range v.Obj {
			obj[key] = convertToGo(value)
		}
		return obj
	}
}
