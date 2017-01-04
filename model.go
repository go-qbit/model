package model

import (
	"context"

	"github.com/go-qbit/model/expr"
)

type IModel interface {
	GetId() string
	GetPKFieldsNames() []string
	GetFieldsNames() []string
	GetFieldDefinition(string) IFieldDefinition
	//AddField(IFieldDefinition)
	AddRelation(Relation, []IFieldDefinition)
	AddMulti(context.Context, []map[string]interface{}) ([]interface{}, error)
	GetAll(context.Context, []string, GetAllOptions) ([]map[string]interface{}, error)
	FieldsToString([]string, map[string]interface{}) string
}

type GetAllOptions struct {
	Filter expr.IExpression
	Limit  uint64
	Offset uint64
}
