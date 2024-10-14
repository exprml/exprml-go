package exprml_test

import (
	"fmt"
	"github.com/exprml/exprml-go"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestEvaluator_Evaluate(t *testing.T) {
	type Testcase struct {
		YamlInput string
		WantValue *pb.Value
		WantError bool
	}
	testcases := map[string]*Testcase{}
	err := filepath.Walk("testdata/evaluator", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".in.yaml") {
			key := strings.TrimSuffix(path, ".in.yaml")
			b, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("fail to read file: %+v", err)
			}
			if _, ok := testcases[key]; !ok {
				testcases[key] = &Testcase{}
			}
			testcases[key].YamlInput = string(b)
			return nil
		}
		if strings.HasSuffix(path, ".want.yaml") {
			key := strings.TrimSuffix(path, ".want.yaml")
			b, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("fail to read file: %+v", err)
			}
			want := exprml.NewDecoder().Decode(&pb.DecodeInput{Yaml: string(b)})
			if want.IsError {
				return fmt.Errorf("fail to decode yaml: %+v", want.ErrorMessage)
			}
			if _, ok := testcases[key]; !ok {
				testcases[key] = &Testcase{}
			}
			if v, ok := want.Value.Obj["want_value"]; ok {
				testcases[key].WantValue = v
			}
			if v, ok := want.Value.Obj["want_error"]; ok {
				testcases[key].WantError = v.Bool
			}
			if testcases[key].WantValue == nil && !testcases[key].WantError {
				return fmt.Errorf("want_value or want_error is not found in %v", path)
			}
			return nil
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	testcaseKeys := []string{}
	for key := range testcases {
		testcaseKeys = append(testcaseKeys, key)
	}
	slices.Sort(testcaseKeys)
	for _, name := range testcaseKeys {
		testcase := testcases[name]
		t.Run(name, func(t *testing.T) {
			decodeResult := exprml.NewDecoder().Decode(&pb.DecodeInput{Yaml: testcase.YamlInput})
			require.False(t, decodeResult.IsError)

			parseResult := exprml.NewParser().Parse(&pb.ParseInput{Value: decodeResult.Value})
			require.False(t, parseResult.IsError)

			got := exprml.NewEvaluator(nil).EvaluateExpr(&pb.EvaluateInput{Expr: parseResult.Expr})

			if testcase.WantError {
				require.NotEqual(t, pb.EvaluateOutput_OK, got.Status)
			} else {
				require.Equal(t, pb.EvaluateOutput_OK, got.Status)
				require.Nil(t, checkEqual([]string{}, testcase.WantValue, got.Value))
			}
		})
	}
}

func TestEvaluator_Extension(t *testing.T) {
	sut := exprml.NewEvaluator(&exprml.EvaluatorConfig{
		Extension: map[string]func(path *pb.Expr_Path, args map[string]*pb.Value) *pb.EvaluateOutput{
			"$test_func": func(path *pb.Expr_Path, args map[string]*pb.Value) *pb.EvaluateOutput {
				return &pb.EvaluateOutput{
					Value: &pb.Value{Type: pb.Value_OBJ, Obj: args},
				}
			},
		},
	})
	type Testcase struct {
		YamlInput string
		WantValue *pb.Value
		WantError bool
	}
	testcases := map[string]*Testcase{
		"Ref": {
			YamlInput: "$test_func",
			WantValue: &pb.Value{
				Type: pb.Value_OBJ,
				Obj:  map[string]*pb.Value{},
			},
		},
		"Call": {
			YamlInput: "$test_func: { $arg: '`value`' }",
			WantValue: &pb.Value{
				Type: pb.Value_OBJ,
				Obj: map[string]*pb.Value{
					"$arg": {Type: pb.Value_STR, Str: "value"},
				},
			},
		},
	}

	testcaseKeys := []string{}
	for key := range testcases {
		testcaseKeys = append(testcaseKeys, key)
	}
	slices.Sort(testcaseKeys)
	for _, name := range testcaseKeys {
		testcase := testcases[name]
		t.Run(name, func(t *testing.T) {
			decodeResult := exprml.NewDecoder().Decode(&pb.DecodeInput{Yaml: testcase.YamlInput})
			require.False(t, decodeResult.IsError)

			parseResult := exprml.NewParser().Parse(&pb.ParseInput{Value: decodeResult.Value})
			require.False(t, parseResult.IsError)

			got := sut.EvaluateExpr(&pb.EvaluateInput{Expr: parseResult.Expr})
			if testcase.WantError {
				require.NotEqual(t, pb.EvaluateOutput_OK, got.Status)
			} else {
				require.Equal(t, pb.EvaluateOutput_OK, got.Status)
				require.Nil(t, checkEqual([]string{}, testcase.WantValue, got.Value))
			}
		})
	}
}

func TestEvaluator_BeforeEvaluate(t *testing.T) {
	evalPaths := []string{}

	decodeResult := exprml.NewDecoder().Decode(&pb.DecodeInput{Yaml: "cat: ['`Hello`', '`, `', '`ExprML`', '`!`']"})
	require.False(t, decodeResult.IsError)

	parseResult := exprml.NewParser().Parse(&pb.ParseInput{Value: decodeResult.Value})
	require.False(t, parseResult.IsError)

	config := &exprml.EvaluatorConfig{
		BeforeEvaluate: func(input *pb.EvaluateInput) error {
			evalPaths = append(evalPaths, exprml.Format(input.Expr.Path))
			return nil
		},
	}
	result := exprml.NewEvaluator(config).EvaluateExpr(&pb.EvaluateInput{Expr: parseResult.Expr})
	require.Equal(t, pb.EvaluateOutput_OK, result.Status)
	require.ElementsMatch(t, []string{"/", "/cat/0", "/cat/1", "/cat/2", "/cat/3"}, evalPaths)
}

func TestEvaluator_AfterEvaluate(t *testing.T) {
	evalTypes := []pb.Value_Type{}

	decodeResult := exprml.NewDecoder().Decode(&pb.DecodeInput{Yaml: "cat: ['`Hello`', '`, `', '`ExprML`', '`!`']"})
	require.False(t, decodeResult.IsError)

	parseResult := exprml.NewParser().Parse(&pb.ParseInput{Value: decodeResult.Value})
	require.False(t, parseResult.IsError)

	config := &exprml.EvaluatorConfig{
		AfterEvaluate: func(input *pb.EvaluateInput, output *pb.EvaluateOutput) error {
			evalTypes = append(evalTypes, output.Value.Type)
			return nil
		},
	}
	result := exprml.NewEvaluator(config).EvaluateExpr(&pb.EvaluateInput{Expr: parseResult.Expr})
	require.Equal(t, pb.EvaluateOutput_OK, result.Status)

	wantTypes := []pb.Value_Type{pb.Value_STR, pb.Value_STR, pb.Value_STR, pb.Value_STR, pb.Value_STR}
	require.ElementsMatch(t, wantTypes, evalTypes)
}

func checkEqual(path []string, want, got *pb.Value) error {
	p := "/" + strings.Join(path, "/")
	if want.Type != got.Type {
		return fmt.Errorf("type mismatch: path=%v, got=%v, want=%v", p, got.Type, want.Type)
	}
	switch want.Type {
	default:
		return fmt.Errorf("unexpected type: path=%v, got=%v, want=%v", p, got.Type, want.Type)
	case pb.Value_NULL:
		return nil
	case pb.Value_BOOL:
		if want.Bool != got.Bool {
			return fmt.Errorf("boolean mismatch: path=%v, got=%v, want=%v", p, got.Bool, want.Bool)
		}
		return nil
	case pb.Value_NUM:
		if want.Num != got.Num {
			return fmt.Errorf("number mismatch: path=%v, got=%v, want=%v", p, got.Num, want.Num)
		}
		return nil
	case pb.Value_STR:
		if want.Str != got.Str {
			return fmt.Errorf("string mismatch: path=%v, got=%v, want=%v", p, got.Str, want.Str)
		}
		return nil
	case pb.Value_ARR:
		if len(want.Arr) != len(got.Arr) {
			return fmt.Errorf("array length mismatch: path=%v, got=%v, want=%v", p, len(got.Arr), len(want.Arr))
		}
		for i := 0; i < len(want.Arr); i++ {
			path := append([]string{}, path...)
			if err := checkEqual(append(path, fmt.Sprintf("%v", i)), want.Arr[i], got.Arr[i]); err != nil {
				return err
			}
		}
		return nil
	case pb.Value_OBJ:
		wk, gk := sortedKeys(want.Obj), sortedKeys(got.Obj)
		if !slices.Equal(wk, gk) {
			return fmt.Errorf("object keys mismatch: path=%v, got=[%v], want=[%v]", p, strings.Join(wk, ","), strings.Join(gk, ","))
		}
		for _, k := range wk {
			path := append([]string{}, path...)
			if err := checkEqual(append(path, k), want.Obj[k], got.Obj[k]); err != nil {
				return err
			}
		}
		return nil
	}
}

func sortedKeys(m map[string]*pb.Value) []string {
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
