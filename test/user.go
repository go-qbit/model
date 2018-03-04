package test

import (
	"github.com/go-qbit/model"
	"context"
)

type User struct {
	*model.BaseModel
}

func NewUser(storage model.IStorage) *User {
	return &User{
		BaseModel: model.NewBaseModel(
			"user",
			[]model.IFieldDefinition{
				&model.IntField{
					Id:      "id",
					Caption: "ID",
				},

				&model.StringField{
					Id:       "name",
					Caption:  "Name",
					Required: true,
				},

				&model.StringField{
					Id:       "lastname",
					Caption:  "Lastame",
					Required: true,
				},

				&model.DerivableField{
					Id:        "fullname",
					Caption:   "Full name",
					DependsOn: []string{"name", "lastname"},
					Get: func(ctx context.Context, row map[string]interface{}) (interface{}, error) {
						return row["name"].(string) + " " + row["lastname"].(string), nil
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
