package model

import (
	"context"
)

type IStorage interface {
	NewModel(string, []IFieldDefinition, []string) IModel
	RegisterModel(IModel) error
	Add(context.Context, IModel, *Data, AddOptions) (*Data, error)
	Query(context.Context, IModel, []string, GetAllOptions) (*Data, error)
	Edit(context.Context, IModel, IExpression, map[string]interface{}) error
	Delete(context.Context, IModel, IExpression) error
}
