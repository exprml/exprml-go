package exprml

import (
	"encoding/json"
	"fmt"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
	"github.com/goccy/go-yaml"
)

type Decoder interface {
	Decode(input *pb.DecodeInput) *pb.DecodeOutput
}

func NewDecoder() Decoder {
	return &decoder{}
}

type decoder struct{}

func (d decoder) Decode(input *pb.DecodeInput) *pb.DecodeOutput {
	b, err := yaml.YAMLToJSON([]byte(input.Text))
	if err != nil {
		return &pb.DecodeOutput{
			IsError:      true,
			ErrorMessage: fmt.Sprintf("fail to convertFromGo yaml to json: %+v", err),
		}
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return &pb.DecodeOutput{
			IsError:      true,
			ErrorMessage: fmt.Sprintf("fail to unmarshal json: %+v", err),
		}
	}
	return &pb.DecodeOutput{Value: convertFromGo(v)}
}

func convertFromGo(v any) *pb.Value {
	switch v := v.(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T", v))
	case nil:
		return &pb.Value{Type: pb.Value_NULL}
	case bool:
		return &pb.Value{Type: pb.Value_BOOL, Bool: v}
	case float64:
		return &pb.Value{Type: pb.Value_NUM, Num: v}
	case string:
		return &pb.Value{Type: pb.Value_STR, Str: v}
	case []interface{}:
		arr := []*pb.Value{}
		for _, elem := range v {
			arr = append(arr, convertFromGo(elem))
		}
		return &pb.Value{Type: pb.Value_ARR, Arr: arr}
	case map[string]interface{}:
		obj := map[string]*pb.Value{}
		for key, value := range v {
			obj[key] = convertFromGo(value)
		}
		return &pb.Value{Type: pb.Value_OBJ, Obj: obj}
	}
}
