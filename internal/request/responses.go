package request

import (
	"errors"
	"net/http"
)

var (
	ErrValidation           = errors.New("one or more validation errors occurred")
	ErrRecordNotFound       = errors.New("record not found")
	ErrCouldNotSaveRecord   = errors.New("could not save record")
	ErrCouldNotParseRequest = errors.New("could not parse request")
	ErrProcessingRequest    = errors.New("there was an error processing the request")
	ErrNotAuthorised        = errors.New("not authorised")

	NotFoundMsg = ErrorMessage{"message": "not found"}
)

type ErrorMessage map[string]string

func GetSingleItemResp(data interface{}) SingleItemResp {
	return SingleItemResp{
		Data: data,
	}
}

func GeneralErrResp(msg string, statusCode int) GeneralErrorResp {
	return GeneralErrorResp{
		Name:    "there was an error processing request",
		Message: msg,
		Code:    statusCode,
		Status:  "error",
		Errors:  nil,
	}
}

func ServerErrResp(msg string) GeneralErrorResp {
	return GeneralErrorResp{
		Name:    "there was an error processing request",
		Message: msg,
		Code:    http.StatusInternalServerError,
		Status:  "error",
		Errors:  nil,
	}
}

func GetNotAuthorisedResp() GeneralErrorResp {
	return GeneralErrorResp{
		Name:    "there was an error processing request",
		Message: ErrNotAuthorised.Error(),
		Code:    http.StatusForbidden,
		Status:  "error",
		Errors:  nil,
	}
}

// GetNotFoundResp for returning 404 messages.
func GetNotFoundResp() GeneralErrorResp {
	return GeneralErrorResp{
		Name:    ErrRecordNotFound.Error(),
		Message: ErrRecordNotFound.Error(),
		Code:    http.StatusNotFound,
		Status:  "error",
		Errors:  nil,
	}
}

func GetListResp(data interface{}, pagination *Request) CollResp {
	return CollResp{
		Data:    data,
		Request: *pagination,
	}
}

// GetValidateErrResp will prepare the error response. It will default to a predefined error for Message but
// will override it if one is supplied.
func GetValidateErrResp(errors ErrMsgs, errs ...string) GeneralErrorResp {
	err := ErrValidation.Error()
	if len(errs) > 0 {
		err = errs[0]
	}

	return GeneralErrorResp{
		Name:    "Validation failed",
		Message: err,
		Code:    http.StatusBadRequest,
		Status:  "error",
		Errors:  errors,
	}
}

// UnableToParseResp will return a message indicating that the JSON request could not be parsed.
func UnableToParseResp() GeneralErrorResp {
	return GeneralErrorResp{
		Name:    "Parsing error",
		Message: ErrCouldNotParseRequest.Error(),
		Code:    http.StatusBadRequest,
		Status:  "error",
		Errors:  nil,
	}
}

// ErrorProcessingRequestResp will return a message indicating that there was an error processing request.
func ErrorProcessingRequestResp() GeneralErrorResp {
	return GeneralErrorResp{
		Name:    "Parsing error",
		Message: ErrProcessingRequest.Error(),
		Code:    http.StatusInternalServerError,
		Status:  "error",
		Errors:  nil,
	}
}
