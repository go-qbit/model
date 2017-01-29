package model

import "github.com/go-qbit/qerror"

type Data struct {
	fields    []string
	fieldNums map[string]int
	data      [][]interface{}
}

func NewData(fields []string, data [][]interface{}) *Data {
	d := &Data{
		fields:    fields,
		fieldNums: make(map[string]int),
		data:      data,
	}

	for i, field := range fields {
		d.fieldNums[field] = i
	}

	return d
}

func NewEmptyData(fields []string) *Data {
	return NewData(fields, nil)
}

func (d *Data) FieldNum(field string) int {
	n, exists := d.fieldNums[field]

	if !exists {
		return -1
	}

	return n
}

func (d *Data) Len() int {
	return len(d.data)
}

func (d *Data) Add(row []interface{}) error {
	if len(d.fields) != len(row) {
		qerror.New("Invalid row size, must be %d", len(d.fields))
	}

	d.data = append(d.data, row)

	return nil
}

func (d *Data) Fields() []string {
	return d.fields
}

func (d *Data) Data() [][]interface{} {
	return d.data
}

func (d *Data) Maps() []map[string]interface{} {
	res := make([]map[string]interface{}, len(d.data))

	for i, row := range d.data {
		res[i] = make(map[string]interface{})
		for j, field := range d.fields {
			if row[j] != nil {
				res[i][field] = row[j]
			}
		}
	}

	return res
}

func (d *Data) GetFieldsData(fields []string) *Data {
	res := NewEmptyData(fields)

	fieldsNums := make([]int, len(fields))
	for i, field := range fields {
		fieldsNums[i] = d.FieldNum(field)
	}

	res.data = make([][]interface{}, len(d.data))

	for i, row := range d.data {
		res.data[i] = make([]interface{}, len(fields))
		for j, fieldN := range fieldsNums {
			if fieldN != -1 {
				res.data[i][j] = row[fieldN]
			}
		}
	}

	return res
}
