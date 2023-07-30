package goaerrors

import (
	goa "goa.design/goa/v3/pkg"
)

func NewInternalError() *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "internal_error",
		ID:      goa.NewErrorID(),
		Message: "unexpected error occurred",
		Fault:   true,
	}
}
