package expr

import "github.com/go-qbit/model"

// Eq
type eq struct {
	op1, op2 model.IExpression
}

func Eq(op1, op2 model.IExpression) *eq { return &eq{op1, op2} }

func (e *eq) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Eq(e.op1, e.op2)
}

// Ne
type ne struct {
	op1, op2 model.IExpression
}

func Ne(op1, op2 model.IExpression) *ne { return &ne{op1, op2} }

func (e *ne) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Ne(e.op1, e.op2)
}

// Lt
type lt struct {
	op1, op2 model.IExpression
}

func Lt(op1, op2 model.IExpression) *lt { return &lt{op1, op2} }
func (e *lt) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Lt(e.op1, e.op2)
}

// Le
type le struct {
	op1, op2 model.IExpression
}

func Le(op1, op2 model.IExpression) *le { return &le{op1, op2} }
func (e *le) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Le(e.op1, e.op2)
}

// Gt
type gt struct {
	op1, op2 model.IExpression
}

func Gt(op1, op2 model.IExpression) *gt { return &gt{op1, op2} }
func (e *gt) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Gt(e.op1, e.op2)
}

// Ge
type ge struct {
	op1, op2 model.IExpression
}

func Ge(op1, op2 model.IExpression) *ge { return &ge{op1, op2} }
func (e *ge) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Ge(e.op1, e.op2)
}

// In
type in struct {
	op     model.IExpression
	values []model.IExpression
}

func In(op model.IExpression) *in         { return &in{op, nil} }
func (e *in) Add(value model.IExpression) { e.values = append(e.values, value) }
func (e *in) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.In(e.op, e.values)
}

// And
type and struct {
	ops []model.IExpression
}

func And(op1, op2 model.IExpression, extraOps ...model.IExpression) *and {
	return &and{
		ops: append([]model.IExpression{op1, op2}, extraOps...),
	}
}
func (e *and) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.And(e.ops)
}

// Or
type or struct {
	ops []model.IExpression
}

func Or(op1, op2 model.IExpression, extraOps ...model.IExpression) *or {
	return &or{
		ops: append([]model.IExpression{op1, op2}, extraOps...),
	}
}
func (e *or) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Or(e.ops)
}

// Any
type any struct {
	localModel model.IModel
	extModel   model.IModel
	filter     model.IExpression
}

func Any(localModel, extModel model.IModel, filter model.IExpression) *any {
	return &any{localModel, extModel, filter}
}

func (e *any) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Any(e.localModel, e.extModel, e.filter)
}

// Model field
type modelField struct {
	m     model.IModel
	field string
}

func ModelField(m model.IModel, field string) *modelField { return &modelField{m, field} }
func (e modelField) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.ModelField(e.m, e.field)
}

// Value
type value struct {
	data interface{}
}

func Value(v interface{}) *value { return &value{v} }
func (e *value) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Value(e.data)
}

// Func
type function struct {
	name   string
	params []model.IExpression
}

func Func(name string, params ...model.IExpression) *function { return &function{name, params} }
func (e *function) GetProcessor(processor model.IExpressionProcessor) interface{} {
	return processor.Func(e.name, e.params...)
}
