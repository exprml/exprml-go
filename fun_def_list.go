package exprml

import pb "github.com/exprml/exprml-go/pb/exprml/v1"

func Register(funDefList *pb.FunDefList, funDef *pb.Node) *pb.FunDefList {
	return &pb.FunDefList{
		Parent: funDefList,
		FunDef: funDef,
	}
}

func Find(funDefList *pb.FunDefList, ident string) *pb.FunDefList {
	if funDefList == nil || funDefList.FunDef == nil {
		return nil
	}
	if funDefList.FunDef.Object["def"] == nil {
		panic("invalid FunDef")
	}
	if funDefList.FunDef.Object["def"].Value.Str == ident {
		return funDefList
	}
	return Find(funDefList.Parent, ident)
}
