package errors

import "net/http"

//NotFound represents a HTTP status code 404 error
type NotFound struct {
	httpError
}

//NewNotFound constructs a NotFound HTTP error
func NewNotFound(messages ...string) *NotFound {
	return &NotFound{
		httpError{
			http.StatusNotFound,
			new("Resource not found.", messages...),
		},
	}
}

//BadRequest represents a HTTP status code 400 error
type BadRequest struct {
	httpError
}

//NewBadRequest constructs a BadRequest HTTP error
func NewBadRequest(messages ...string) *BadRequest {
	return &BadRequest{
		httpError{
			http.StatusBadRequest,
			new("Invalid request sent.", messages...),
		},
	}
}

//InternalServer represents a HTTP status code 500 error
type InternalServer struct {
	httpError
}

//NewInternalServer constructs a InternalServer HTTP error
func NewInternalServer(messages ...string) *InternalServer {
	return &InternalServer{
		httpError{
			http.StatusInternalServerError,
			new("Something was wrong with server.", messages...),
		},
	}
}

//UnprocessableEntity represents a HTTP status code 422 error
type UnprocessableEntity struct {
	httpError
}

//NewUnprocessableEntity constructs a UnprocessableEntity HTTP error
func NewUnprocessableEntity(messages ...string) *UnprocessableEntity {
	return &UnprocessableEntity{
		httpError{
			http.StatusUnprocessableEntity,
			new("Unprocessable entity.", messages...),
		},
	}
}

//MethodNotAllowed represents a HTTP status code 422 error
type MethodNotAllowed struct {
	httpError
}

//NewMethodNotAllowed constructs a MethodNotAllowed HTTP error
func NewMethodNotAllowed(messages ...string) *MethodNotAllowed {
	return &MethodNotAllowed{
		httpError{
			http.StatusMethodNotAllowed,
			new("Method not alowed.", messages...),
		},
	}
}
