package expr

type IModelRow interface {
	GetValue(string) (interface{}, error)
}

type IExpression interface {
	GetProcessor(IExpressionProcessor) interface{}
}

type IExpressionProcessor interface {
	Eq(op1, op2 IExpression) interface{}
	Lt(op1, op2 IExpression) interface{}
	Le(op1, op2 IExpression) interface{}
	Gt(op1, op2 IExpression) interface{}
	Ge(op1, op2 IExpression) interface{}
	In(op IExpression, arr []IExpression) interface{}
	And(operators []IExpression) interface{}
	Or(operators []IExpression) interface{}
	ModelField(fieldName string) interface{}
	Value(value interface{}) interface{}
}

// Eq
type eq struct {
	op1, op2 IExpression
}

func Eq(op1, op2 IExpression) *eq { return &eq{op1, op2} }

func (e *eq) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.Eq(e.op1, e.op2)
}

// Lt
type lt struct {
	op1, op2 IExpression
}

func Lt(op1, op2 IExpression) *lt { return &lt{op1, op2} }
func (e *lt) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.Lt(e.op1, e.op2)
}

// Le
type le struct {
	op1, op2 IExpression
}

func Le(op1, op2 IExpression) *le { return &le{op1, op2} }
func (e *le) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.Le(e.op1, e.op2)
}

// Gt
type gt struct {
	op1, op2 IExpression
}

func Gt(op1, op2 IExpression) *gt { return &gt{op1, op2} }
func (e *gt) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.Gt(e.op1, e.op2)
}

// Ge
type ge struct {
	op1, op2 IExpression
}

func Ge(op1, op2 IExpression) *ge { return &ge{op1, op2} }
func (e *ge) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.Ge(e.op1, e.op2)
}

// In
type in struct {
	op     IExpression
	values []IExpression
}

func In(op IExpression) *in         { return &in{op, nil} }
func (e *in) Add(value IExpression) { e.values = append(e.values, value) }
func (e *in) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.In(e.op, e.values)
}

// And
type and struct {
	ops []IExpression
}

func And(op1, op2 IExpression, extraOps ...IExpression) *and {
	return &and{
		ops: append([]IExpression{op1, op2}, extraOps...),
	}
}
func (e *and) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.And(e.ops)
}

// Or
type or struct {
	ops []IExpression
}

func Or(op1, op2 IExpression, extraOps ...IExpression) *or {
	return &or{
		ops: append([]IExpression{op1, op2}, extraOps...),
	}
}
func (e *or) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.Or(e.ops)
}

// Model field
type ModelField string

func (e ModelField) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.ModelField(string(e))
}

// Value
type value struct {
	data interface{}
}

func Value(v interface{}) *value { return &value{v} }
func (e *value) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.Value(e.data)
}
