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
	GetRelations() []string
	GetRelation(string) *Relation
	AddMulti(context.Context, []string, [][]interface{}, AddOptions) ([]interface{}, error)
	GetAll(context.Context, []string, GetAllOptions) ([]map[string]interface{}, error)
	FieldsToString([]string, map[string]interface{}) string
}

type AddOptions struct {
	Replace bool
}

type GetAllOptions struct {
	Filter  expr.IExpression
	OrderBy []Order
	Limit   uint64
	Offset  uint64
}

type Order struct {
	FieldName string
	Desc      bool
}
