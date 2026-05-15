package request

import (
	"errors"
	"net/http"
)

// ErrValidation and related vars are standard error values used in API response helpers.
var (
	ErrValidation            = errors.New("one or more validation errors occurred")
	ErrRecordNotFound        = errors.New("record not found")
	ErrCouldNotSaveRecord    = errors.New("could not save record")
	ErrCouldNotParseRequest  = errors.New("could not parse request")
	ErrProcessingRequest     = errors.New("there was an error processing the request")
	ErrBadUsernameOrPassword = errors.New("bad username or password")
	ErrNotAuthorised         = errors.New("not authorised")

	NotFoundMsg = ErrorMessage{"message": "not found"}
)

const (
	errResponseName   = "there was an error processing request"
	errResponseStatus = "error"
)

// ErrorMessage is a map of field names to error message strings for API error responses.
type ErrorMessage map[string]string

// GetSingleItemResp wraps a single data item in a standard response envelope.
func GetSingleItemResp(data interface{}) SingleItemResp {
	return SingleItemResp{
		Data: data,
	}
}

// GeneralErrResp builds a general error response with the provided message and HTTP status code.
func GeneralErrResp(msg string, statusCode int) GeneralErrorResp {
	return GeneralErrorResp{
		Name:    errResponseName,
		Message: msg,
		Code:    statusCode,
		Status:  errResponseStatus,
		Errors:  nil,
	}
}

// ServerErrResp builds a 500 internal server error response.
func ServerErrResp(msg string) GeneralErrorResp {
	return GeneralErrorResp{
		Name:    errResponseName,
		Message: msg,
		Code:    http.StatusInternalServerError,
		Status:  errResponseStatus,
		Errors:  nil,
	}
}

// GetNotAuthorisedResp builds a 403 forbidden response.
func GetNotAuthorisedResp() GeneralErrorResp {
	return GeneralErrorResp{
		Name:    errResponseName,
		Message: ErrNotAuthorised.Error(),
		Code:    http.StatusForbidden,
		Status:  errResponseStatus,
		Errors:  nil,
	}
}

// GetNotFoundResp for returning 404 messages.
func GetNotFoundResp() GeneralErrorResp {
	return GeneralErrorResp{
		Name:    ErrRecordNotFound.Error(),
		Message: ErrRecordNotFound.Error(),
		Code:    http.StatusNotFound,
		Status:  errResponseStatus,
		Errors:  nil,
	}
}

// GetListResp wraps a collection and its pagination metadata in a standard response envelope.
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
		Status:  errResponseStatus,
		Errors:  errors,
	}
}

// UnableToParseResp will return a message indicating that the JSON request could not be parsed.
func UnableToParseResp() GeneralErrorResp {
	return GeneralErrorResp{
		Name:    "Parsing error",
		Message: ErrCouldNotParseRequest.Error(),
		Code:    http.StatusBadRequest,
		Status:  errResponseStatus,
		Errors:  nil,
	}
}

// ErrorProcessingRequestResp will return a message indicating that there was an error processing request.
func ErrorProcessingRequestResp() GeneralErrorResp {
	return GeneralErrorResp{
		Name:    "Parsing error",
		Message: ErrProcessingRequest.Error(),
		Code:    http.StatusInternalServerError,
		Status:  errResponseStatus,
		Errors:  nil,
	}
}
