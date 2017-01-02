package expr

import (
	"fmt"
	"reflect"
)

type in struct {
	op     IExpression
	values []IExpression
}

func (e *in) Add(value IExpression) {
	e.values = append(e.values, value)
}

func (e *in) Eval(row IModelRow) (interface{}, error) {
	v1, err1 := e.op.Eval(row)
	if err1 != nil {
		return nil, err1
	}

	rv1 := reflect.ValueOf(v1)

	for _, op2 := range e.values {
		v2, err2 := op2.Eval(row)
		if err2 != nil {
			return nil, err2
		}

		rv2 := reflect.ValueOf(v2)
		if rv1.Kind() != rv2.Kind() {
			return nil, fmt.Errorf("%s and %s cannot be compared", rv1.Kind(), rv2.Kind())
		}

		switch rv1.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if rv1.Int() == rv2.Int() {
				return true, nil
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if rv1.Uint() == rv2.Uint() {
				return true, nil
			}
		case reflect.String:
			if rv1.String() == rv2.String() {
				return true, nil
			}
		default:
			return nil, fmt.Errorf("%s and %s cannot be compared", rv1.Kind(), rv2.Kind())
		}
	}

	return false, nil
}

func In(op IExpression) *in { return &in{op, nil} }
