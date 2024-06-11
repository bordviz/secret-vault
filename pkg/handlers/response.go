package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	StatusError = "error"
	StatusOk    = "ok"
)

type Response struct {
	Status string      `json:"status"`
	Detail interface{} `json:"detail"`
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, detail interface{}) {
	render.Status(r, status)
	render.JSON(w, r, Response{Status: StatusError, Detail: detail})
}

func ValidationError(w http.ResponseWriter, r *http.Request, errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	render.Status(r, 422)
	render.JSON(w, r, Response{
		Status: StatusError,
		Detail: strings.Join(errMsgs, ", "),
	})

	return strings.Join(errMsgs, ", ")
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, status)
	render.JSON(w, r, data)
}
