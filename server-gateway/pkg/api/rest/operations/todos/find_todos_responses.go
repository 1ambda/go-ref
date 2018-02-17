// Code generated by go-swagger; DO NOT EDIT.

package todos

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	model "github.com/1ambda/go-ref/server-gateway/pkg/api/model"
)

// FindTodosOKCode is the HTTP code returned for type FindTodosOK
const FindTodosOKCode int = 200

/*FindTodosOK list the todo operations

swagger:response findTodosOK
*/
type FindTodosOK struct {

	/*
	  In: Body
	*/
	Payload model.FindTodosOKBody `json:"body,omitempty"`
}

// NewFindTodosOK creates FindTodosOK with default headers values
func NewFindTodosOK() *FindTodosOK {

	return &FindTodosOK{}
}

// WithPayload adds the payload to the find todos o k response
func (o *FindTodosOK) WithPayload(payload model.FindTodosOKBody) *FindTodosOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find todos o k response
func (o *FindTodosOK) SetPayload(payload model.FindTodosOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindTodosOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make(model.FindTodosOKBody, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*FindTodosDefault generic error response

swagger:response findTodosDefault
*/
type FindTodosDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *model.Error `json:"body,omitempty"`
}

// NewFindTodosDefault creates FindTodosDefault with default headers values
func NewFindTodosDefault(code int) *FindTodosDefault {
	if code <= 0 {
		code = 500
	}

	return &FindTodosDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the find todos default response
func (o *FindTodosDefault) WithStatusCode(code int) *FindTodosDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the find todos default response
func (o *FindTodosDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the find todos default response
func (o *FindTodosDefault) WithPayload(payload *model.Error) *FindTodosDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find todos default response
func (o *FindTodosDefault) SetPayload(payload *model.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindTodosDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
