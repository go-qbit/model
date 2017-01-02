package expr

import (
	"fmt"
	"reflect"
)

type eq struct {
	op1, op2 IExpression
}

func (e *eq) Eval(row IModelRow) (interface{}, error) {
	v1, err1 := e.op1.Eval(row)
	if err1 != nil {
		return nil, err1
	}

	v2, err2 := e.op2.Eval(row)
	if err2 != nil {
		return nil, err2
	}

	rv1 := reflect.ValueOf(v1)
	rv2 := reflect.ValueOf(v2)
	if rv1.Kind() != rv2.Kind() {
		return nil, fmt.Errorf("%s and %s cannot be compared", rv1.Kind(), rv2.Kind())
	}

	switch rv1.Kind() {
	case reflect.Bool:
		return rv1.Bool() == rv2.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv1.Int() == rv2.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv1.Uint() == rv2.Uint(), nil
	case reflect.String:
		return rv1.String() == rv2.String(), nil
	default:
		return nil, fmt.Errorf("%s and %s cannot be compared", rv1.Kind(), rv2.Kind())
	}
}

func Eq(op1, op2 IExpression) *eq { return &eq{op1, op2} }
