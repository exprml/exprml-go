package exprml

import pb "github.com/exprml/exprml-go/pb/exprml/v1"

func ObjValue(obj map[string]*pb.Value) *pb.Value {
	return &pb.Value{
		Type: pb.Value_OBJ,
		Obj:  obj,
	}
}

func ArrValue(arr []*pb.Value) *pb.Value {
	return &pb.Value{
		Type: pb.Value_ARR,
		Arr:  arr,
	}
}

func StrValue(str string) *pb.Value {
	return &pb.Value{
		Type: pb.Value_STR,
		Str:  str,
	}
}

func NumValue(num float64) *pb.Value {
	return &pb.Value{
		Type: pb.Value_NUM,
		Num:  num,
	}
}

func BoolValue(b bool) *pb.Value {
	return &pb.Value{
		Type: pb.Value_BOOL,
		Bool: b,
	}
}
