package exprml_test

import (
	"fmt"

	"github.com/exprml/exprml-go"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
)

func ExampleEvaluator() {
	// Evaluate an expression

	// ExprML source code in JSON-compatible YAML format.
	source := "cat: ['`Hello`', '`, `', '`ExprML`', '`!`']"
	// Decode source as a JSON value
	decoded := exprml.NewDecoder().Decode(&pb.DecodeInput{Text: source})
	// Parse an expression
	parsed := exprml.NewParser().Parse(&pb.ParseInput{Value: decoded.Value})
	// Evaluate parsed expression
	evaluated := exprml.NewEvaluator(nil).Evaluate(&pb.EvaluateInput{Expr: parsed.Expr})
	// Encode evaluated result
	encoded := exprml.NewEncoder().Encode(&pb.EncodeInput{Value: evaluated.Value})

	fmt.Println(encoded.Text)
	// Output: "Hello, ExprML!"
}

func ExampleEvaluator_extension() {
	// Call Go functions from ExprML

	source := "$hello: { $name: '`ExprML Extension`' }"
	decoded := exprml.NewDecoder().Decode(&pb.DecodeInput{Text: source})
	parsed := exprml.NewParser().Parse(&pb.ParseInput{Value: decoded.Value})
	evaluator := exprml.NewEvaluator(&exprml.Config{
		Extension: map[string]func(path *pb.Expr_Path, args map[string]*pb.Value) *pb.EvaluateOutput{
			// Define an extension function named $hello,
			// which takes an argument $name and returns a greeting string.
			"$hello": func(path *pb.Expr_Path, args map[string]*pb.Value) *pb.EvaluateOutput {
				name, ok := args["$name"]
				if !ok || name.Type != pb.Value_STR {
					return &pb.EvaluateOutput{
						Status:       pb.EvaluateOutput_UNKNOWN_ERROR,
						ErrorPath:    path,
						ErrorMessage: "invalid argument: $name",
					}
				}
				return &pb.EvaluateOutput{
					Value: &pb.Value{Type: pb.Value_STR, Str: "Hello, " + name.Str + "!"},
				}
			},
		},
	})
	evaluated := evaluator.Evaluate(&pb.EvaluateInput{Expr: parsed.Expr})
	encoded := exprml.NewEncoder().Encode(&pb.EncodeInput{Value: evaluated.Value})
	fmt.Println(encoded.Text)
	// Output: "Hello, ExprML Extension!"
}

func ExampleEvaluator_beforeEvaluate() {
	// Hook Go functions before each evaluation of nested expressions

	source := "cat: ['`Hello`', '`, `', '`ExprML`', '`!`']"
	decoded := exprml.NewDecoder().Decode(&pb.DecodeInput{Text: source})
	parsed := exprml.NewParser().Parse(&pb.ParseInput{Value: decoded.Value})
	evaluator := exprml.NewEvaluator(&exprml.Config{
		/* Hook a function before the evaluation of each expression. */
		BeforeEvaluate: func(input *pb.EvaluateInput) error {
			fmt.Printf("Before evaluation: %q\n", exprml.Format(input.Expr.Path))
			return nil
		},
	})
	_ = evaluator.Evaluate(&pb.EvaluateInput{Expr: parsed.Expr})
	// Output:
	// Before evaluation: "/"
	// Before evaluation: "/cat/0"
	// Before evaluation: "/cat/1"
	// Before evaluation: "/cat/2"
	// Before evaluation: "/cat/3"
}

func ExampleEvaluator_afterEvaluate() {
	// Hook Go functions after each evaluation of nested expressions

	source := "cat: ['`Hello`', '`, `', '`ExprML`', '`!`']"
	decoded := exprml.NewDecoder().Decode(&pb.DecodeInput{Text: source})
	parsed := exprml.NewParser().Parse(&pb.ParseInput{Value: decoded.Value})
	evaluator := exprml.NewEvaluator(&exprml.Config{
		/* Hook a function after the evaluation of each expression. */
		AfterEvaluate: func(input *pb.EvaluateInput, output *pb.EvaluateOutput) error {
			fmt.Printf("After evaluation: %q: %v\n", exprml.Format(input.Expr.Path), output.Value.String())
			return nil
		},
	})
	_ = evaluator.Evaluate(&pb.EvaluateInput{Expr: parsed.Expr})
	// Output:
	// After evaluation: "/cat/0": type:STR  str:"Hello"
	// After evaluation: "/cat/1": type:STR  str:", "
	// After evaluation: "/cat/2": type:STR  str:"ExprML"
	// After evaluation: "/cat/3": type:STR  str:"!"
	// After evaluation: "/": type:STR  str:"Hello, ExprML!"
}
