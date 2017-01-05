package test

import (
	"fmt"
	"reflect"

	"github.com/go-qbit/model/expr"
)

type ExprProcessor struct{}

type EvalFunc func(row expr.IModelRow) (interface{}, error)

func (p *ExprProcessor) Eq(op1, op2 expr.IExpression) interface{} {
	return EvalFunc(func(row expr.IModelRow) (interface{}, error) {
		lt, err := p.Lt(op1, op2).(EvalFunc)(row)
		if err != nil {
			return nil, err
		}

		gt, err := p.Lt(op2, op1).(EvalFunc)(row)
		if err != nil {
			return nil, err
		}

		return !(lt.(bool) || gt.(bool)), nil
	})
}

func (p *ExprProcessor) Lt(op1, op2 expr.IExpression) interface{} {
	return EvalFunc(func(row expr.IModelRow) (interface{}, error) {
		v1, err1 := op1.GetProcessor(p).(EvalFunc)(row)
		if err1 != nil {
			return nil, err1
		}

		v2, err2 := op2.GetProcessor(p).(EvalFunc)(row)
		if err2 != nil {
			return nil, err2
		}

		rv1 := reflect.ValueOf(v1)
		rv2 := reflect.ValueOf(v2)
		if rv1.Kind() != rv2.Kind() {
			return nil, fmt.Errorf("%s and %s cannot be compared", rv1.Kind(), rv2.Kind())
		}

		switch rv1.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return rv1.Int() < rv2.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return rv1.Uint() < rv2.Uint(), nil
		case reflect.String:
			return rv1.String() < rv2.String(), nil
		default:
			return nil, fmt.Errorf("%s and %s cannot be compared", rv1.Kind(), rv2.Kind())
		}
	})
}

func (p *ExprProcessor) Le(op1, op2 expr.IExpression) interface{} {
	return EvalFunc(func(row expr.IModelRow) (interface{}, error) {
		return nil, nil
	})
}

func (p *ExprProcessor) Gt(op1, op2 expr.IExpression) interface{} {
	return EvalFunc(func(row expr.IModelRow) (interface{}, error) {
		return p.Lt(op2, op1).(EvalFunc)(row)
	})
}

func (p *ExprProcessor) Ge(op1, op2 expr.IExpression) interface{} {
	return EvalFunc(func(row expr.IModelRow) (interface{}, error) {
		return p.Le(op2, op1).(EvalFunc)(row)
	})
}

func (p *ExprProcessor) In(op expr.IExpression, values []expr.IExpression) interface{} {
	return EvalFunc(func(row expr.IModelRow) (interface{}, error) {
		for _, op2 := range values {
			eq, err := p.Eq(op, op2).(EvalFunc)(row)
			if err != nil {
				return nil, err
			}
			if eq.(bool) {
				return true, nil
			}
		}

		return false, nil
	})
}

func (p *ExprProcessor) ModelField(fieldName string) interface{} {
	return EvalFunc(func(row expr.IModelRow) (interface{}, error) {
		return row.GetValue(fieldName)
	})
}

func (p *ExprProcessor) Value(value interface{}) interface{} {
	return EvalFunc(func(row expr.IModelRow) (interface{}, error) {
		return value, nil
	})
}
