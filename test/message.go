package test

import (
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
				},
			},
			[]string{"id"},
			storage,
		),
	}
}
