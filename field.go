package model

import "reflect"

type IFieldDefinition interface {
	GetId() string
	GetCaption() string
	GetType() reflect.Type
	IsDerivable() bool
	IsRequired() bool
	GetDependsOn() []string
	Calc(map[string]interface{}) (interface{}, error)
	Check(interface{}) error
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
}

func (f *IntField) GetId() string                                    { return f.Id }
func (f *IntField) GetCaption() string                               { return f.Caption }
func (f *IntField) GetType() reflect.Type                            { return reflect.TypeOf(int(0)) }
func (f *IntField) IsDerivable() bool                                { return false }
func (f *IntField) IsRequired() bool                                 { return f.Required }
func (f *IntField) GetDependsOn() []string                           { return nil }
func (f *IntField) Calc(map[string]interface{}) (interface{}, error) { return nil, nil }
func (f *IntField) Check(v interface{}) error {
	if f.CheckFunc != nil {
		return f.CheckFunc(v)
	} else {
		return nil
	}
}
func (f *IntField) CloneForFK(id string, caption string, required bool) IFieldDefinition {
	return &IntField{id, caption, required, f.CheckFunc}
}

type StringField struct {
	Id        string
	Caption   string
	Required  bool
	CheckFunc func(interface{}) error
}

func (f *StringField) GetId() string                                    { return f.Id }
func (f *StringField) GetCaption() string                               { return f.Caption }
func (f *StringField) GetType() reflect.Type                            { return reflect.TypeOf(int(0)) }
func (f *StringField) IsRequired() bool                                 { return f.Required }
func (f *StringField) IsDerivable() bool                                { return false }
func (f *StringField) GetDependsOn() []string                           { return nil }
func (f *StringField) Calc(map[string]interface{}) (interface{}, error) { return nil, nil }
func (f *StringField) Check(v interface{}) error {
	if f.CheckFunc != nil {
		return f.CheckFunc(v)
	} else {
		return nil
	}
}
func (f *StringField) CloneForFK(id string, caption string, required bool) IFieldDefinition {
	return &StringField{id, caption, required, f.CheckFunc}
}

type DerivableField struct {
	Id        string
	Caption   string
	DependsOn []string
	Get       func(map[string]interface{}) (interface{}, error)
}

func (f *DerivableField) GetId() string                                        { return f.Id }
func (f *DerivableField) GetCaption() string                                   { return f.Caption }
func (f *DerivableField) GetType() reflect.Type                                { return reflect.TypeOf(int(0)) }
func (f *DerivableField) IsRequired() bool                                     { return false }
func (f *DerivableField) IsDerivable() bool                                    { return true }
func (f *DerivableField) GetDependsOn() []string                               { return f.DependsOn }
func (f *DerivableField) Calc(row map[string]interface{}) (interface{}, error) { return f.Get(row) }
func (f *DerivableField) Check(v interface{}) error                            { return nil }
func (f *DerivableField) CloneForFK(id string, caption string, required bool) IFieldDefinition {
	panic("Derivable field cannot be FK")
}
