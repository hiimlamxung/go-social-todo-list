package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"error_key"`
}

func NewErrorResponse(root error, msg string, log string, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}
	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootError().Error()
}

func NewFullErrorResponse(statusCode int, root error, msg string, log string, key string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func NewAuthorized(root error, msg string, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    root,
		Message:    msg,
		Key:        key,
	}
}

func NewCustomError(statusCode int, root error, msg string, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, msg, root.Error(), key)
	}
	return NewErrorResponse(errors.New(msg), msg, msg, key)
}

func ErrDB(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "Something went wrong with DB", err.Error(), "DB_ERROR")
}

func ErrInvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "Invalid request", err.Error(), "ErrInvalidRequest")
}

func ErrInternal(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "Something went wrong in the server", err.Error(), "ErrInternal")
}

func ErrCannotListEntity(entity string, err error) *AppError {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("Cannot list %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotList%s", entity),
	)
}

func ErrorCannotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("Cannot delete %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotDelete%s", entity),
	)
}

func ErrorCannotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("Cannot update %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotUpdate%s", entity),
	)
}

func ErrorCannotGetEntity(entity string, err error) *AppError {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("Cannot get %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotGet%s", entity),
	)
}

func ErrorCannotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("Cannot create %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotCreate%s", entity),
	)
}

func ErrorEntityDeleted(entity string, err error) *AppError {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("%s is deleted", strings.ToLower(entity)),
		fmt.Sprintf("Err%sDeleted", entity),
	)
}

func ErrorEntityExisted(entity string, err error) *AppError {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("%s is existed", strings.ToLower(entity)),
		fmt.Sprintf("Err%sExisted", entity),
	)
}

func ErrorEntityNotFound(entity string, err error) *AppError {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("%s is not found", strings.ToLower(entity)),
		fmt.Sprintf("Err%sNotFound", entity),
	)
}

func ErrorNoPermission(err error) *AppError {
	return NewCustomError(
		http.StatusForbidden,
		err,
		"You have no permission",
		"ErrNoPermission",
	)
}

var RecordNotFound = errors.New("Record not found !!!")
