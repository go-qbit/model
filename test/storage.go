package test

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/go-qbit/model"
	"github.com/go-qbit/qerror"
	"github.com/go-qbit/timelog"
)

var exprProcessor = &ExprProcessor{}

type Storage struct {
	data      map[string][]DataRow
	mtx       sync.RWMutex
	models    map[string]model.IModel
	modelsMtx sync.RWMutex
}

type DataRow map[string]interface{}

func (r DataRow) GetValue(name string) (interface{}, error) {
	return r[name], nil
}

func (r DataRow) SetValue(name string, value interface{}) error {
	r[name] = value
	return nil
}

func NewStorage() *Storage {
	s := &Storage{
		data:   make(map[string][]DataRow),
		models: make(map[string]model.IModel),
	}

	return s
}

func (s *Storage) NewModel(id string, fields []model.IFieldDefinition, pk []string) model.IModel {
	return model.NewBaseModel(id, fields, pk, s)
}

func (s *Storage) RegisterModel(m model.IModel) error {
	s.modelsMtx.Lock()
	defer s.modelsMtx.Unlock()

	if _, exists := s.models[m.GetId()]; exists {
		return qerror.New("Model '%s' is already exists", m.GetId())
	}

	s.models[m.GetId()] = m

	return nil
}

func (s *Storage) GetModelsNames() []string {
	s.modelsMtx.RLock()
	defer s.modelsMtx.RUnlock()

	res := make([]string, 0, len(s.models))
	for k, _ := range s.models {
		res = append(res, k)
	}

	sort.Strings(res)

	return res
}

func (s *Storage) Add(ctx context.Context, m model.IModel, data *model.Data, opts model.AddOptions) (*model.Data, error) {
	ctx = timelog.Start(ctx, "Storage.Add")
	defer timelog.Finish(ctx)

	pKeys := model.NewEmptyData(m.GetPKFieldsNames())

	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, row := range data.Data() {
		dataRow := make(map[string]interface{})
		for i, field := range data.Fields() {
			dataRow[field] = row[i]
		}

		s.data[m.GetId()] = append(s.data[m.GetId()], dataRow)

		pk := make([]interface{}, len(m.GetPKFieldsNames()))
		for i, pkName := range m.GetPKFieldsNames() {
			pk[i] = dataRow[pkName]
		}

		pKeys.Add(pk)
	}

	return pKeys, nil
}

func (s *Storage) Query(ctx context.Context, m model.IModel, fieldsNames []string, options model.GetAllOptions) (*model.Data, error) {
	ctx = timelog.Start(ctx, "Storage.Query")
	defer timelog.Finish(ctx)

	res := model.NewEmptyData(fieldsNames)

	s.mtx.RLock()
	defer s.mtx.RUnlock()

	for _, row := range s.data[m.GetId()] {
		resRow := make([]interface{}, len(fieldsNames))

		if options.Filter != nil {
			filterRes, err := options.Filter.GetProcessor(exprProcessor).(EvalFunc)(row)
			if err != nil {
				return nil, err
			}
			matched, ok := filterRes.(bool)
			if !ok {
				return nil, fmt.Errorf("Invalid filter, must have bool type, not %t", filterRes)
			}
			if !matched {
				continue
			}
		}

		for i, fieldName := range fieldsNames {
			resRow[i] = row[fieldName]
		}

		res.Add(resRow)
	}

	return res, nil
}

func (s *Storage) Edit(ctx context.Context, m model.IModel, filter model.IExpression, newValues map[string]interface{}) error {
	ctx = timelog.Start(ctx, "Storage.Edit")
	defer timelog.Finish(ctx)

	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, row := range s.data[m.GetId()] {
		filterRes, err := filter.GetProcessor(exprProcessor).(EvalFunc)(row)
		if err != nil {
			return err
		}
		matched, ok := filterRes.(bool)
		if !ok {
			return fmt.Errorf("Invalid filter, must have bool type, not %t", filterRes)
		}
		if !matched {
			continue
		}

		for name, value := range newValues {
			if err := row.SetValue(name, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, m model.IModel, filter model.IExpression) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	newData := make([]DataRow, 0)
	for _, row := range s.data[m.GetId()] {
		filterRes, err := filter.GetProcessor(exprProcessor).(EvalFunc)(row)
		if err != nil {
			return err
		}
		matched, ok := filterRes.(bool)
		if !ok {
			return fmt.Errorf("Invalid filter, must have bool type, not %t", filterRes)
		}
		if matched {
			continue
		}

		newData = append(newData, row)
	}

	s.data[m.GetId()] = newData

	return nil
}
