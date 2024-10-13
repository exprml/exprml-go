package exprml

import (
	"fmt"
	pb "github.com/exprml/exprml-go/pb/exprml/v1"
	"regexp"
	"strings"
)

type Parser interface {
	Parse(input *pb.ParseInput) *pb.ParseOutput
}

func NewParser() Parser {
	return &parser{}
}

type parser struct{}

func (p parser) Parse(input *pb.ParseInput) *pb.ParseOutput {
	expr, err := parse(&pb.Expr_Path{}, input.Value)
	if err != nil {
		return &pb.ParseOutput{
			IsError:      true,
			ErrorMessage: fmt.Sprintf("fail to parse: %v", err),
		}
	}
	return &pb.ParseOutput{Expr: expr}
}

var identRegexp = regexp.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*$`)

func parse(path *pb.Expr_Path, value *pb.Value) (expr *pb.Expr, err error) {
	expr = &pb.Expr{Path: path, Value: value}
	switch value.Type {
	default:
		return nil, fmt.Errorf("invalid Expr: %v: one of string, number, boolean, or object required but got %v", Format(path), value.Type)
	case pb.Value_STR:
		s := value.Str
		switch {
		default:
			return nil, fmt.Errorf("invalid Scalar: %v: string literal must enclosed by '`'", Format(path))
		case identRegexp.MatchString(s):
			expr.Kind = pb.Expr_REF
			expr.Ref = &pb.Ref{Ident: s}
			return expr, nil
		case len(s) > 1 && strings.HasPrefix(s, "`") && strings.HasSuffix(s, "`"):
			expr.Kind = pb.Expr_SCALAR
			expr.Scalar = &pb.Scalar{Scalar: StrValue(s[1 : len(s)-1])}
			return expr, nil
		}
	case pb.Value_NUM, pb.Value_BOOL:
		expr.Kind = pb.Expr_SCALAR
		expr.Scalar = &pb.Scalar{Scalar: value}
		return expr, nil
	case pb.Value_OBJ:
		switch {
		case hasKey(value, "eval"):
			expr.Kind = pb.Expr_EVAL
			expr.Eval = &pb.Eval{}
			expr.Eval.Eval, err = parse(Append(path, "eval"), value.Obj["eval"])
			if err != nil {
				return nil, err
			}
			if hasKey(value, "where") {
				where := value.Obj["where"]
				if where.Type != pb.Value_ARR {
					return nil, fmt.Errorf("invalid Expr: %v: where clause must be an array", Format(Append(path, "where")))
				}
				for i, def := range where.Arr {
					if def.Type != pb.Value_OBJ {
						return nil, fmt.Errorf("invalid definition: %v: where clause must contain only objects but got %v", Format(Append(path, "where", i)), def.Type)
					}
					keys := Keys(def)
					if len(keys) != 1 {
						return nil, fmt.Errorf("invalid definition: %v: definition must contain one property", Format(Append(path, "where", i)))
					}
					prop := keys[0]
					// `$a`: Match
					// `$a()`: Match
					// `$a($b)`: Match
					// `$a($b, $c)`: Match
					// `$a($b, $c, $d)`: Match
					// `$a($b,)`: Match
					// `$a($b, $c,)`: Match
					// `$a($b, $c, $d,)`: Match
					r := regexp.MustCompile(`^\$[_a-zA-Z][_a-zA-Z0-9]*(\(\s*\)|\(\s*\$[_a-zA-Z][_a-zA-Z0-9]*(\s*,\s*\$[_a-zA-Z][_a-zA-Z0-9]*)*(\s*,)?\s*\))?$`)
					if !r.MatchString(prop) {
						return nil, fmt.Errorf("invalid definition: %v: definition must match %q", Format(Append(path, "where", i, prop)), r.String())
					}
					var idents []string
					for _, r := range []rune(keys[0]) {
						switch {
						case r == '$':
							idents = append(idents, string(r))
						case ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') || r == '_':
							idents[len(idents)-1] += string(r)
						default:
							// skip
						}
					}
					body, err := parse(Append(path, "where", i, prop), def.Obj[prop])
					if err != nil {
						return nil, err
					}
					expr.Eval.Where = append(expr.Eval.Where, &pb.Eval_Definition{
						Ident: idents[0],
						Args:  idents[1:],
						Body:  body,
					})
				}
			}
			return expr, nil
		case hasKey(value, "obj"):
			if value.Obj["obj"].Type != pb.Value_OBJ {
				return nil, fmt.Errorf("invalid Obj: %v: 'obj' property must be an object", Format(Append(path, "obj")))
			}
			expr.Kind = pb.Expr_OBJ
			expr.Obj = &pb.Obj{Obj: map[string]*pb.Expr{}}
			for key, val := range value.Obj["obj"].Obj {
				expr.Obj.Obj[key], err = parse(Append(path, "obj", key), val)
				if err != nil {
					return nil, err
				}
			}
			return expr, nil
		case hasKey(value, "arr"):
			if value.Obj["arr"].Type != pb.Value_ARR {
				return nil, fmt.Errorf("invalid Arr: %v: 'arr' property must be an array", Format(Append(path, "arr")))
			}
			expr.Kind = pb.Expr_ARR
			arr := value.Obj["arr"].Arr
			expr.Arr = &pb.Arr{Arr: make([]*pb.Expr, len(arr))}
			for i, val := range arr {
				expr.Arr.Arr[i], err = parse(Append(path, "arr", fmt.Sprintf("%v", i)), val)
				if err != nil {
					return nil, err
				}
			}
			return expr, nil
		case hasKey(value, "json"):
			expr.Kind = pb.Expr_JSON
			if err := checkNonNull(value.Obj["json"]); err != nil {
				return nil, fmt.Errorf("invalid Json: %v: 'json' property cannot contain null", Format(Append(path, "json")))
			}
			expr.Json = &pb.Json{Json: value.Obj["json"]}
			return expr, nil
		case hasKey(value, "do"):
			expr.Kind = pb.Expr_ITER
			expr.Iter = &pb.Iter{}
			// `for($a, $b)`: Match
			r := regexp.MustCompile(`^for\(\s*\$[_a-zA-Z][_a-zA-Z0-9]*\s*,\s*\$[_a-zA-Z][_a-zA-Z0-9]*\s*\)$`)
			for _, prop := range Keys(value) {
				switch {
				default:
					return nil, fmt.Errorf("invalid Iter: %v: inalid property %q", Format(Append(path, "do", prop)), prop)
				case prop == "do":
					expr.Iter.Do, err = parse(Append(path, "do"), value.Obj["do"])
					if err != nil {
						return nil, err
					}
				case prop == "if":
					expr.Iter.If, err = parse(Append(path, "if"), value.Obj["if"])
					if err != nil {
						return nil, err
					}
				case r.MatchString(prop):
					var idents []string
					for _, r := range []rune(prop[3:]) {
						switch {
						case r == '$':
							idents = append(idents, string(r))
						case ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') || r == '_':
							idents[len(idents)-1] += string(r)
						default:
							// skip
						}
					}
					expr.Iter.PosIdent, expr.Iter.ElemIdent = idents[0], idents[1]
					expr.Iter.Col, err = parse(Append(path, prop), value.Obj[prop])
				}
			}
			if expr.Iter.Col == nil {
				return nil, fmt.Errorf("invalid Iter: %v: 'for(...vars...)' property is required", Format(path))
			}
			return expr, nil
		case hasKey(value, "get"):
			expr.Kind = pb.Expr_ELEM
			expr.Elem = &pb.Elem{}
			expr.Elem.Get, err = parse(Append(path, "get"), value.Obj["get"])
			if err != nil {
				return nil, err
			}
			if !hasKey(value, "from") {
				return nil, fmt.Errorf("invalid Elem: %v: 'from' property is required", Format(path))
			}
			expr.Elem.From, err = parse(Append(path, "from"), value.Obj["from"])
			if err != nil {
				return nil, err
			}
			return expr, nil
		case hasKey(value, "cases"):
			expr.Kind = pb.Expr_CASES
			expr.Cases = &pb.Cases{}
			if value.Obj["cases"].Type != pb.Value_ARR {
				return nil, fmt.Errorf("invalid Cases: %v: 'cases' property must be an array", Format(Append(path, "cases")))
			}
			for i, c := range value.Obj["cases"].Arr {
				if c.Type != pb.Value_OBJ {
					return nil, fmt.Errorf("invalid Case: %v: 'cases' property must contain only objects but got %v", Format(Append(path, "cases", i)), c.Type)
				}
				if _, ok := c.Obj["otherwise"]; ok {
					otherwise, err := parse(Append(path, "cases", i, "otherwise"), c.Obj["otherwise"])
					if err != nil {
						return nil, err
					}
					expr.Cases.Cases = append(expr.Cases.Cases, &pb.Cases_Case{
						IsOtherwise: true,
						Otherwise:   otherwise,
					})
				} else {
					if !hasKey(c, "when") {
						return nil, fmt.Errorf("invalid Case: %v: 'when' property is required", Format(Append(path, "cases", i)))
					}
					when, err := parse(Append(path, "cases", i, "when"), c.Obj["when"])
					if err != nil {
						return nil, err
					}
					if !hasKey(c, "then") {
						return nil, fmt.Errorf("invalid Case: %v: 'then' property is required", Format(Append(path, "cases", i)))
					}
					then, err := parse(Append(path, "cases", i, "then"), c.Obj["then"])
					if err != nil {
						return nil, err
					}
					expr.Cases.Cases = append(expr.Cases.Cases, &pb.Cases_Case{
						When: when,
						Then: then,
					})
				}
			}
			return expr, nil
		default:
			if len(value.Obj) != 1 {
				return nil, fmt.Errorf("invalid Expr: %v: operation or function call must contain only one property", Format(path))
			}
			prop := Keys(value)[0]

			expr.OpUnary = map[string]*pb.OpUnary{
				"len":   {Op: pb.OpUnary_LEN},
				"not":   {Op: pb.OpUnary_NOT},
				"flat":  {Op: pb.OpUnary_FLAT},
				"floor": {Op: pb.OpUnary_FLOOR},
				"ceil":  {Op: pb.OpUnary_CEIL},
				"abort": {Op: pb.OpUnary_ABORT},
			}[prop]
			if expr.OpUnary != nil {
				expr.Kind = pb.Expr_OP_UNARY
				expr.OpUnary.Operand, err = parse(Append(path, prop), value.Obj[prop])
				if err != nil {
					return nil, err
				}
				return expr, nil
			}

			expr.OpBinary = map[string]*pb.OpBinary{
				"sub": {Op: pb.OpBinary_SUB},
				"div": {Op: pb.OpBinary_DIV},
				"eq":  {Op: pb.OpBinary_EQ},
				"neq": {Op: pb.OpBinary_NEQ},
				"lt":  {Op: pb.OpBinary_LT},
				"lte": {Op: pb.OpBinary_LTE},
				"gt":  {Op: pb.OpBinary_GT},
				"gte": {Op: pb.OpBinary_GTE},
			}[prop]
			if expr.OpBinary != nil {
				expr.Kind = pb.Expr_OP_BINARY
				if value.Obj[prop].Type != pb.Value_ARR {
					return nil, fmt.Errorf("invalid OpBinary: %v: '%v' property must be an array", Format(Append(path, prop)), prop)
				}
				if len(value.Obj[prop].Arr) != 2 {
					return nil, fmt.Errorf("invalid OpBinary: %v: '%v' property must contain two elements", Format(Append(path, prop)), prop)
				}
				expr.OpBinary.Left, err = parse(Append(path, prop), value.Obj[prop].Arr[0])
				if err != nil {
					return nil, err
				}
				expr.OpBinary.Right, err = parse(Append(path, prop), value.Obj[prop].Arr[1])
				if err != nil {
					return nil, err
				}
				return expr, nil
			}

			expr.OpVariadic = map[string]*pb.OpVariadic{
				"add":   {Op: pb.OpVariadic_ADD},
				"mul":   {Op: pb.OpVariadic_MUL},
				"and":   {Op: pb.OpVariadic_AND},
				"or":    {Op: pb.OpVariadic_OR},
				"cat":   {Op: pb.OpVariadic_CAT},
				"min":   {Op: pb.OpVariadic_MIN},
				"max":   {Op: pb.OpVariadic_MAX},
				"merge": {Op: pb.OpVariadic_MERGE},
			}[prop]
			if expr.OpVariadic != nil {
				expr.Kind = pb.Expr_OP_VARIADIC
				if value.Obj[prop].Type != pb.Value_ARR {
					return nil, fmt.Errorf("invalid OpVariadic: %v: '%v' property must be an array", Format(Append(path, prop)), prop)
				}
				if (prop == "min" || prop == "max") && len(value.Obj[prop].Arr) == 0 {
					return nil, fmt.Errorf("invalid OpVariadic: %v: '%v' property must contain at least one element", Format(Append(path, prop)), prop)
				}
				for i, v := range value.Obj[prop].Arr {
					operand, err := parse(Append(path, prop, i), v)
					if err != nil {
						return nil, err
					}
					expr.OpVariadic.Operands = append(expr.OpVariadic.Operands, operand)
				}
				return expr, nil
			}
			if !identRegexp.MatchString(prop) {
				return nil, fmt.Errorf("invalid Call: %v: function call property %q must match %q", Format(path), prop, identRegexp.String())
			}
			args := value.Obj[prop]
			if args.Type != pb.Value_OBJ {
				return nil, fmt.Errorf("invalid Call: %v: arguments must be given as an object", Format(Append(path, prop)))
			}

			expr.Kind = pb.Expr_CALL
			expr.Call = &pb.Call{Ident: prop}
			expr.Call.Args = map[string]*pb.Expr{}
			for key, val := range args.Obj {
				if !identRegexp.MatchString(key) {
					return nil, fmt.Errorf("invalid Call: %v: argument property %q must match %q", Format(Append(path, prop, key)), key, identRegexp.String())
				}
				expr.Call.Args[key], err = parse(Append(path, prop, key), val)
				if err != nil {
					return nil, err
				}
			}
			return expr, nil
		}
	}
}

func hasKey(v *pb.Value, key string) bool {
	if v.Type != pb.Value_OBJ {
		return false
	}
	_, ok := v.Obj[key]
	return ok
}

func checkNonNull(value *pb.Value) error {
	if value.Type == pb.Value_NULL {
		return fmt.Errorf("null value is not allowed")
	}
	for _, v := range value.Obj {
		if err := checkNonNull(v); err != nil {
			return err
		}
	}
	for _, v := range value.Arr {
		if err := checkNonNull(v); err != nil {
			return err
		}
	}
	return nil
}
