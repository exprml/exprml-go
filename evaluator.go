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
	EvaluateRef(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateObj(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateArr(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateJson(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateIter(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateElem(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateCall(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateCases(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateOpUnary(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateOpBinary(input *pb.EvaluateInput) *pb.EvaluateOutput
	EvaluateOpVariadic(input *pb.EvaluateInput) *pb.EvaluateOutput
}

func NewEvaluator() Evaluator {
	return &BasicEvaluator{}
}

type BasicEvaluator struct{}

func (e BasicEvaluator) EvaluateExpr(input *pb.EvaluateInput) *pb.EvaluateOutput {
	n := input.Expr
	switch n.Kind {
	default:
		panic("given expression must be validated")
	case pb.Expr_EVAL:
		return e.EvaluateEval(input)
	case pb.Expr_SCALAR:
		return e.EvaluateScalar(input)
	case pb.Expr_REF:
		return e.EvaluateRef(input)
	case pb.Expr_OBJ:
		return e.EvaluateObj(input)
	case pb.Expr_ARR:
		return e.EvaluateArr(input)
	case pb.Expr_JSON:
		return e.EvaluateJson(input)
	case pb.Expr_ITER:
		return e.EvaluateIter(input)
	case pb.Expr_ELEM:
		return e.EvaluateElem(input)
	case pb.Expr_CALL:
		return e.EvaluateCall(input)
	case pb.Expr_CASES:
		return e.EvaluateCases(input)
	case pb.Expr_OP_UNARY:
		return e.EvaluateOpUnary(input)
	case pb.Expr_OP_BINARY:
		return e.EvaluateOpBinary(input)
	case pb.Expr_OP_VARIADIC:
		return e.EvaluateOpVariadic(input)
	}
}

func (e BasicEvaluator) EvaluateEval(input *pb.EvaluateInput) *pb.EvaluateOutput {
	st := input.DefStack
	where := input.Expr.Eval.Where
	for _, def := range where {
		st = Register(st, def)
	}
	return e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: input.Expr.Eval.Eval})
}

func (e BasicEvaluator) EvaluateScalar(input *pb.EvaluateInput) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{Value: input.Expr.Scalar.Scalar}
}

func (e BasicEvaluator) EvaluateObj(input *pb.EvaluateInput) *pb.EvaluateOutput {
	result := map[string]*pb.Value{}
	for pos, expr := range input.Expr.Obj.Obj {
		val := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: expr})
		if val.Status != pb.EvaluateOutput_OK {
			return val
		}
		result[pos] = val.Value
	}
	return &pb.EvaluateOutput{Value: ObjValue(result)}
}

func (e BasicEvaluator) EvaluateArr(input *pb.EvaluateInput) *pb.EvaluateOutput {
	result := []*pb.Value{}
	for _, expr := range input.Expr.Arr.Arr {
		val := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: expr})
		if val.Status != pb.EvaluateOutput_OK {
			return val
		}
		result = append(result, val.Value)
	}
	return &pb.EvaluateOutput{Value: ArrValue(result)}
}

func (e BasicEvaluator) EvaluateJson(input *pb.EvaluateInput) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{Value: input.Expr.Json.Json}
}

func (e BasicEvaluator) EvaluateIter(input *pb.EvaluateInput) *pb.EvaluateOutput {
	iter := input.Expr.Iter
	forPos, forElem := iter.PosIdent, iter.ElemIdent
	inVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: iter.Col})
	switch inVal.Value.Type {
	default:
		return errorUnexpectedType(iter.Col.Path, inVal.Value.Type, []pb.Value_Type{pb.Value_ARR, pb.Value_OBJ})
	case pb.Value_STR:
		result := []*pb.Value{}
		for i, c := range []rune(inVal.Value.Str) {
			st := input.DefStack
			st = Register(st, NewDefinition(input.Expr.Path, forPos, NumValue(float64(i))))
			st = Register(st, NewDefinition(input.Expr.Path, forElem, StrValue(string(c))))
			if iter.If != nil {
				ifVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: iter.If})
				if ifVal.Status != pb.EvaluateOutput_OK {
					return ifVal
				}
				if ifVal.Value.Type != pb.Value_BOOL {
					return errorUnexpectedType(iter.If.Path, ifVal.Value.Type, []pb.Value_Type{pb.Value_BOOL})
				}
				if !ifVal.Value.Bool {
					continue
				}
			}
			v := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: iter.Do})
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
			st = Register(st, NewDefinition(input.Expr.Path, forPos, NumValue(float64(i))))
			st = Register(st, NewDefinition(input.Expr.Path, forElem, elemVal))
			if iter.If != nil {
				ifVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: iter.If})
				if ifVal.Status != pb.EvaluateOutput_OK {
					return ifVal
				}
				if ifVal.Value.Type != pb.Value_BOOL {
					return errorUnexpectedType(iter.If.Path, ifVal.Value.Type, []pb.Value_Type{pb.Value_BOOL})
				}
				if !ifVal.Value.Bool {
					continue
				}
			}
			v := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: iter.Do})
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
			st = Register(st, NewDefinition(input.Expr.Path, forPos, StrValue(key)))
			st = Register(st, NewDefinition(input.Expr.Path, forElem, inVal.Value.Obj[key]))
			if iter.If != nil {
				ifVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: iter.If})
				if ifVal.Status != pb.EvaluateOutput_OK {
					return ifVal
				}
				if ifVal.Value.Type != pb.Value_BOOL {
					return errorUnexpectedType(iter.If.Path, ifVal.Value.Type, []pb.Value_Type{pb.Value_BOOL})
				}
				if !ifVal.Value.Bool {
					continue
				}
			}
			v := e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: iter.Do})
			if v.Status != pb.EvaluateOutput_OK {
				return v
			}
			result[key] = v.Value
		}
		return &pb.EvaluateOutput{Value: ObjValue(result)}
	}
}

func (e BasicEvaluator) EvaluateElem(input *pb.EvaluateInput) *pb.EvaluateOutput {
	elem := input.Expr.Elem
	getVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: elem.Get})
	if getVal.Status != pb.EvaluateOutput_OK {
		return getVal
	}
	pos := getVal.Value
	fromVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: elem.From})
	if fromVal.Status != pb.EvaluateOutput_OK {
		return fromVal
	}
	col := fromVal.Value

	switch col.Type {
	default:
		return errorUnexpectedType(elem.From.Path, col.Type, []pb.Value_Type{pb.Value_STR, pb.Value_ARR, pb.Value_OBJ})
	case pb.Value_STR:
		if pos.Type != pb.Value_NUM {
			return errorUnexpectedType(elem.Get.Path, pos.Type, []pb.Value_Type{pb.Value_NUM})
		}
		if !canInt(pos) {
			return errorIndexNotInteger(elem.Get.Path, pos.Num)
		}
		idx := int(pos.Num)
		if idx < 0 || idx >= len([]rune(col.Str)) {
			return errorIndexOutOfBounds(elem.Get.Path, pos, 0, len(elem.From.Value.Arr))
		}
		return &pb.EvaluateOutput{Value: StrValue(string([]rune(col.Str)[idx]))}
	case pb.Value_ARR:
		if pos.Type != pb.Value_NUM {
			return errorUnexpectedType(elem.Get.Path, pos.Type, []pb.Value_Type{pb.Value_NUM})
		}
		if !canInt(pos) {
			return errorIndexNotInteger(elem.Get.Path, pos.Num)
		}
		idx := int(pos.Num)
		if idx < 0 || idx >= len(col.Arr) {
			return errorIndexOutOfBounds(elem.Get.Path, pos, 0, len(col.Arr))
		}
		return &pb.EvaluateOutput{Value: col.Arr[idx]}
	case pb.Value_OBJ:
		if pos.Type != pb.Value_STR {
			return errorUnexpectedType(elem.Get.Path, pos.Type, []pb.Value_Type{pb.Value_STR})
		}
		key := pos.Str
		if _, ok := col.Obj[key]; !ok {
			return errorInvalidKey(elem.Get.Path, key, sortedKeys(col.Obj))
		}
		return &pb.EvaluateOutput{Value: col.Obj[key]}
	}
}

func (e BasicEvaluator) EvaluateCall(input *pb.EvaluateInput) *pb.EvaluateOutput {
	call := input.Expr.Call
	st := Find(input.DefStack, call.Ident)
	if st == nil {
		return errorReferenceNotFound(input.Expr.Path, call.Ident)
	}
	def := st.Def
	for _, argName := range def.Args {
		arg, ok := call.Args[argName]
		if !ok {
			return errorArgumentMismatch(input.Expr.Path, argName)
		}
		argVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: arg})
		if argVal.Status != pb.EvaluateOutput_OK {
			return argVal
		}
		st = Register(st, NewDefinition(Append(input.Expr.Path, call.Ident, argName), argName, argVal.Value))
	}
	return e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: def.Body})
}

func (e BasicEvaluator) EvaluateCases(input *pb.EvaluateInput) *pb.EvaluateOutput {
	cases := input.Expr.Cases.Cases
	for _, case_ := range cases {
		if case_.IsOtherwise {
			return e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: case_.Otherwise})
		} else {
			boolVal := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: case_.When})
			if boolVal.Status != pb.EvaluateOutput_OK {
				return boolVal
			}
			if boolVal.Value.Type != pb.Value_BOOL {
				return errorUnexpectedType(case_.When.Path, boolVal.Value.Type, []pb.Value_Type{pb.Value_BOOL})
			}
			if boolVal.Value.Bool {
				return e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: case_.Then})
			}
		}
	}
	return errorCasesNotExhaustive(Append(input.Expr.Path, "cases"))
}

func (e BasicEvaluator) EvaluateOpUnary(input *pb.EvaluateInput) *pb.EvaluateOutput {
	op := input.Expr.OpUnary
	o := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: op.Operand})
	if o.Status != pb.EvaluateOutput_OK {
		return o
	}
	operand := o.Value
	switch op.Op {
	default:
		panic(fmt.Sprintf("unexpected unary operator %q", op.Op.String()))
	case pb.OpUnary_LEN:
		switch operand.Type {
		default:
			return errorUnexpectedType(Append(input.Expr.Path, "len"), operand.Type, []pb.Value_Type{pb.Value_STR, pb.Value_ARR, pb.Value_OBJ})
		case pb.Value_STR:
			return &pb.EvaluateOutput{Value: NumValue(float64(len([]rune(operand.Str))))}
		case pb.Value_ARR:
			return &pb.EvaluateOutput{Value: NumValue(float64(len(operand.Arr)))}
		case pb.Value_OBJ:
			return &pb.EvaluateOutput{Value: NumValue(float64(len(operand.Obj)))}
		}
	case pb.OpUnary_NOT:
		if operand.Type != pb.Value_BOOL {
			return errorUnexpectedType(Append(input.Expr.Path, "not"), operand.Type, []pb.Value_Type{pb.Value_BOOL})
		}
		return &pb.EvaluateOutput{Value: BoolValue(!operand.Bool)}
	case pb.OpUnary_FLAT:
		if operand.Type != pb.Value_ARR {
			return errorUnexpectedType(Append(input.Expr.Path, "flat"), operand.Type, []pb.Value_Type{pb.Value_ARR})
		}
		v := []*pb.Value{}
		for _, elem := range operand.Arr {
			if elem.Type != pb.Value_ARR {
				return errorUnexpectedType(Append(input.Expr.Path, "flat"), elem.Type, []pb.Value_Type{pb.Value_ARR})
			}
			v = append(v, elem.Arr...)
		}
		return &pb.EvaluateOutput{Value: ArrValue(v)}
	case pb.OpUnary_FLOOR:
		if operand.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, "floor"), operand.Type, []pb.Value_Type{pb.Value_NUM})
		}
		return &pb.EvaluateOutput{Value: NumValue(math.Floor(operand.Num))}
	case pb.OpUnary_CEIL:
		if operand.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, "ceil"), operand.Type, []pb.Value_Type{pb.Value_NUM})
		}
		return &pb.EvaluateOutput{Value: NumValue(math.Ceil(operand.Num))}
	case pb.OpUnary_ABORT:
		if operand.Type != pb.Value_STR {
			return errorUnexpectedType(Append(input.Expr.Path, "abort"), operand.Type, []pb.Value_Type{pb.Value_STR})
		}
		return &pb.EvaluateOutput{Status: pb.EvaluateOutput_ABORTED, ErrorMessage: operand.Str}
	}
}

func (e BasicEvaluator) EvaluateOpBinary(input *pb.EvaluateInput) *pb.EvaluateOutput {
	op := input.Expr.OpBinary
	ol := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: op.Left})
	if ol.Status != pb.EvaluateOutput_OK {
		return ol
	}
	or := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: op.Right})
	if or.Status != pb.EvaluateOutput_OK {
		return or
	}
	operandL, operandR := ol.Value, or.Value
	switch op.Op {
	default:
		panic(fmt.Sprintf("unexpected binary operator %q", op.Op.String()))
	case pb.OpBinary_SUB:
		if operandL.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, "sub", 0), operandL.Type, []pb.Value_Type{pb.Value_NUM})
		}
		if operandR.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, "sub", 1), operandR.Type, []pb.Value_Type{pb.Value_NUM})
		}
		v := operandL.Num - operandR.Num
		if !isFiniteNumber(v) {
			return errorNotFiniteNumber(input.Expr.Path)
		}
		return &pb.EvaluateOutput{Value: NumValue(v)}
	case pb.OpBinary_DIV:
		if operandL.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, "div", 0), operandL.Type, []pb.Value_Type{pb.Value_NUM})
		}
		if operandR.Type != pb.Value_NUM {
			return errorUnexpectedType(Append(input.Expr.Path, "div", 1), operandR.Type, []pb.Value_Type{pb.Value_NUM})
		}
		v := operandL.Num / operandR.Num
		if !isFiniteNumber(v) {
			return errorNotFiniteNumber(input.Expr.Path)
		}
		return &pb.EvaluateOutput{Value: NumValue(v)}
	case pb.OpBinary_EQ:
		return &pb.EvaluateOutput{Value: equal(operandL, operandR)}
	case pb.OpBinary_NEQ:
		return &pb.EvaluateOutput{Value: BoolValue(!equal(operandL, operandR).Bool)}
	case pb.OpBinary_LT:
		cmpVal := compare(Append(input.Expr.Path, "lt"), operandL, operandR)
		if cmpVal.Status != pb.EvaluateOutput_OK {
			return cmpVal
		}
		return &pb.EvaluateOutput{Value: BoolValue(cmpVal.Value.Num < 0)}
	case pb.OpBinary_LTE:
		cmpVal := compare(Append(input.Expr.Path, "lte"), operandL, operandR)
		if cmpVal.Status != pb.EvaluateOutput_OK {
			return cmpVal
		}
		return &pb.EvaluateOutput{Value: BoolValue(cmpVal.Value.Num <= 0)}
	case pb.OpBinary_GT:
		cmpVal := compare(Append(input.Expr.Path, "gt"), operandL, operandR)
		if cmpVal.Status != pb.EvaluateOutput_OK {
			return cmpVal
		}
		return &pb.EvaluateOutput{Value: BoolValue(cmpVal.Value.Num > 0)}
	case pb.OpBinary_GTE:
		cmpVal := compare(Append(input.Expr.Path, "gte"), operandL, operandR)
		if cmpVal.Status != pb.EvaluateOutput_OK {
			return cmpVal
		}
		return &pb.EvaluateOutput{Value: BoolValue(cmpVal.Value.Num >= 0)}
	}
}

func (e BasicEvaluator) EvaluateOpVariadic(input *pb.EvaluateInput) *pb.EvaluateOutput {
	op := input.Expr.OpVariadic
	operands := []*pb.Value{}
	for _, elem := range op.Operands {
		val := e.EvaluateExpr(&pb.EvaluateInput{DefStack: input.DefStack, Expr: elem})
		if val.Status != pb.EvaluateOutput_OK {
			return val
		}
		operands = append(operands, val.Value)
	}
	switch op.Op {
	default:
		panic(fmt.Sprintf("unexpected variadic operator %q", op.Op.String()))
	case pb.OpVariadic_ADD:
		addVal := 0.0
		for i, operand := range operands {
			if operand.Type != pb.Value_NUM {
				return errorUnexpectedType(Append(input.Expr.Path, "add", i), operand.Type, []pb.Value_Type{pb.Value_NUM})
			}
			addVal += operand.Num
		}
		if !isFiniteNumber(addVal) {
			return errorNotFiniteNumber(Append(input.Expr.Path, "add"))
		}
		return &pb.EvaluateOutput{Value: NumValue(addVal)}
	case pb.OpVariadic_MUL:
		mulVal := 1.0
		for i, operand := range operands {
			if operand.Type != pb.Value_NUM {
				return errorUnexpectedType(Append(input.Expr.Path, "mul", i), operand.Type, []pb.Value_Type{pb.Value_NUM})
			}
			mulVal *= operand.Num
		}
		if !isFiniteNumber(mulVal) {
			return errorNotFiniteNumber(Append(input.Expr.Path, "mul"))
		}
		return &pb.EvaluateOutput{Value: NumValue(mulVal)}
	case pb.OpVariadic_AND:
		for i, operand := range operands {
			if operand.Type != pb.Value_BOOL {
				return errorUnexpectedType(Append(input.Expr.Path, "and", i), operand.Type, []pb.Value_Type{pb.Value_BOOL})
			}
			if !operand.Bool {
				return &pb.EvaluateOutput{Value: BoolValue(false)}
			}
		}
		return &pb.EvaluateOutput{Value: BoolValue(true)}
	case pb.OpVariadic_OR:
		for i, operand := range operands {
			if operand.Type != pb.Value_BOOL {
				return errorUnexpectedType(Append(input.Expr.Path, "or", i), operand.Type, []pb.Value_Type{pb.Value_BOOL})
			}
			if operand.Bool {
				return &pb.EvaluateOutput{Value: BoolValue(true)}
			}
		}
		return &pb.EvaluateOutput{Value: BoolValue(false)}
	case pb.OpVariadic_CAT:
		catVal := ""
		for i, operand := range operands {
			if operand.Type != pb.Value_STR {
				return errorUnexpectedType(Append(input.Expr.Path, "cat", i), operand.Type, []pb.Value_Type{pb.Value_STR})
			}
			catVal += operand.Str
		}
		return &pb.EvaluateOutput{Value: StrValue(catVal)}
	case pb.OpVariadic_MIN:
		minVal := math.Inf(1)
		for i, operand := range operands {
			if operand.Type != pb.Value_NUM {
				return errorUnexpectedType(Append(input.Expr.Path, "min", i), operand.Type, []pb.Value_Type{pb.Value_NUM})
			}
			minVal = math.Min(minVal, operand.Num)
		}
		return &pb.EvaluateOutput{Value: NumValue(minVal)}
	case pb.OpVariadic_MAX:
		maxVal := math.Inf(-1)
		for i, operand := range operands {
			if operand.Type != pb.Value_NUM {
				return errorUnexpectedType(Append(input.Expr.Path, "max", i), operand.Type, []pb.Value_Type{pb.Value_NUM})
			}
			maxVal = math.Max(maxVal, operand.Num)
		}
		return &pb.EvaluateOutput{Value: NumValue(maxVal)}
	case pb.OpVariadic_MERGE:
		mergeVal := map[string]*pb.Value{}
		for i, operand := range operands {
			if operand.Type != pb.Value_OBJ {
				return errorUnexpectedType(Append(input.Expr.Path, "merge", i), operand.Type, []pb.Value_Type{pb.Value_OBJ})
			}
			for k, v := range operand.Obj {
				mergeVal[k] = v
			}
		}
		return &pb.EvaluateOutput{Value: ObjValue(mergeVal)}
	}
}

func (e BasicEvaluator) EvaluateRef(input *pb.EvaluateInput) *pb.EvaluateOutput {
	ref := input.Expr.Ref
	st := Find(input.DefStack, ref.Ident)
	if st == nil {
		return errorReferenceNotFound(input.Expr.Path, ref.Ident)
	}
	return e.EvaluateExpr(&pb.EvaluateInput{DefStack: st, Expr: st.Def.Body})
}

func errorIndexOutOfBounds(path *pb.Expr_Path, index *pb.Value, begin, end int) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_INVALID_INDEX,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("invalid index: index out of bounds: %v not in [%v, %v)", int64(index.Num), begin, end),
	}
}

func errorIndexNotInteger(path *pb.Expr_Path, index float64) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_INVALID_INDEX,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("invalid index: non integer index: %v", index),
	}
}

func errorInvalidKey(path *pb.Expr_Path, key string, keys []string) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_INVALID_INDEX,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("invalid key: %q not in {%v}", key, strings.Join(keys, ",")),
	}
}

func errorUnexpectedType(path *pb.Expr_Path, got pb.Value_Type, want []pb.Value_Type) *pb.EvaluateOutput {
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

func errorArgumentMismatch(path *pb.Expr_Path, arg string) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_ARGUMENT_MISMATCH,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("argument mismatch: argument %q required", arg),
	}
}

func errorReferenceNotFound(path *pb.Expr_Path, ref string) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_REFERENCE_NOT_FOUND,
		ErrorPath:    path,
		ErrorMessage: fmt.Sprintf("reference not found: %q", ref),
	}
}

func errorCasesNotExhaustive(path *pb.Expr_Path) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_CASES_NOT_EXHAUSTIVE,
		ErrorPath:    path,
		ErrorMessage: "cases not exhaustive",
	}
}

func errorNotComparable(path *pb.Expr_Path) *pb.EvaluateOutput {
	return &pb.EvaluateOutput{
		Status:       pb.EvaluateOutput_NOT_COMPARABLE,
		ErrorPath:    path,
		ErrorMessage: "not comparable",
	}
}

func errorNotFiniteNumber(path *pb.Expr_Path) *pb.EvaluateOutput {
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
func compare(path *pb.Expr_Path, l, r *pb.Value) *pb.EvaluateOutput {
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
