package test

import (
	"github.com/go-qbit/model"
	"context"
	"strconv"
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

				&model.DerivableField{
					Id:        "stringid",
					DependsOn: []string{"id"},
					Caption:   "Map data",
					Get: func(ctx context.Context, row map[string]interface{}) (interface{}, error) {
						return model.GetDerivableFieldsData(ctx, "mapdata").(map[int]string)[row["id"].(int)], nil
					},
				},
			},
			storage,
			model.BaseModelOpts{
				PkFieldsNames: []string{"id"},
				PrepareDerivableFieldsCtx: func(ctx context.Context, m model.IModel, requestedFields map[string]struct{}, rows []map[string]interface{}) {
					if _, exists := requestedFields["stringid"]; exists {
						mapData := make(map[int]string)
						for _, row := range rows {
							mapData[row["id"].(int)] = strconv.FormatInt(int64(row["id"].(int)), 10)
						}
						model.SetDerivableFieldsData(ctx, "mapdata", mapData)
					}
				},
			},
		),
	}
}
