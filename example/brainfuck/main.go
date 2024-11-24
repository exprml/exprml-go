package main

import (
	_ "embed"
	"fmt"
	"github.com/exprml/exprml-go"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
	"os"
)

//go:embed brainfuck.exprml.yaml
var brainfuckExprMLYAML string

//go:embed source.bf
var sourceBF string

func main() {
	decoded := exprml.NewDecoder().Decode(&pb.DecodeInput{Text: brainfuckExprMLYAML})
	parsed := exprml.NewParser().Parse(&pb.ParseInput{Value: decoded.Value})
	evaluator := exprml.NewEvaluator(&exprml.Config{
		Extension: map[string]func(path *pb.Expr_Path, args map[string]*pb.Value) *pb.EvaluateOutput{
			"$input": func(path *pb.Expr_Path, _ map[string]*pb.Value) *pb.EvaluateOutput {
				b := []byte{0}
				_, _ = os.Stdin.Read(b)
				return &pb.EvaluateOutput{
					Value: &pb.Value{Type: pb.Value_NUM, Num: float64(int64(b[0]))},
				}
			},
			"$output": func(path *pb.Expr_Path, args map[string]*pb.Value) *pb.EvaluateOutput {
				value := args["$value"]
				fmt.Printf("%c", byte(int64(value.Num)))
				return &pb.EvaluateOutput{Value: value}
			},
			"$source": func(path *pb.Expr_Path, _ map[string]*pb.Value) *pb.EvaluateOutput {
				return &pb.EvaluateOutput{
					Value: &pb.Value{Type: pb.Value_STR, Str: sourceBF},
				}
			},
		},
	})
	evaluated := evaluator.Evaluate(&pb.EvaluateInput{Expr: parsed.Expr})
	if evaluated.Status != pb.EvaluateOutput_OK {
		fmt.Println(evaluated.String())
	} else {
		encoded := exprml.NewEncoder().Encode(&pb.EncodeInput{Value: evaluated.Value})
		fmt.Println(encoded.Text)
	}
}
