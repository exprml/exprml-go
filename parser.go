package exprml

import pb "github.com/exprml/exprml-go/pb/exprml/v1"

type Parser interface {
	Parse(input *pb.ParseInput) *pb.ParseOutput
}

func NewParser() Parser {
	return &parser{}
}

type parser struct{}

func (p parser) Parse(input *pb.ParseInput) *pb.ParseOutput {
	return &pb.ParseOutput{
		Node: ConstructNode(&pb.Node_Path{}, input.Value),
	}
}
