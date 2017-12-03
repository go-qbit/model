package test

import (
	"github.com/go-qbit/model"
)

type Address struct {
	*model.BaseModel
}

func NewAddress(storage model.IStorage) *Address {
	return &Address{
		model.NewBaseModel(
			"address",
			[]model.IFieldDefinition{
				&model.IntField{
					Id:      "id",
					Caption: "ID",
				},

				&model.StringField{
					Id:       "country",
					Caption:  "Country",
					Required: true,
				},

				&model.StringField{
					Id:      "city",
					Caption: "City",
				},

				&model.StringField{
					Id:      "address",
					Caption: "Address",
				},
			},
			storage,
			model.BaseModelOpts{
				PkFieldsNames: []string{"id"},
			},
		),
	}
}
