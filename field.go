package model

import (
	"reflect"
	"context"
)

type IFieldDefinition interface {
	GetId() string
	GetCaption() string
	GetType() reflect.Type
	GetStorageType() string
	IsDerivable() bool
	IsRequired() bool
	GetDependsOn() []string
	Calc(context.Context, map[string]interface{}) (interface{}, error)
	Check(context.Context, interface{}) error
	Clean(context.Context, interface{}) (interface{}, error)
	CloneForFK(id string, caption string, required bool) IFieldDefinition
}

type IFieldsDefinitions interface {
	GetFieldsDefinitions() []IFieldDefinition
	GetPKFieldsNames() []string
	AddField(IFieldDefinition)
}

type IntField struct {
	Id        string
	Caption   string
	Required  bool
	CheckFunc func(interface{}) error
	CleanFunc func(interface{}) (interface{}, error)
}

func (f *IntField) GetId() string                                                     { return f.Id }
func (f *IntField) GetCaption() string                                                { return f.Caption }
func (f *IntField) GetType() reflect.Type                                             { return reflect.TypeOf(int(0)) }
func (f *IntField) GetStorageType() string                                            { return "int" }
func (f *IntField) IsDerivable() bool                                                 { return false }
func (f *IntField) IsRequired() bool                                                  { return f.Required }
func (f *IntField) GetDependsOn() []string                                            { return nil }
func (f *IntField) Calc(context.Context, map[string]interface{}) (interface{}, error) { return nil, nil }
func (f *IntField) Check(ctx context.Context, v interface{}) error {
	if f.CheckFunc != nil {
		return f.CheckFunc(v)
	} else {
		return nil
	}
}
func (f *IntField) Clean(ctx context.Context, v interface{}) (interface{}, error) {
	if f.CleanFunc != nil {
		return f.CleanFunc(v)
	} else {
		return v, nil
	}
}
func (f *IntField) CloneForFK(id string, caption string, required bool) IFieldDefinition {
	return &IntField{id, caption, required, f.CheckFunc, f.CleanFunc}
}

type StringField struct {
	Id        string
	Caption   string
	Required  bool
	CheckFunc func(interface{}) error
	CleanFunc func(interface{}) (interface{}, error)
}

func (f *StringField) GetId() string                                                     { return f.Id }
func (f *StringField) GetCaption() string                                                { return f.Caption }
func (f *StringField) GetType() reflect.Type                                             { return reflect.TypeOf(int(0)) }
func (f *StringField) GetStorageType() string                                            { return "string" }
func (f *StringField) IsRequired() bool                                                  { return f.Required }
func (f *StringField) IsDerivable() bool                                                 { return false }
func (f *StringField) GetDependsOn() []string                                            { return nil }
func (f *StringField) Calc(context.Context, map[string]interface{}) (interface{}, error) { return nil, nil }
func (f *StringField) Check(ctx context.Context, v interface{}) error {
	if f.CheckFunc != nil {
		return f.CheckFunc(v)
	} else {
		return nil
	}
}
func (f *StringField) Clean(ctx context.Context, v interface{}) (interface{}, error) {
	if f.CleanFunc != nil {
		return f.CleanFunc(v)
	} else {
		return v, nil
	}
}
func (f *StringField) CloneForFK(id string, caption string, required bool) IFieldDefinition {
	return &StringField{id, caption, required, f.CheckFunc, f.CleanFunc}
}

type DerivableField struct {
	Id        string
	Caption   string
	DependsOn []string
	Get       func(context.Context, map[string]interface{}) (interface{}, error)
}

func (f *DerivableField) GetId() string                                                             { return f.Id }
func (f *DerivableField) GetCaption() string                                                        { return f.Caption }
func (f *DerivableField) GetType() reflect.Type                                                     { var v interface{}; return reflect.TypeOf(v) }
func (f *DerivableField) GetStorageType() string                                                    { return "interface{}" }
func (f *DerivableField) IsRequired() bool                                                          { return false }
func (f *DerivableField) IsDerivable() bool                                                         { return true }
func (f *DerivableField) GetDependsOn() []string                                                    { return f.DependsOn }
func (f *DerivableField) Calc(ctx context.Context, row map[string]interface{}) (interface{}, error) { return f.Get(ctx, row) }
func (f *DerivableField) Check(ctx context.Context, v interface{}) error                            { return nil }
func (f *DerivableField) Clean(ctx context.Context, v interface{}) (interface{}, error)             { return v, nil }
func (f *DerivableField) CloneForFK(id string, caption string, required bool) IFieldDefinition {
	panic("Derivable field cannot be FK")
}
