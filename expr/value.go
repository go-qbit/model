package expr

type value struct {
	data interface{}
}

func (e *value) Eval(IModelRow) (interface{}, error) {
	return e.data, nil
}

func Value(v interface{}) *value {
	return &value{v}
}
