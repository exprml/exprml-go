package exprml

import pb "github.com/exprml/exprml-go/pb/exprml/v1"

func ConstructNode(path *pb.Node_Path, value *pb.Value) *pb.Node {
	n := &pb.Node{
		Path:  path,
		Value: value,
	}
	switch value.Type {
	case pb.Value_ARR:
		n.Kind = pb.Node_ARRAY
		for i, v := range value.Arr {
			n.Array = append(n.Array, ConstructNode(Append(path, i), v))
		}
	case pb.Value_OBJ:
		n.Kind = pb.Node_OBJECT
		n.Object = make(map[string]*pb.Node)
		for k, v := range value.Obj {
			n.Object[k] = ConstructNode(Append(path, k), v)
		}
	default:
		n.Kind = pb.Node_SCALAR
	}

	return n
}
