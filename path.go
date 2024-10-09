package exprml

import pb "github.com/exprml/exprml-go/pb/exprml/v1"

func Append(path *pb.Node_Path, pos ...any) *pb.Node_Path {
	if path == nil {
		path = &pb.Node_Path{}
	}
	for _, pos := range pos {
		switch pos := pos.(type) {
		default:
			panic("pos must be int or string")
		case int:
			path.Pos = append(path.Pos, &pb.Node_Path_Pos{Index: int64(pos)})
		case int64:
			path.Pos = append(path.Pos, &pb.Node_Path_Pos{Index: pos})
		case string:
			path.Pos = append(path.Pos, &pb.Node_Path_Pos{Key: pos})
		}
	}
	return path
}
