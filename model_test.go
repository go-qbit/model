package model_test

import (
	"context"
	"testing"

	"github.com/go-qbit/model"
	"github.com/go-qbit/model/expr"
	"github.com/go-qbit/model/test"

	"github.com/stretchr/testify/suite"
)

type ModelTestSuite struct {
	suite.Suite
	storage *test.Storage
	user    *test.User
	phone   *test.Phone
	address *test.Address
	message *test.Message
}

func TestModelTestSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}

func (s *ModelTestSuite) SetupTest() {
	s.storage = test.NewStorage()

	s.user = test.NewUser(s.storage)
	s.phone = test.NewPhone(s.storage)
	s.message = test.NewMessage(s.storage)
	s.address = test.NewAddress(s.storage)

	model.AddOneToOneRelation(s.user, s.phone, false)
	model.AddOneToManyRelation(s.user, s.message)
	model.AddManyToManyRelation(s.user, s.address, s.storage)

	_, err := s.user.AddMulti(context.Background(), []map[string]interface{}{
		{"id": 1, "name": "Ivan", "lastname": "Sidorov"},
		{"id": 2, "name": "Petr", "lastname": "Ivanov"},
		{"id": 3, "name": "James", "lastname": "Bond"},
		{"id": 4, "name": "John", "lastname": "Connor"},
		{"id": 5, "name": "Sara", "lastname": "Connor"},
	})
	s.NoError(err)

	_, err = s.phone.AddMulti(context.Background(), []map[string]interface{}{
		{"id": 1, "country_code": 1, "code": 111, "number": 1111111},
		{"id": 3, "country_code": 3, "code": 333, "number": 3333333},
	})
	s.NoError(err)

	_, err = s.message.AddMulti(context.Background(), []map[string]interface{}{
		{"id": 10, "text": "Message 1", "fk__user__id": 1},
		{"id": 20, "text": "Message 2", "fk__user__id": 1},
		{"id": 30, "text": "Message 3", "fk__user__id": 1},
		{"id": 40, "text": "Message 4", "fk__user__id": 2},
	})
	s.NoError(err)

	_, err = s.address.AddMulti(context.Background(), []map[string]interface{}{
		{"id": 100, "country": "USA", "city": "Arlington", "address": "1022 Bridges Dr"},
		{"id": 200, "country": "USA", "city": "Fort Worth", "address": "7105 Plover Circle"},
		{"id": 300, "country": "USA", "city": "Crowley", "address": "524 Pecan Street"},
		{"id": 400, "country": "USA", "city": "Arlington", "address": "1022 Bridges Dr"},
		{"id": 500, "country": "USA", "city": "Louisville", "address": "1246 Everett Avenue"},
	})
	s.NoError(err)

	_, err = s.user.GetRelation("address").JunctionModel.AddMulti(context.Background(), []map[string]interface{}{
		{"fk__user__id": 1, "fk__address__id": 100},
		{"fk__user__id": 1, "fk__address__id": 200},
		{"fk__user__id": 2, "fk__address__id": 200},
		{"fk__user__id": 2, "fk__address__id": 300},
		{"fk__user__id": 3, "fk__address__id": 300},
		{"fk__user__id": 4, "fk__address__id": 400},
		{"fk__user__id": 5, "fk__address__id": 500},
	})
	s.NoError(err)
}

func (s *ModelTestSuite) TestModel_CheckRegisteredModels() {
	s.Equal(
		[]string{"_junction__user__address", "address", "message", "phone", "user"},
		s.storage.GetModelsNames(),
	)
}

func (s *ModelTestSuite) TestModel_GetFieldsNames() {
	s.Equal(
		[]string{"id", "name", "lastname", "fullname"},
		s.user.GetFieldsNames(),
	)

	s.Equal(
		[]string{"id", "text", "fk__user__id"},
		s.message.GetFieldsNames(),
	)
}

func (s *ModelTestSuite) TestModel_GetAll() {
	data, err := s.user.GetAll(
		context.Background(),
		[]string{
			// Self fields
			"id", "lastname", "fullname",
			// 1 to 1
			"phone.formated_number",
			// 1 to many
			"message.text",
			// many to many
			"address.city", "address.address",
		},
		model.GetAllOptions{
			Filter: expr.Lt(expr.ModelField("id"), expr.Value(4)),
		},
	)
	s.NoError(err)

	s.Equal([]map[string]interface{}{
		{
			"id":       1,
			"lastname": "Sidorov",
			"fullname": "Ivan Sidorov",
			"phone":    map[string]interface{}{"formated_number": "+1 (111) 1111111"},
			"message": []map[string]interface{}{
				{"text": "Message 1"},
				{"text": "Message 2"},
				{"text": "Message 3"},
			},
			"address": []map[string]interface{}{
				{"city": "Arlington", "address": "1022 Bridges Dr"},
				{"city": "Fort Worth", "address": "7105 Plover Circle"},
			},
		},
		{
			"id":       2,
			"lastname": "Ivanov",
			"fullname": "Petr Ivanov",
			"message": []map[string]interface{}{
				{"text": "Message 4"},
			},
			"address": []map[string]interface{}{
				{"city": "Fort Worth", "address": "7105 Plover Circle"},
				{"city": "Crowley", "address": "524 Pecan Street"},
			},
		},
		{
			"id":       3,
			"lastname": "Bond",
			"fullname": "James Bond",
			"phone":    map[string]interface{}{"formated_number": "+3 (333) 3333333"},
			"address": []map[string]interface{}{
				{"city": "Crowley", "address": "524 Pecan Street"},
			},
		},
	}, data)
}

func (s *ModelTestSuite) TestModel_GetAllToStruct() {
	type PhoneType struct {
		FormattedNumber string `field:"formated_number"`
	}

	type MessageType struct {
		Text string `field:"text"`
	}

	type AddressTYpe struct {
		City    string `field:"city"`
		Address string `field:"address"`
	}

	type UserType struct {
		Id        int           `field:"id"`
		Lastname  string        `field:"lastname"`
		Fullname  string        `field:"fullname"`
		Phone     PhoneType     `field:"phone"`
		PhonePtr  *PhoneType    `field:"phone"`
		Messages  []MessageType `field:"message"`
		Addresses []AddressTYpe `field:"address"`
	}

	var res []UserType

	s.NoError(s.user.GetAllToStruct(
		context.Background(),
		&res,
		model.GetAllOptions{
			Filter: expr.Lt(expr.ModelField("id"), expr.Value(4)),
		},
	))

	s.Equal(
		[]UserType{
			{
				Id:       1,
				Lastname: "Sidorov",
				Fullname: "Ivan Sidorov",
				Phone: PhoneType{
					FormattedNumber: "+1 (111) 1111111",
				},
				PhonePtr: &PhoneType{
					FormattedNumber: "+1 (111) 1111111",
				},
				Messages: []MessageType{
					{Text: "Message 1"},
					{Text: "Message 2"},
					{Text: "Message 3"},
				},
				Addresses: []AddressTYpe{
					{City: "Arlington", Address: "1022 Bridges Dr"},
					{City: "Fort Worth", Address: "7105 Plover Circle"},
				},
			},
			{
				Id:       2,
				Lastname: "Ivanov",
				Fullname: "Petr Ivanov",
				Messages: []MessageType{
					{Text: "Message 4"},
				},
				Addresses: []AddressTYpe{
					{City: "Fort Worth", Address: "7105 Plover Circle"},
					{City: "Crowley", Address: "524 Pecan Street"},
				},
			},
			{
				Id:       3,
				Lastname: "Bond",
				Fullname: "James Bond",
				Phone: PhoneType{
					FormattedNumber: "+3 (333) 3333333",
				},
				PhonePtr: &PhoneType{
					FormattedNumber: "+3 (333) 3333333",
				},
				Addresses: []AddressTYpe{
					{City: "Crowley", Address: "524 Pecan Street"},
				},
			},
		},
		res,
	)
}
