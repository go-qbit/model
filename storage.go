package model

import (
	"context"

	"github.com/go-qbit/model/expr"
)

type IStorage interface {
	RegisterModel(IModel) error
	Add(context.Context, IModel, []map[string]interface{}) ([]interface{}, error)
	Query(context.Context, IModel, []string, QueryOptions) ([]map[string]interface{}, error)
}

type QueryOptions struct {
	Filter  expr.IExpression
	OrderBy []Order
	Limit   uint64
	Offset  uint64
}
