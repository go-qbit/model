package test

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-qbit/model"
	"github.com/go-qbit/qerror"
	"sort"
)

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

func NewStorage() *Storage {
	s := &Storage{
		data:   make(map[string][]DataRow),
		models: make(map[string]model.IModel),
	}

	return s
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

func (s *Storage) Add(ctx context.Context, m model.IModel, data []map[string]interface{}) ([]interface{}, error) {
	pKeys := make([]interface{}, len(data))

	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, row := range data {
		s.data[m.GetId()] = append(s.data[m.GetId()], row)

		pk := make(map[string]interface{})
		for _, pkName := range m.GetPKFieldsNames() {
			pk[pkName] = row[pkName]
		}
	}

	return pKeys, nil
}

func (s *Storage) Query(ctx context.Context, m model.IModel, fieldsNames []string, options model.QueryOptions) ([]map[string]interface{}, error) {
	var res []map[string]interface{}

	s.mtx.RLock()
	s.mtx.RUnlock()

	for _, row := range s.data[m.GetId()] {
		resRow := make(map[string]interface{})

		if options.Filter != nil {
			filterRes, err := options.Filter.Eval(row)
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

		for _, fieldName := range fieldsNames {
			resRow[fieldName] = row[fieldName]
		}

		res = append(res, resRow)
	}

	return res, nil
}
