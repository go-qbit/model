package model

import (
	"fmt"

	"github.com/go-qbit/qerror"
)

type AddError struct {
	*qerror.BaseError
	Message string
}

func AddErrorf(message string, a ...interface{}) *AddError {
	return &AddError{qerror.New(1), fmt.Sprintf(message, a...)}
}

func (e *AddError) Error() string {
	return e.Message + "\n" + e.BaseError.Error()
}

type EditError struct {
	*qerror.BaseError
	Message string
}

func EditErrorf(message string, a ...interface{}) *EditError {
	return &EditError{qerror.New(1), fmt.Sprintf(message, a...)}
}

func (e *EditError) Error() string {
	return e.Message + "\n" + e.BaseError.Error()
}

type DeleteError struct {
	*qerror.BaseError
	Message string
}

func DeleteErrorf(message string, a ...interface{}) *DeleteError {
	return &DeleteError{qerror.New(1), fmt.Sprintf(message, a...)}
}

func (e *DeleteError) Error() string {
	return e.Message + "\n" + e.BaseError.Error()
}

type FieldError struct {
	*qerror.BaseError
	Field   string
	Message string
}

func FieldErrorf(field, message string, a ...interface{}) *FieldError {
	return &FieldError{qerror.New(1), field, fmt.Sprintf(message, a...)}
}

func (e *FieldError) Error() string {
	return e.Message + "\n" + e.BaseError.Error()
}
