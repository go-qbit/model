package test

import (
	"strings"

	"github.com/go-qbit/model"
)

type Message struct {
	*model.BaseModel
}

func NewMessage(storage model.IStorage) *Message {
	return &Message{
		model.NewBaseModel(
			"message",
			[]model.IFieldDefinition{
				&model.IntField{
					Id:      "id",
					Caption: "ID",
				},

				&model.StringField{
					Id:       "text",
					Caption:  "Message text",
					Required: true,
					CleanFunc: func(v interface{}) (interface{}, error) {
						return strings.Trim(v.(string), " \t"), nil
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
