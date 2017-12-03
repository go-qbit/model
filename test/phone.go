package test

import (
	"fmt"

	"github.com/go-qbit/model"
)

type Phone struct {
	*model.BaseModel
}

func NewPhone(storage model.IStorage) *Phone {
	return &Phone{
		model.NewBaseModel(
			"phone",
			[]model.IFieldDefinition{
				&model.IntField{
					Id:      "id",
					Caption: "ID",
				},

				&model.IntField{
					Id:       "country_code",
					Caption:  "Country code",
					Required: true,
				},

				&model.IntField{
					Id:       "code",
					Caption:  "Code",
					Required: true,
				},

				&model.IntField{
					Id:       "number",
					Caption:  "Number",
					Required: true,
				},

				&model.DerivableField{
					Id:        "formated_number",
					Caption:   "Formated number",
					DependsOn: []string{"country_code", "code", "number"},
					Get: func(row map[string]interface{}) (interface{}, error) {
						return fmt.Sprintf("+%d (%d) %d", row["country_code"], row["code"], row["number"]), nil
					},
				},
			},
			storage,
			model.BaseModelOpts{
				PkFieldsNames: []string{"id"},
			},
		),
	}
}
