package exprml

import (
	"fmt"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
	"math"
	"slices"
	"strings"
)

type Evaluator interface {
	EvaluateExpr(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateEval(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateScalar(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateObj(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateArr(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateJson(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateRangeIter(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateGetElem(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateFunCall(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateCases(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateOpUnary(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateOpBinary(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateOpVariadic(input *pb.EvaluateInput) *pb.EvaluateOutput
}

func NewEvaluator() Evaluator {
	return &BasicEvaluator{}
}

type BasicEvaluator struct{}

func hasKey(v *pb.Node, key string) bool {
	if v.Kind != pb.Node_OBJECT {
		return false
	}
	_, ok := v.Object[key]
	return ok
}
func (e BasicEvaluator) EvaluateExpr(input *pb.EvaluateInput) *pb.EvaluateOutput {
	n := input.Expr
	switch n.Kind {
	case pb.Node_SCALAR:
		return e.EvaluateScalar(input)
	case pb.Node_OBJECT:
		switch {
		case hasKey(n, "eval"):
			return e.EvaluateEval(input)
		case hasKey(n, "obj"):
			return e.EvaluateObj(input)
		case hasKey(n, "arr"):
			return e.EvaluateArr(input)
		case hasKey(n, "json"):
			return e.EvaluateJson(input)
		case hasKey(n, "for"):
			return e.EvaluateRangeIter(input)
		case hasKey(n, "get"):
			return e.EvaluateGetElem(input)
		case hasKey(n, "ref"):
			return e.EvaluateFunCall(input)
		case hasKey(n, "cases"):
			return e.EvaluateCases(input)
		case hasKey(n, "len"),
			hasKey(n, "not"),
			hasKey(n, "flat"),
			hasKey(n, "floor"),
			hasKey(n, "ceil"),
			hasKey(n, "abort"):
			return e.EvaluateOpUnary(input)
		case hasKey(n, "sub"),
			hasKey(n, "div"),
			hasKey(n, "eq"),
			hasKey(n, "neq"),
			hasKey(n, "lt"),
			hasKey(n, "lte"),
			hasKey(n, "gt"),
			hasKey(n, "gte"):
			return e.EvaluateOpBinary(input)
		case hasKey(n, "add"),
			hasKey(n, "mul"),
			hasKey(n, "and"),
			hasKey(n, "or"),
			hasKey(n, "cat"),
			hasKey(n, "min"),
			hasKey(n, "max"),
			hasKey(n, "merge"):
			return e.EvaluateOpVariadic(input)
		}
	}
	panic("given expression must be validated")
}

func (e BasicEvaluator) EvaluateEval(input *pb.EvaluateInput) *pb.EvaluateOutput {
	st := input.DefStack
	if where, ok := input.Expr.Object["where"]; ok {
		for _, funDef := range where.Array {
			st = Register(st, funDef)
		}
	}
	return e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: input.Expr.Object["eval"]})
}

func (e BasicEvaluator) EvaluateScalar(input *pb.EvaluateInput) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{Value: input.Expr.Value}
}

func (e BasicEvaluator) EvaluateObj(input *pb.EvaluateInput) *pb.EvaluateOutput {
	obj := input.Expr.Object["obj"]
	result := map[string]*pb.Value{}
	for pos, expr := range obj.Object {
		val := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: expr})
		if val.Status != pb.EvaluateOutput_OK {
			return val
		}
		result[pos] = val.Value
	}
	return &pb.EvaluateOutput{Value: ObjValue(result)}
}

func (e BasicEvaluator) EvaluateArr(input *pb.EvaluateInput) *pb.EvaluateOutput {
	arr := input.Expr.Object["arr"]
	result := []*pb.Value{}
	for _, expr := range arr.Array {
		val := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: expr})
		if val.Status != pb.EvaluateOutput_OK {
			return val
		}
		result = append(result, val.Value)
	}
	return &pb.EvaluateOutput{Value: ArrValue(result)}
}

func (e BasicEvaluator) EvaluateJson(input *pb.EvaluateInput) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{Value: input.Expr.Value.Obj["json"]}
}

func (e BasicEvaluator) EvaluateRangeIter(input *pb.EvaluateInput) *pb.EvaluateOutput {
	in := input.Expr.Object["in"]
	inVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: in})
	forPos, forElem := input.Expr.Object["for"].Array[0], input.Expr.Object["for"].Array[1]
	switch inVal.Value.Type {
	default:
		return errorUnexpectedType(in.Path, inVal.Value.Type, []pb.Value_Type{pb.Value_ARR, pb.Value_OBJ})
	case pb.Value_STR:
		result := []*pb.Value{}
		for i, c := range []rune(inVal.Value.Str) {
			st := input.DefStack
			st = Register(st, NewFunDef(forPos.Path, forPos.Value.Str, NumValue(float64(i))))
			st = Register(st, NewFunDef(forElem.Path, forElem.Value.Str, StrValue(string(c))))
			if if_, ok := input.Expr.Object["if"]; ok {
				ifVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: if_})
				if ifVal.Status != pb.EvaluateOutput_OK {
					return ifVal
				}
				if ifVal.Value.Type != pb.Value_BOOL {
					return errorUnexpectedType(if_.Path, ifVal.Value.Type, []pb.Value_Type{pb.Value_BOOL})
				}
				if !ifVal.Value.Bool {
					continue
				}
			}
			v := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: input.Expr.Object["do"]})
			if v.Status != pb.EvaluateOutput_OK {
				return v
			}
			result = append(result, v.Value)
		}
		return &pb.EvaluateOutput{Value: ArrValue(result)}
	case pb.Value_ARR:
		result := []*pb.Value{}
		for i, elemVal := range inVal.Value.Arr {
			st := input.DefStack
			st = Register(st, NewFunDef(forPos.Path, forPos.Value.Str, NumValue(float64(i))))
			st = Register(st, NewFunDef(forElem.Path, forElem.Value.Str, elemVal))
			if if_, ok := input.Expr.Object["if"]; ok {
				ifVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: if_})
				if ifVal.Status != pb.EvaluateOutput_OK {
					return ifVal
				}
				if ifVal.Value.Type != pb.Value_BOOL {
					return errorUnexpectedType(if_.Path, ifVal.Value.Type, []pb.Value_Type{pb.Value_BOOL})
				}
				if !ifVal.Value.Bool {
					continue
				}
			}
			v := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: input.Expr.Object["do"]})
			if v.Status != pb.EvaluateOutput_OK {
				return v
			}
			result = append(result, v.Value)
		}
		return &pb.EvaluateOutput{Value: ArrValue(result)}
	case pb.Value_OBJ:
		result := map[string]*pb.Value{}
		for _, key := range sortedKeys(inVal.Value.Obj) {
			st := input.DefStack
			st = Register(st, NewFunDef(forPos.Path, forPos.Value.Str, StrValue(key)))
			st = Register(st, NewFunDef(forElem.Path, forElem.Value.Str, inVal.Value.Obj[key]))
			if if_, ok := input.Expr.Object["if"]; ok {
				ifVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: if_})
				if ifVal.Status != pb.EvaluateOutput_OK {
					return ifVal
				}
				if ifVal.Value.Type != pb.Value_BOOL {
					return errorUnexpectedType(if_.Path, ifVal.Value.Type, []pb.Value_Type{pb.Value_BOOL})
				}
				if !ifVal.Value.Bool {
					continue
				}
			}
			v := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: input.Expr.Object["do"]})
			if v.Status != pb.EvaluateOutput_OK {
				return v
			}
			result[key] = v.Value
		}
		return &pb.EvaluateOutput{Value: ObjValue(result)}
	}
}

func (e BasicEvaluator) EvaluateGetElem(input *pb.EvaluateInput) *pb.EvaluateOutput {
	var pos, col *pb.Value
	get := input.Expr.Object["get"]
	{
		getVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: get})
		if getVal.Status != pb.EvaluateOutput_OK {
			return getVal
		}
		pos = getVal.Value
	}
	from := input.Expr.Object["from"]
	{
		fromVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: from})
		if fromVal.Status != pb.EvaluateOutput_OK {
			return fromVal
		}
		col = fromVal.Value
	}

	switch col.Type {
	default:
		return errorUnexpectedType(from.Path, col.Type, []pb.Value_Type{pb.Value_STR, pb.Value_ARR, pb.Value_OBJ})
	case pb.Value_STR:
		if pos.Type != pb.Value_NUM {
			return errorUnexpectedType(get.Path, pos.Type, []pb.Value_Type{pb.Value_NUM})
		}
		if !canInt(pos) {
			return errorIndexNotInteger(get.Path, pos.Num)
		}
		idx := int(pos.Num)
		if idx < 0 || idx >= len([]rune(col.Str)) {
			return errorIndexOutOfBounds(get.Path, pos, 0, len(from.Value.Arr))
		}
		return &pb.EvaluateOutput{Value: StrValue(string([]rune(col.Str)[idx]))}
	case pb.Value_ARR:
		if pos.Type != pb.Value_NUM {
			return errorUnexpectedType(get.Path, pos.Type, []pb.Value_Type{pb.Value_NUM})
		}
		if !canInt(pos) {
			return errorIndexNotInteger(get.Path, pos.Num)
		}
		idx := int(pos.Num)
		if idx < 0 || idx >= len(col.Arr) {
			return errorIndexOutOfBounds(get.Path, pos, 0, len(col.Arr))
		}
		return &pb.EvaluateOutput{Value: col.Arr[idx]}
	case pb.Value_OBJ:
		if pos.Type != pb.Value_STR {
			return errorUnexpectedType(get.Path, pos.Type, []pb.Value_Type{pb.Value_STR})
		}
		key := pos.Str
		if _, ok := col.Obj[key]; !ok {
			return errorInvalidKey(get.Path, key, sortedKeys(col.Obj))
		}
		return &pb.EvaluateOutput{Value: col.Obj[key]}
	}
}

func (e BasicEvaluator) EvaluateFunCall(input *pb.EvaluateInput) *pb.EvaluateOutput {
	funCall := input.Expr
	ref := funCall.Object["ref"]
	st := Find(input.DefStack, ref.Value.Str)
	if st == nil {
		return errorReferenceNotFound(ref.Path, ref.Value.Str)
	}
	funDef := st.FunDef
	if funDefWith, ok := funDef.Object["with"]; ok {
		for _, argName := range funDefWith.Array {
			with, ok := funCall.Object["with"]
			if !ok {
				return errorArgumentMismatch(funCall.Path, argName.Value.Str)
			}
			argExpr, ok := with.Object[argName.Value.Str]
			if !ok {
				return errorArgumentMismatch(with.Path, argName.Value.Str)
			}
			argVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: argExpr})
			if argVal.Status != pb.EvaluateOutput_OK {
				return argVal
			}
			st = Register(st, NewFunDef(argExpr.Path, argName.Value.Str, argVal.Value))
		}
	}
	return e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: funDef.Object["value"]})
}

func (e BasicEvaluator) EvaluateCases(input *pb.EvaluateInput) *pb.EvaluateOutput {
	cases := input.Expr.Object["cases"]
	for _, case_ := range cases.Array {
		switch {
		case hasKey(case_, "when"):
			boolVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: case_.Object["when"]})
			if boolVal.Status != pb.EvaluateOutput_OK {
				return boolVal
			}
			if boolVal.Value.Type != pb.Value_BOOL {
				return errorUnexpectedType(case_.Object["when"].Path, boolVal.Value.Type, []pb.Value_Type{pb.Value_BOOL})
			}
			if boolVal.Value.Bool {
				return e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: case_.Object["then"]})
			}
		case hasKey(case_, "otherwise"):
			return e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: case_.Object["otherwise"]})
		}
	}
	return errorCasesNotExhaustive(cases.Path)
}

func (e BasicEvaluator) EvaluateOpUnary(input *pb.EvaluateInput) *pb.EvaluateOutput {
	var (
		operator string
		operand  *pb.Value
	)
	for k, v := range input.Expr.Object { // only one property exists
		operator = k
		o := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: v})
		if o.Status != pb.EvaluateOutput_OK {
			return o
		}
		operand = o.Value
	}
	switch operator {
	default:
		panic(fmt.Sprintf("unexpected unary operator %q", operator))
	case "len":
		switch operand.Type {
		default:
			return errorUnexpectedType(Append(input.Expr.Path, operator), operand.Type, []pb.Value_Type{pb.Value_STR, pb.Value_ARR, pb.Value_OBJ})
		case pb.Value_STR:
			return &pb.EvaluateOutput{Value: NumValue(float64(len([]rune(operand.Str))))}
		case pb.Value_ARR:
			return &pb.EvaluateOutput{Value: NumValue(float64(len(operand.Arr)))}
		case pb.Value_OBJ:
			return &pb.EvaluateOutput{Value: NumValue(float64(len(operand.Obj)))}
		}
	case "not":
		if operand.Type != pb.Value_BOOL {
			return errorUnexpectedType(Append(input.Expr.Path, operator), operand.Type, []pb.Value_Type{pb.Value_BOOL})
		}
		return &pb.EvaluateOutput{Value: BoolValue(!operand.Bool)}
	case "flat":
		if operand.Type != pb.Value_ARR {
			return errorUnexpectedType(Append(input.Expr.Path, operator), operand.Type, []pb.Value_Type{pb.Value_ARR})
		}
		v := []*pb.Value{}
		for _, elem := range operand.Arr {
			if elem.Type != pb.Value_ARR {
				return errorUnexpectedType(Append(input.Expr.Path, operator), elem.Type, []pb.Value_Type{pb.Value_ARR})
			}
			v = append(v, elem.Arr...)
		}
		return &pb.EvaluateOutput{Value: ArrValue(v)}
	case "floor":
		if operand.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, operator), operand.Type, []pb.Value_Type{pb.Value_NUM})
		}
		return &pb.EvaluateOutput{Value: NumValue(math.Floor(operand.Num))}
	case "ceil":
		if operand.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, operator), operand.Type, []pb.Value_Type{pb.Value_NUM})
		}
		return &pb.EvaluateOutput{Value: NumValue(math.Ceil(operand.Num))}
	case "abort":
		if operand.Type != pb.Value_STR {
			return errorUnexpectedType(Append(input.Expr.Path, operator), operand.Type, []pb.Value_Type{pb.Value_STR})
		}
		return &pb.EvaluateOutput{Status: pb.EvaluateOutput_ABORTED, ErrorMessage: operand.Str}
	}
}

func (e BasicEvaluator) EvaluateOpBinary(input *pb.EvaluateInput) *pb.EvaluateOutput {
	var (
		operator           string
		operandL, operandR *pb.Value
	)
	for k, v := range input.Expr.Object { // only one property exists
		operator = k
		ol := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: v.Array[0]})
		if ol.Status != pb.EvaluateOutput_OK {
			return ol
		}
		or := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: v.Array[1]})
		if or.Status != pb.EvaluateOutput_OK {
			return or
		}
		operandL, operandR = ol.Value, or.Value
	}
	switch operator {
	default:
		panic(fmt.Sprintf("unexpected binary operator %q", operator))
	case "sub":
		if operandL.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, operator, 0), operandL.Type, []pb.Value_Type{pb.Value_NUM})
		}
		if operandR.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, operator, 1), operandR.Type, []pb.Value_Type{pb.Value_NUM})
		}
		v := operandL.Num - operandR.Num
		if !isFiniteNumber(v) {
			return errorNotFiniteNumber(input.Expr.Path)
		}
		return &pb.EvaluateOutput{Value: NumValue(v)}
	case "div":
		if operandL.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, operator, 0), operandL.Type, []pb.Value_Type{pb.Value_NUM})
		}
		if operandR.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, operator, 1), operandR.Type, []pb.Value_Type{pb.Value_NUM})
		}
		v := operandL.Num / operandR.Num
		if !isFiniteNumber(v) {
			return errorNotFiniteNumber(input.Expr.Path)
		}
		return &pb.EvaluateOutput{Value: NumValue(v)}
	case "eq":
		return &pb.EvaluateOutput{Value: equal(operandL, operandR)}
	case "neq":
		return &pb.EvaluateOutput{Value: BoolValue(!equal(operandL, operandR).Bool)}
	case "lt":
		cmpVal := compare(Append(input.Expr.Path, operator), operandL, operandR)
		if cmpVal.Status != pb.EvaluateOutput_OK {
			return cmpVal
		}
		return &pb.EvaluateOutput{Value: BoolValue(cmpVal.Value.Num < 0)}
	case "lte":
		cmpVal := compare(Append(input.Expr.Path, operator), operandL, operandR)
		if cmpVal.Status != pb.EvaluateOutput_OK {
			return cmpVal
		}
		return &pb.EvaluateOutput{Value: BoolValue(cmpVal.Value.Num <= 0)}
	case "gt":
		cmpVal := compare(Append(input.Expr.Path, operator), operandL, operandR)
		if cmpVal.Status != pb.EvaluateOutput_OK {
			return cmpVal
		}
		return &pb.EvaluateOutput{Value: BoolValue(cmpVal.Value.Num > 0)}
	case "gte":
		cmpVal := compare(Append(input.Expr.Path, operator), operandL, operandR)
		if cmpVal.Status != pb.EvaluateOutput_OK {
			return cmpVal
		}
		return &pb.EvaluateOutput{Value: BoolValue(cmpVal.Value.Num >= 0)}
	}
}

func (e BasicEvaluator) EvaluateOpVariadic(input *pb.EvaluateInput) *pb.EvaluateOutput {
	var (
		operator string
		operands []*pb.Value
	)
	for k, v := range input.Expr.Object { // only one property exists
		operator = k
		for _, elem := range v.Array {
			val := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: elem})
			if val.Status != pb.EvaluateOutput_OK {
				return val
			}
			operands = append(operands, val.Value)
		}
	}
	switch operator {
	default:
		panic(fmt.Sprintf("unexpected variadic operator %q", operator))
	case "add":
		addVal := 0.0
		for i, operand := range operands {
			if operand.Type != pb.Value_NUM {
				return errorUnexpectedType(Append(input.Expr.Path, operator, i), operand.Type, []pb.Value_Type{pb.Value_NUM})
			}
			addVal += operand.Num
		}
		if !isFiniteNumber(addVal) {
			return errorNotFiniteNumber(Append(input.Expr.Path, operator))
		}
		return &pb.EvaluateOutput{Value: NumValue(addVal)}
	case "mul":
		mulVal := 1.0
		for i, operand := range operands {
			if operand.Type != pb.Value_NUM {
				return errorUnexpectedType(Append(input.Expr.Path, operator, i), operand.Type, []pb.Value_Type{pb.Value_NUM})
			}
			mulVal *= operand.Num
		}
		if !isFiniteNumber(mulVal) {
			return errorNotFiniteNumber(Append(input.Expr.Path, operator))
		}
		return &pb.EvaluateOutput{Value: NumValue(mulVal)}
	case "and":
		for i, operand := range operands {
			if operand.Type != pb.Value_BOOL {
				return errorUnexpectedType(Append(input.Expr.Path, operator, i), operand.Type, []pb.Value_Type{pb.Value_BOOL})
			}
			if !operand.Bool {
				return &pb.EvaluateOutput{Value: BoolValue(false)}
			}
		}
		return &pb.EvaluateOutput{Value: BoolValue(true)}
	case "or":
		for i, operand := range operands {
			if operand.Type != pb.Value_BOOL {
				return errorUnexpectedType(Append(input.Expr.Path, operator, i), operand.Type, []pb.Value_Type{pb.Value_BOOL})
			}
			if operand.Bool {
				return &pb.EvaluateOutput{Value: BoolValue(true)}
			}
		}
		return &pb.EvaluateOutput{Value: BoolValue(false)}
	case "cat":
		catVal := ""
		for i, operand := range operands {
			if operand.Type != pb.Value_STR {
				return errorUnexpectedType(Append(input.Expr.Path, operator, i), operand.Type, []pb.Value_Type{pb.Value_STR})
			}
			catVal += operand.Str
		}
		return &pb.EvaluateOutput{Value: StrValue(catVal)}
	case "min":
		minVal := math.Inf(1)
		for i, operand := range operands {
			if operand.Type != pb.Value_NUM {
				return errorUnexpectedType(Append(input.Expr.Path, operator, i), operand.Type, []pb.Value_Type{pb.Value_NUM})
			}
			minVal = math.Min(minVal, operand.Num)
		}
		return &pb.EvaluateOutput{Value: NumValue(minVal)}
	case "max":
		maxVal := math.Inf(-1)
		for i, operand := range operands {
			if operand.Type != pb.Value_NUM {
				return errorUnexpectedType(Append(input.Expr.Path, operator, i), operand.Type, []pb.Value_Type{pb.Value_NUM})
			}
			maxVal = math.Max(maxVal, operand.Num)
		}
		return &pb.EvaluateOutput{Value: NumValue(maxVal)}
	case "merge":
		mergeVal := map[string]*pb.Value{}
		for i, operand := range operands {
			if operand.Type != pb.Value_OBJ {
				return errorUnexpectedType(Append(input.Expr.Path, operator, i), operand.Type, []pb.Value_Type{pb.Value_OBJ})
			}
			for k, v := range operand.Obj {
				mergeVal[k] = v
			}
		}
		return &pb.EvaluateOutput{Value: ObjValue(mergeVal)}
	}
}

func errorIndexOutOfBounds(path *pb.Node_Path, index *pb.Value, begin, end int) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_INVALID_INDEX,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("invalid index: index out of bounds: %v not in [%v, %v)", int64(index.Num), begin, end),
	}
}

func errorIndexNotInteger(path *pb.Node_Path, index float64) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_INVALID_INDEX,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("invalid index: non integer index: %v", index),
	}
}

func errorInvalidKey(path *pb.Node_Path, key string, keys []string) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_INVALID_INDEX,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("invalid key: %q not in {%v}", key, strings.Join(keys, ",")),
	}
}

func errorUnexpectedType(path *pb.Node_Path, got pb.Value_Type, want []pb.Value_Type) *pb.EvaluateOutput {
	wantStr := make([]string, len(want))
	for i, t := range want {
		wantStr[i] = t.String()
	}
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_UNEXPECTED_TYPE,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("unexpected type: got %v, want {%v}", got.String(), strings.Join(wantStr, ",")),
	}
}

func errorArgumentMismatch(path *pb.Node_Path, arg string) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_ARGUMENT_MISMATCH,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("argument mismatch: argument %q required", arg),
	}
}

func errorReferenceNotFound(path *pb.Node_Path, ref string) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_REFERENCE_NOT_FOUND,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("reference not found: %q", ref),
	}
}

func errorCasesNotExhaustive(path *pb.Node_Path) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_CASES_NOT_EXHAUSTIVE,
		ErrorPath:    path,
		ErrorMessage: "cases not exhaustive",
	}
}

func errorNotComparable(path *pb.Node_Path) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_NOT_COMPARABLE,
		ErrorPath:    path,
		ErrorMessage: "not comparable",
	}
}

func errorNotFiniteNumber(path *pb.Node_Path) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_NOT_FINITE_NUMBER,
		ErrorPath:    path,
		ErrorMessage: "not finite number",
	}
}

func canInt(v *pb.Value) bool {
	return v.Type == pb.Value_NUM && v.Num == float64(int(v.Num))
}

func sortedKeys(m map[string]*pb.Value) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func equal(l, r *pb.Value) *pb.Value {
	falseValue, trueValue := BoolValue(false), BoolValue(true)
	switch {
	default:
		panic(fmt.Sprintf("unexpected type %v", l.Type))
	case l.Type != r.Type:
		return falseValue
	case l.Type == pb.Value_NUM:
		return BoolValue(l.Num == r.Num)
	case l.Type == pb.Value_BOOL:
		return BoolValue(l.Bool == r.Bool)
	case l.Type == pb.Value_STR:
		return BoolValue(l.Str == r.Str)
	case l.Type == pb.Value_ARR:
		if len(l.Arr) != len(r.Arr) {
			return falseValue
		}
		for i, l := range l.Arr {
			eq := equal(l, r.Arr[i])
			if !eq.Bool {
				return falseValue
			}
		}
		return trueValue
	case l.Type == pb.Value_OBJ:
		lk, rk := sortedKeys(l.Obj), sortedKeys(r.Obj)
		if !slices.Equal(lk, rk) {
			return falseValue
		}
		for k, l := range l.Obj {
			eq := equal(l, r.Obj[k])
			if !eq.Bool {
				return falseValue
			}
		}
		return trueValue
	}
}
func compare(path *pb.Node_Path, l, r *pb.Value) *pb.EvaluateOutput {
	ltValue := &pb.EvaluateOutput{Value: NumValue(-1)}
	gtValue := &pb.EvaluateOutput{Value: NumValue(1)}
	eqValue := &pb.EvaluateOutput{Value: NumValue(0)}
	switch {
	default:
		return errorNotComparable(path)
	case l.Type == pb.Value_NUM && r.Type == pb.Value_NUM:
		if l.Num < r.Num {
			return ltValue
		}
		if l.Num > r.Num {
			return gtValue
		}
		return eqValue
	case l.Type == pb.Value_BOOL && r.Type == pb.Value_BOOL:
		if !l.Bool && r.Bool {
			return ltValue
		}
		if l.Bool && !r.Bool {
			return gtValue
		}
		return eqValue
	case l.Type == pb.Value_STR && r.Type == pb.Value_STR:
		if l.Str < r.Str {
			return ltValue
		}
		if l.Str > r.Str {
			return gtValue
		}
		return eqValue
	case l.Type == pb.Value_ARR && r.Type == pb.Value_ARR:
		n := len(l.Arr)
		if n > len(r.Arr) {
			n = len(r.Arr)
		}
		for i := 0; i < n; i++ {
			cmp := compare(path, l.Arr[i], r.Arr[i])
			if cmp.Status != pb.EvaluateOutput_OK {
				return cmp
			}
			if cmp.Value.Num != 0 {
				return cmp
			}
		}
		if len(l.Arr) < len(r.Arr) {
			return ltValue
		}
		if len(l.Arr) > len(r.Arr) {
			return gtValue
		}
		return eqValue
	}
}

func isFiniteNumber(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
