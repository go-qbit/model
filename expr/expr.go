package expr

type IModelRow interface {
	GetValue(string) (interface{}, error)
}

type IExpression interface {
	Eval(IModelRow) (interface{}, error)
}
