package response

import (
	"fmt"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    string = "OK"
	StatusError string = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("Field %s is required", err.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("Field %s is invalid url", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("Field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
