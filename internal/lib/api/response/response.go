package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Ok() Response {
	return Response{
		Status: "Successful",
	}
}

func Error(msg string) Response {
	return Response{
		Status: "Error",
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response { // принимаем список ошибкок валидации
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URl", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: "Error",
		Error:  strings.Join(errMsgs, ", "),
	}
}
