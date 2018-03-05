package model

import (
	"context"
)

type IModel interface {
	GetId() string
	GetPKFieldsNames() []string
	GetFieldsNames() []string
	GetFieldDefinition(string) IFieldDefinition
	GetSharedData(string) interface{}
	FieldExpr(string) *exprModelFieldS
	//AddField(IFieldDefinition)
	AddRelation(Relation, string, []IFieldDefinition)
	GetRelations() []string
	GetRelation(string) *Relation
	AddMulti(context.Context, *Data, AddOptions) (*Data, error)
	GetAll(context.Context, []string, GetAllOptions) (*Data, error)
	Edit(context.Context, IExpression, map[string]interface{}) error
	Delete(context.Context, IExpression) error
	FieldsToString([]string, map[string]interface{}) string
}

type IModelRow interface {
	GetValue(string) (interface{}, error)
}

type IExpression interface {
	GetProcessor(IExpressionProcessor) interface{}
}

type IExpressionProcessor interface {
	Eq(op1, op2 IExpression) interface{}
	Ne(op1, op2 IExpression) interface{}
	Lt(op1, op2 IExpression) interface{}
	Le(op1, op2 IExpression) interface{}
	Gt(op1, op2 IExpression) interface{}
	Ge(op1, op2 IExpression) interface{}
	In(op IExpression, arr []IExpression) interface{}
	And(operators []IExpression) interface{}
	Or(operators []IExpression) interface{}
	Any(localModel, extModel IModel, filter IExpression) interface{}
	ModelField(model IModel, fieldName string) interface{}
	Value(value interface{}) interface{}
}

type AddOptions struct {
	Replace bool
}

type GetAllOptions struct {
	Distinct    bool
	Filter      IExpression
	OrderBy     []Order
	Limit       uint64
	Offset      uint64
	RowsWoLimit *uint64
	ForUpdate   bool
}

type Order struct {
	FieldName string
	Desc      bool
}

type Relation struct {
	ExtModel                 IModel
	RelationType             RelationType
	LocalFieldsNames         []string
	FkFieldsNames            []string
	JunctionModel            IModel
	JunctionLocalFieldsNames []string
	JunctionFkFieldsNames    []string
	IsRequired               bool
	IsBack                   bool
}

type RelationType int

const (
	RELATION_ONE_TO_ONE   RelationType = iota
	RELATION_ONE_TO_MANY
	RELATION_MANY_TO_ONE
	RELATION_MANY_TO_MANY
)

func (r RelationType) String() string {
	switch r {
	case RELATION_ONE_TO_ONE:
		return "OneToOne"
	case RELATION_ONE_TO_MANY:
		return "OneToMany"
	case RELATION_MANY_TO_ONE:
		return "ManyToOne"
	case RELATION_MANY_TO_MANY:
		return "ManyToMany"
	default:
		return "Unknown"
	}
}
