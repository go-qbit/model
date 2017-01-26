package model

import (
	"context"
)

type IStorage interface {
	NewModel(string, []IFieldDefinition, []string) IModel
	RegisterModel(IModel) error
	Add(context.Context, IModel, []string, [][]interface{}, AddOptions) ([]interface{}, error)
	Query(context.Context, IModel, []string, GetAllOptions) ([]map[string]interface{}, error)
}
