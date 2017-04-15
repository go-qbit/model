package model

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/go-qbit/qerror"
	"github.com/go-qbit/timelog"
)

// In
type exprInS struct {
	op     IExpression
	values []IExpression
}

func exprIn(op IExpression) *exprInS     { return &exprInS{op, nil} }
func (e *exprInS) Add(value IExpression) { e.values = append(e.values, value) }
func (e *exprInS) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.In(e.op, e.values)
}

// Model field
type exprModelFieldS struct {
	m     IModel
	field string
}

//func exprModelField(m IModel, field string) *exprModelFieldS {
//	return &exprModelFieldS{m, field}
//}
func (e *exprModelFieldS) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.ModelField(e.m, e.field)
}

// Value
type exprValueS struct {
	data interface{}
}

func exprValue(v interface{}) *exprValueS { return &exprValueS{v} }
func (e *exprValueS) GetProcessor(processor IExpressionProcessor) interface{} {
	return processor.Value(e.data)
}

type ModelLink struct {
	Pk  []interface{}
	Fks [][]interface{}
}

type BaseModel struct {
	id            string
	fields        []IFieldDefinition
	fieldsMtx     sync.RWMutex
	nameToField   map[string]IFieldDefinition
	pkFieldsNames []string
	extModels     map[string]Relation
	extModelsMtx  sync.RWMutex
	storage       IStorage
}

func NewBaseModel(id string, fields []IFieldDefinition, pkFieldsNames []string, storage IStorage) *BaseModel {
	m := &BaseModel{
		id:            id,
		fields:        fields,
		nameToField:   make(map[string]IFieldDefinition),
		pkFieldsNames: pkFieldsNames,
		extModels:     make(map[string]Relation),
		storage:       storage,
	}

	if err := storage.RegisterModel(m); err != nil {
		panic(err)
	}

	for _, field := range fields {
		m.nameToField[field.GetId()] = field
	}

	return m
}

func (m *BaseModel) GetId() string {
	return m.id
}

func (m *BaseModel) GetFieldsNames() []string {
	m.fieldsMtx.RLock()
	defer m.fieldsMtx.RUnlock()

	res := make([]string, len(m.fields))

	for i, field := range m.fields {
		res[i] = field.GetId()
	}

	return res
}

func (m *BaseModel) GetPKFieldsNames() []string {
	return m.pkFieldsNames
}

func (m *BaseModel) GetFieldDefinition(name string) IFieldDefinition {
	m.fieldsMtx.RLock()
	defer m.fieldsMtx.RUnlock()

	return m.nameToField[name]
}

func (m *BaseModel) AddField(field IFieldDefinition) {
	m.fieldsMtx.Lock()
	defer m.fieldsMtx.Unlock()

	if _, exists := m.nameToField[field.GetId()]; exists {
		panic(fmt.Sprintf("The model '%s' already has a field '%s'", m.GetId(), field.GetId()))
	}

	m.nameToField[field.GetId()] = field
	m.fields = append(m.fields, field)
}

func (m *BaseModel) GetAllFieldDependencies(fieldName string) ([]string, error) {
	resMap := make(map[string]struct{})
	if fieldDefinition, exists := m.nameToField[fieldName]; exists {
		for _, parentFieldName := range fieldDefinition.GetDependsOn() {
			resMap[parentFieldName] = struct{}{}
			parentDependencies, err := m.GetAllFieldDependencies(parentFieldName)
			if err != nil {
				return nil, err
			}
			for _, parent2FieldName := range parentDependencies {
				resMap[parent2FieldName] = struct{}{}
			}
		}
	} else {
		return nil, qerror.New("Unknown field '%s' in modeld `%s`", fieldName, m.id)
	}

	res := make([]string, 0, len(resMap))
	for name, _ := range resMap {
		res = append(res, name)
	}

	return res, nil
}

func (m *BaseModel) AddRelation(relation Relation, newFields []IFieldDefinition) {
	for _, newField := range newFields {
		m.AddField(newField)
	}

	m.extModelsMtx.Lock()
	defer m.extModelsMtx.Unlock()

	m.extModels[relation.ExtModel.GetId()] = relation
}

func (m *BaseModel) GetRelations() []string {
	m.extModelsMtx.RLock()
	defer m.extModelsMtx.RUnlock()

	res := make([]string, 0, len(m.extModels))
	for name, _ := range m.extModels {
		res = append(res, name)
	}

	sort.Strings(res)

	return res
}

func (m *BaseModel) GetRelation(modelName string) *Relation {
	m.extModelsMtx.RLock()
	defer m.extModelsMtx.RUnlock()

	relation, exists := m.extModels[modelName]
	if exists {
		return &relation
	} else {
		return nil
	}
}

func (m *BaseModel) AddMulti(ctx context.Context, data *Data, opts AddOptions) (*Data, error) {
	ctx = timelog.Start(ctx, m.GetId()+": AddMulti")
	defer timelog.Finish(ctx)

	if data.Len() == 0 {
		return nil, nil
	}

	fieldsMap := make(map[string]struct{})
	fields := make([]IFieldDefinition, len(data.Fields()))
	for i, fieldName := range data.Fields() {
		fields[i] = m.GetFieldDefinition(fieldName)
		if fields[i] == nil {
			return nil, qerror.New("Unknown field '%s' in model '%s'", fieldName, m.id)
		}

		if fields[i].IsDerivable() {
			return nil, qerror.New("The field '%s' is derivable in model %s, it cannot be added", fieldName, m.id)
		}

		fieldsMap[fieldName] = struct{}{}
	}

	for _, fieldName := range m.GetFieldsNames() {
		field := m.GetFieldDefinition(fieldName)
		if field.IsRequired() {
			if _, exists := fieldsMap[fieldName]; !exists {
				return nil, qerror.New("Missed required field '%s' in model '%s'", fieldName, m.id)
			}
		}
	}

	for _, row := range data.Data() {
		if len(row) != len(fields) {
			return nil, qerror.New("Invalid columns number")
		}

		for i, field := range fields {
			if field.IsRequired() && row[i] == nil {
				return nil, qerror.New("Missed required field '%s' value in model '%s'", field.GetId(), m.id)
			}

			var err error
			if row[i], err = field.Clean(row[i]); err != nil {
				return nil, err
			}

			if err := field.Check(row[i]); err != nil {
				return nil, err
			}
		}
	}

	return m.storage.Add(ctx, m, data, opts)
}

func (m *BaseModel) Link(ctx context.Context, extModel IModel, links []ModelLink) error {
	ctx = timelog.Start(ctx, m.GetId()+": Link to "+extModel.GetId())
	defer timelog.Finish(ctx)

	relation := m.GetRelation(extModel.GetId())
	if relation == nil {
		return qerror.New("No relation found between %s and %s", m.GetId(), extModel.GetId())
	}

	switch relation.RelationType {
	case RELATION_MANY_TO_MANY:
		data := NewEmptyData(append(relation.JunctionLocalFieldsNames, relation.JunctionFkFieldsNames...))
		for _, link := range links {
			for _, fk := range link.Fks {
				if err := data.Add(append(link.Pk, fk...)); err != nil {
					return err
				}
			}
		}

		_, err := relation.JunctionModel.AddMulti(ctx, data, AddOptions{Replace: true})

		return err

	default:
		panic("Not implemented") //ToDo
	}

	return nil
}

func (m *BaseModel) AddFromStructs(ctx context.Context, data interface{}, opts AddOptions) (*Data, error) {
	rt := reflect.TypeOf(data)

	if rt.Kind() != reflect.Slice || rt.Elem().Kind() != reflect.Struct {
		return nil, qerror.New("Invalid type '%s', must to slice of struct", rt.String())
	}

	fieldsNames := make([]string, 0, rt.Elem().NumField())
	fieldsNums := make([]int, 0, rt.Elem().NumField())
	for i := 0; i < rt.Elem().NumField(); i++ {
		if rt.Elem().Field(i).Tag.Get("field") == "-" {
			continue
		}

		fieldsNames = append(fieldsNames, m.structFieldToFieldName(rt.Elem().Field(i)))
		fieldsNums = append(fieldsNums, i)
	}

	rData := reflect.ValueOf(data)
	flatData := NewEmptyData(fieldsNames)

	for i := 0; i < rData.Len(); i++ {
		flatRow := make([]interface{}, len(fieldsNames))
		for j := 0; j < len(fieldsNames); j++ {
			flatRow[j] = rData.Index(i).Field(fieldsNums[j]).Interface()
		}
		flatData.Add(flatRow)
	}

	return m.AddMulti(ctx, flatData, opts)
}

func (m *BaseModel) GetAll(ctx context.Context, fieldsNames []string, opts GetAllOptions) (*Data, error) {
	ctx = timelog.Start(ctx, m.GetId()+": GetAll")
	defer timelog.Finish(ctx)

	requestedLocalFields := make(map[string]struct{})
	requestedExtFields := make(map[string]map[string]struct{})

	needLocalFields := make(map[string]struct{})
	needDerivableFieldsNames := make(map[string]struct{})
	extFields := make(map[string][]string)

	for _, fieldName := range fieldsNames {
		splittedFieldName := strings.SplitN(fieldName, ".", 2)

		// Local fields
		if len(splittedFieldName) == 1 {
			requestedLocalFields[fieldName] = struct{}{}

			field := m.GetFieldDefinition(fieldName)
			if field == nil {
				return nil, qerror.New("Unknown field '%s' in model '%s'", fieldName, m.id)
			}

			if field.IsDerivable() {
				needDerivableFieldsNames[fieldName] = struct{}{}
				dependencies, err := m.GetAllFieldDependencies(fieldName)
				if err != nil {
					return nil, err
				}
				for _, name := range dependencies {
					needLocalFields[name] = struct{}{}
				}
			}
			needLocalFields[fieldName] = struct{}{}
		} else {
			// External fields
			extFields[splittedFieldName[0]] = append(extFields[splittedFieldName[0]], splittedFieldName[1])
			if _, exists := requestedExtFields[splittedFieldName[0]]; !exists {
				requestedExtFields[splittedFieldName[0]] = make(map[string]struct{})
			}
			requestedExtFields[splittedFieldName[0]][strings.SplitN(splittedFieldName[1], ".", 2)[0]] = struct{}{}
		}
	}

	for extModelName, _ := range extFields {
		relation, exists := m.extModels[extModelName]
		if !exists {
			return nil, qerror.New("There is no relation between '%s' and '%s'", m.GetId(), extModelName)
		}

		for _, pkFieldName := range relation.LocalFieldsNames {
			needLocalFields[pkFieldName] = struct{}{}
		}
		for _, fkFieldName := range relation.FkFieldsNames {
			extFields[extModelName] = append(extFields[extModelName], fkFieldName)
		}
	}

	needLocalFieldsNamesArr := make([]string, 0, len(needLocalFields))
	for fieldName, _ := range needLocalFields {
		if !m.nameToField[fieldName].IsDerivable() {
			needLocalFieldsNamesArr = append(needLocalFieldsNamesArr, fieldName)
		}
	}
	valuesData, err := m.storage.Query(ctx, m, needLocalFieldsNamesArr, opts)
	if err != nil {
		return nil, err
	}

	resFields := make([]string, 0, len(requestedLocalFields)+len(requestedExtFields))
	for fieldName, _ := range requestedLocalFields {
		resFields = append(resFields, fieldName)
	}
	for fieldName, _ := range requestedExtFields {
		resFields = append(resFields, fieldName)
	}

	if valuesData.Len() == 0 {
		return NewEmptyData(resFields), nil
	}

	values := valuesData.Maps()

	extData := make(map[string]map[string][]map[string]interface{})
	localPk := m.GetPKFieldsNames()
	for extModelName, extFields := range extFields {
		relation := m.extModels[extModelName]
		extModel := relation.ExtModel

		if len(localPk) != 1 { //ToDo: implement PKs with 2 and more keys
			panic("Not implemented")
		}

		extValuesMap := make(map[string][]map[string]interface{})

		if relation.JunctionModel != nil {
			filter := exprIn(relation.JunctionModel.FieldExpr(relation.JunctionLocalFieldsNames[0]))
			for _, row := range values {
				filter.Add(exprValue(row[relation.LocalFieldsNames[0]]))
			}

			junctionFields := append(relation.JunctionLocalFieldsNames, relation.JunctionFkFieldsNames...)
			junctionValues, err := relation.JunctionModel.GetAll(ctx, junctionFields, GetAllOptions{Filter: filter})
			if err != nil {
				return nil, err
			}

			if junctionValues.Len() > 0 {
				junctionModel := relation.JunctionModel
				junctionValuesMap := make(map[string][]string)
				filter = exprIn(extModel.FieldExpr(relation.FkFieldsNames[0]))

				for _, value := range junctionValues.Maps() {
					key := junctionModel.FieldsToString(relation.JunctionFkFieldsNames, value)
					junctionValuesMap[key] = append(junctionValuesMap[key], junctionModel.FieldsToString(relation.JunctionLocalFieldsNames, value))
					filter.Add(exprValue(value[relation.JunctionFkFieldsNames[0]]))
				}

				orderBy := make([]Order, len(relation.FkFieldsNames))
				for i := range relation.FkFieldsNames {
					orderBy[i].FieldName = relation.FkFieldsNames[i]
				}
				extValues, err := extModel.GetAll(ctx, extFields, GetAllOptions{Filter: filter, OrderBy: orderBy})
				if err != nil {
					return nil, err
				}

				for _, extRow := range extValues.Maps() {
					for _, fk := range junctionValuesMap[extModel.FieldsToString(relation.FkFieldsNames, extRow)] {
						extValuesMap[fk] = append(extValuesMap[fk], extRow)
					}
				}
			}
		} else {
			filter := exprIn(extModel.FieldExpr(relation.FkFieldsNames[0]))
			for _, row := range values {
				filter.Add(exprValue(row[relation.LocalFieldsNames[0]]))
			}

			extValues, err := extModel.GetAll(ctx, extFields, GetAllOptions{Filter: filter})
			if err != nil {
				return nil, err
			}

			for _, extRow := range extValues.Maps() {
				stringFk := extModel.FieldsToString(relation.FkFieldsNames, extRow)
				extValuesMap[stringFk] = append(extValuesMap[stringFk], extRow)
			}
		}

		extData[extModelName] = extValuesMap
	}

	for _, value := range values {
		for extModelName, extModelData := range extData {
			relation := m.extModels[extModelName]
			extModel := relation.ExtModel
			extValues, exists := extModelData[extModel.FieldsToString(relation.LocalFieldsNames, value)]
			if !exists {
				continue
			}

			// Delete unrequested external fields
			for i, _ := range extValues {
				for fieldName, _ := range extValues[i] {
					if _, exists := requestedExtFields[extModelName][fieldName]; !exists {
						delete(extValues[i], fieldName)
					}
				}
			}

			switch relation.RelationType {
			case RELATION_ONE_TO_ONE, RELATION_MANY_TO_ONE:
				value[extModelName] = extValues[0]
			case RELATION_ONE_TO_MANY, RELATION_MANY_TO_MANY:
				value[extModelName] = extValues
			}
		}

		for fieldName, _ := range needDerivableFieldsNames {
			var err error
			value[fieldName], err = m.nameToField[fieldName].Calc(value)
			if err != nil {
				return nil, err
			}
		}

		for key, _ := range value {
			_, isLocalField := requestedLocalFields[key]
			_, isExtField := extFields[key]

			if !(isLocalField || isExtField) {
				delete(value, key)
			}
		}
	}

	res := NewEmptyData(resFields)
	for _, row := range values {
		resRow := make([]interface{}, len(resFields))

		for i, fieldName := range resFields {
			resRow[i] = row[fieldName]
		}

		if err := res.Add(resRow); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (m *BaseModel) GetAllToStruct(ctx context.Context, arr interface{}, options GetAllOptions) error {
	ctx = timelog.Start(ctx, m.GetId()+": GetAllToStruct")
	defer timelog.Finish(ctx)

	rt := reflect.TypeOf(arr)
	if rt.Kind() != reflect.Ptr {
		return qerror.New("Invalid type '%s', must be pointer to slice of struct", rt.String())
	}

	rt = rt.Elem()
	if rt.Kind() != reflect.Slice || rt.Elem().Kind() != reflect.Struct {
		return qerror.New("Invalid type '%s', must be pointer to slice of struct", rt.String())
	}

	fields, err := m.getFieldsFromStruct(rt.Elem())
	if err != nil {
		return err
	}

	data, err := m.GetAll(ctx, fields, options)
	if err != nil {
		return err
	}

	arrValue := reflect.MakeSlice(rt, data.Len(), data.Len())

	for i, row := range data.Maps() {
		if err := m.mapToVar(row, arrValue.Index(i)); err != nil {
			return err
		}
	}

	reflect.ValueOf(arr).Elem().Set(arrValue)

	return nil
}

func (m *BaseModel) Edit(ctx context.Context, filter IExpression, newValues map[string]interface{}) error {
	ctx = timelog.Start(ctx, m.GetId()+": Edit")
	defer timelog.Finish(ctx)

	for name := range newValues {
		if m.GetFieldDefinition(name) == nil {
			return qerror.New("Unknown field '%s' in model '%s'", name, m.id)
		}
	}

	return m.storage.Edit(ctx, m, filter, newValues)
}

func (m *BaseModel) Delete(ctx context.Context, filter IExpression) error {
	ctx = timelog.Start(ctx, m.GetId()+": Delete")
	defer timelog.Finish(ctx)

	return m.storage.Delete(ctx, m, filter)
}

func (m *BaseModel) FieldsToString(fieldsNames []string, row map[string]interface{}) string {
	buf := &bytes.Buffer{}

	for i, fieldName := range fieldsNames {
		if i > 0 {
			buf.WriteRune('|')
		}

		if row[fieldName] == nil {
			buf.WriteString("<nil>")
			continue
		}

		switch v := row[fieldName].(type) {
		case int:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int8:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int16:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int32:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int64:
			buf.WriteString(strconv.FormatInt(v, 10))
		case uint:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint8:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint16:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint32:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint64:
			buf.WriteString(strconv.FormatUint(v, 10))
		case string:
			buf.WriteString(v)

		case *int:
			buf.WriteString(strconv.FormatInt(int64(*v), 10))
		case *int8:
			buf.WriteString(strconv.FormatInt(int64(*v), 10))
		case *int16:
			buf.WriteString(strconv.FormatInt(int64(*v), 10))
		case *int32:
			buf.WriteString(strconv.FormatInt(int64(*v), 10))
		case *int64:
			buf.WriteString(strconv.FormatInt(*v, 10))
		case *uint:
			buf.WriteString(strconv.FormatUint(uint64(*v), 10))
		case *uint8:
			buf.WriteString(strconv.FormatUint(uint64(*v), 10))
		case *uint16:
			buf.WriteString(strconv.FormatUint(uint64(*v), 10))
		case *uint32:
			buf.WriteString(strconv.FormatUint(uint64(*v), 10))
		case *uint64:
			buf.WriteString(strconv.FormatUint(*v, 10))
		case *string:
			buf.WriteString(*v)

		default:
			panic(fmt.Sprintf("PkToString is not implemented for type %T", row[fieldName]))
		}
	}

	return buf.String()
}

func (m *BaseModel) FieldExpr(name string) *exprModelFieldS {
	if m.GetFieldDefinition(name) == nil {
		return nil
	}
	return &exprModelFieldS{m, name}
}

func (m *BaseModel) getFieldsFromStruct(t reflect.Type) ([]string, error) {
	if t.Kind() != reflect.Struct {
		return nil, qerror.New("Invalid type `%s` for getting fields", t.String())
	}

	var res []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		fieldName := m.structFieldToFieldName(field)

		switch field.Type.Kind() {
		case reflect.Struct:
			extFields, err := m.getFieldsFromStruct(field.Type)
			if err != nil {
				return nil, err
			}
			for _, extFieldName := range extFields {
				res = append(res, fieldName+"."+extFieldName)
			}
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Struct {
				extFields, err := m.getFieldsFromStruct(field.Type.Elem())
				if err != nil {
					return nil, err
				}
				for _, extFieldName := range extFields {
					res = append(res, fieldName+"."+extFieldName)
				}
			} else {
				res = append(res, fieldName)
			}
		case reflect.Ptr:
			if field.Type.Elem().Kind() == reflect.Struct {
				extFields, err := m.getFieldsFromStruct(field.Type.Elem())
				if err != nil {
					return nil, err
				}
				for _, extFieldName := range extFields {
					res = append(res, fieldName+"."+extFieldName)
				}
			} else {
				res = append(res, fieldName)
			}
		default:
			res = append(res, fieldName)
		}
	}

	return res, nil
}

func (m *BaseModel) mapToVar(v interface{}, s reflect.Value) error {
	switch s.Kind() {
	case reflect.Ptr:
		newVar := reflect.New(s.Type().Elem())
		m.mapToVar(v, newVar.Elem())
		s.Set(newVar)
	case reflect.Struct:
		vMap, ok := v.(map[string]interface{})
		if !ok {
			return qerror.New("Invalid type %T for converting to structure", v)
		}

		for i := 0; i < s.NumField(); i++ {
			field := s.Field(i)
			if fieldMap, exists := vMap[m.structFieldToFieldName(s.Type().Field(i))]; exists {
				if err := m.mapToVar(fieldMap, field); err != nil {
					return err
				}
			}
		}
	case reflect.Slice:
		switch v := v.(type) {
		case []byte:
			s.Set(reflect.ValueOf(v))
		case []map[string]interface{}:
			newSlice := reflect.MakeSlice(s.Type(), len(v), len(v))
			for i, _ := range v {
				if err := m.mapToVar(v[i], newSlice.Index(i)); err != nil {
					return err
				}
			}
			s.Set(newSlice)
		default:
			return qerror.New("Invalid type %T for converting to slice", v)
		}
	default:
		if s.Kind() != reflect.Ptr && reflect.TypeOf(v).Kind() == reflect.Ptr {
			rv := reflect.ValueOf(v)
			if !rv.IsNil() {
				s.Set(rv.Elem())
			}
		} else {
			s.Set(reflect.ValueOf(v))
		}
	}

	return nil
}

func (m *BaseModel) structFieldToFieldName(field reflect.StructField) string {
	fieldName, ok := field.Tag.Lookup("field")
	if !ok {
		fieldName = field.Name[0:1] + regexp.MustCompile("[A-Z]").ReplaceAllStringFunc(field.Name[1:], func(v string) string {
			return "_" + v
		})
		fieldName = strings.ToLower(fieldName)
	}

	return fieldName
}
