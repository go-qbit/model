package expr

type ModelField string

func (e ModelField) Eval(row IModelRow) (interface{}, error) {
	return row.GetValue(string(e))
}
