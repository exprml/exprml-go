package exprml

import pb "github.com/exprml/exprml-go/pb/exprml/v1"

func Register(defStack *pb.DefStack, def *pb.Eval_Definition) *pb.DefStack {
	return &pb.DefStack{
		Parent: defStack,
		Def:    def,
	}
}

func Find(defStack *pb.DefStack, ident string) *pb.DefStack {
	if defStack == nil || defStack.Def == nil {
		return nil
	}
	if defStack.Def.Ident == ident {
		return defStack
	}
	return Find(defStack.Parent, ident)
}

func NewDefinition(path *pb.Expr_Path, ident string, value *pb.Value) *pb.Eval_Definition {
	return &pb.Eval_Definition{
		Ident: ident,
		Body: &pb.Expr{
			Path:  path,
			Kind:  pb.Expr_JSON,
			Value: value,
			Json:  &pb.Json{Json: value},
		},
	}
}
