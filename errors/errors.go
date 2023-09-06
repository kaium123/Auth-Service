package errors

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Body map[string]interface{}

func ValidationErrors(err error) Body {
	return validationResponse(err)
}

func GenerateErrorResponseBody(err error) Body {
	message := err.Error()
	return readFromMap(message)
}

func readFromMap(message string) Body {
	return GenerateResponseBody(message)
}

func GenerateResponseBody(message string) Body {
	return Body{
		"message": message,
	}
}

func validationResponse(err error) Body {
	return Body{
		"validation_error": err,
	}
}

type ErrorIndex struct {
	Error error
	Index string
}

type Queue struct {
	item_value []ErrorIndex
}

func (q *Queue) Enqueue(err error, index string) {
	q.item_value = append(q.item_value, ErrorIndex{Error: err, Index: index}) //used to add items
}

func (q *Queue) Dequeue() ErrorIndex {
	item := q.item_value[0]
	q.item_value = q.item_value[1:] //used to remove items
	return item
}

func (q *Queue) IsEmpty() bool {
	return len(q.item_value) == 0
}

func ErrorTraverseForUnitTest(err error) map[string][]ErrorType {

	errorsMap := make(map[string][]ErrorType, 0)
	if err == nil {
		return errorsMap
	}

	q := Queue{}
	q.Enqueue(err, "estimate_errors")
	for {
		if q.IsEmpty() {
			break
		}
		cur := q.Dequeue()
		if cur.Error == nil {
			continue
		}

		v := ApplicationError{}
		if errors.As(cur.Error, &v) {
			if len(v.Errs) > 0 {
				for _, err := range v.Errs {
					q.Enqueue(err, cur.Index)
				}
				continue
			} else {
				tmpErrorKey := cur.Index
				errorField, _ := v.TranslationParams["field"].(string)
				if errorField != "" {
					tmpErrorKey += "_"
				}
				tmpErrorKey += errorField
				errorsMap[tmpErrorKey] = append(errorsMap[tmpErrorKey], (v.ErrorType))
			}

		} else {
			for idx, e := range cur.Error.(validation.Errors) {
				v := ApplicationError{}
				if errors.As(e, &v) {
					if len(v.Errs) > 0 {
						for _, err := range v.Errs {
							q.Enqueue(err, idx)
						}
						continue
					}
					tmpErrorKey := cur.Index
					errorField, _ := v.TranslationParams["field"].(string)
					if errorField != "" {
						tmpErrorKey += "_"
					}
					tmpErrorKey += errorField
					errorsMap[tmpErrorKey] = append(errorsMap[tmpErrorKey], (v.ErrorType))
				} else {
					q.Enqueue(e, cur.Index+"_"+idx)
				}
			}
		}

	}

	return errorsMap
}

// var StatusCode = make(map[error]int, 0)

// func SetStatusCode(err error, status int) error {
// 	StatusCode[err] = status
// 	fmt.Println(StatusCode[err])
// 	return err
// }

// func GetStatusCode(err error) int {
// 	fmt.Println(err)
// 	errVal := err
// 	fmt.Println(StatusCode[err])
// 	return StatusCode[errVal]
// }

// var v = make(map[*error]int, 0)

// func SetStatusCode(err *error, status int) *error {
// 	v[err] = status
// 	return err
// }

// func GetStatusCode(err *error) int {
// 	return v[err]
// }
// func SetPointer(err error) *error {
// 	return &err
// }

type ErrorType int

const (
	NotFoundErr ErrorType = iota + 1
	InvalidEmail
	UnKnownErr
	RangeValidationErr
	InvalidDate
	Overflow
	InvalidType
	InvalidOrder
	MisMatchSubtotal
	MisMatchTotal
	ZeroEstimateLineItems
	InvalidLineItemType
	RequiredField
	InvalidAmount
	LengthValidation
	Underflow
	InvaidPhone
	InvaidCountryCode
	InvaidMobile
	InvalidStatus
	InvalidOrderBy
	InvalidOrderType
	ExcelFieldError
	ExcelFieldMissing
	ExcelMissingData
)

type ApplicationError struct {
	ErrorType         ErrorType
	TranslationKey    string
	TranslationParams map[string]interface{}
	HttpCode          int
	Errs              []ApplicationError
}

func (e *ApplicationError) Join(errs ...ApplicationError) {
	e.Errs = append(e.Errs, errs...)
	//return e
}

func (e *ApplicationError) Unwrap() []ApplicationError {
	return e.Errs
}

func (e ApplicationError) Error() string {
	return e.TranslationKey
}

func (e ApplicationError) Code() string {
	return e.TranslationKey
}

func (e ApplicationError) Message() string {
	return e.TranslationKey
}

func (e ApplicationError) SetMessage(message string) validation.Error {
	e.TranslationKey = message
	return e
}

func (e ApplicationError) Params() map[string]interface{} {
	return e.TranslationParams
}

func (e ApplicationError) SetParams(params map[string]interface{}) validation.Error {
	e.TranslationParams = params
	return e
}

type GinError struct{}

func GetErrorMessage(err error) string {
	aError := ApplicationError{}
	if errors.As(err, &aError) {
		return aError.TranslationKey
	}
	return fmt.Sprint(err)
}

func (g GinError) GetErrorMessage(err error) map[string]interface{} {
	aError := ApplicationError{}
	if errors.As(err, &aError) {

		return GenerateErrorResponseBody(err)
	}
	return ValidationErrors(err)
}

func (g GinError) GetStatusCode(err error) int {
	aError := ApplicationError{}
	if errors.As(err, &aError) {
		return aError.HttpCode
	}
	return http.StatusBadRequest
}

func GetStatusCode(err error) int {
	aError := ApplicationError{}
	if errors.As(err, &aError) {
		return aError.HttpCode
	}
	return 500
}

func GetErrorType(err error) int {
	aError := ApplicationError{}
	if errors.As(err, &aError) {
		return int(aError.ErrorType)
	}
	return 1000
}

func (g GinError) ErrorTraverse(err error) Body {
	aError := ApplicationError{}
	if errors.As(err, &aError) {
		if len(aError.Errs) == 0 {
			return GenerateErrorResponseBody(err)
		}

		mp := make(map[string]interface{}, 0)
		for _, e := range aError.Errs {
			field, _ := e.TranslationParams["field"].(string)
			if len(e.Errs) > 0 {
				childErr := g.ErrorTraverse(e)
				for k, ee := range childErr {
					mp[k] = ee
				}
			} else {
				mp[field] = e.Error()
			}
		}

		field, _ := aError.TranslationParams["field"].(string)
		errBody := map[string]interface{}{field: mp}
		return errBody
	} else {
		return validationResponse(err)
	}
}

// var (
// 	// Common
// 	ErrEstimateNotFound = SetStatusCode(SetPointer(errors.New("estimateNotFound")), http.StatusBadRequest)
// 	ErrUnkwonError      = SetStatusCode(SetPointer(errors.New("unKnownError")), http.StatusBadRequest)
// 	ErrUnauthorizedUser = SetStatusCode(SetPointer(errors.New("unauthorizedUser")), http.StatusBadRequest)
// )
