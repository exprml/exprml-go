package exprml

import (
	"fmt"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
)

func Append(path *pb.Expr_Path, pos ...any) *pb.Expr_Path {
	if path == nil {
		path = &pb.Expr_Path{}
	}
	for _, pos := range pos {
		switch pos := pos.(type) {
		default:
			panic("pos must be int or string")
		case int:
			path.Pos = append(path.Pos, &pb.Expr_Path_Pos{Index: int64(pos)})
		case int64:
			path.Pos = append(path.Pos, &pb.Expr_Path_Pos{Index: pos})
		case string:
			path.Pos = append(path.Pos, &pb.Expr_Path_Pos{Key: pos})
		}
	}
	return path
}

func Format(path *pb.Expr_Path) string {
	if path == nil || len(path.Pos) == 0 {
		return "/"
	}
	s := ""
	for _, pos := range path.Pos {
		if pos.Key != "" {
			s += fmt.Sprintf("/%s", pos.Key)
		} else {
			s += fmt.Sprintf("/%d", pos.Index)
		}
	}
	return s
}
