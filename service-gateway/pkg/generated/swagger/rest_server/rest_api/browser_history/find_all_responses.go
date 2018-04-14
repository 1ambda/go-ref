// Code generated by go-swagger; DO NOT EDIT.

package browser_history

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	rest_model "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
)

// FindAllOKCode is the HTTP code returned for type FindAllOK
const FindAllOKCode int = 200

/*FindAllOK BrowserHistory records with pagination info

swagger:response findAllOK
*/
type FindAllOK struct {

	/*
	  In: Body
	*/
	Payload *rest_model.BrowserHistoryWithPagination `json:"body,omitempty"`
}

// NewFindAllOK creates FindAllOK with default headers values
func NewFindAllOK() *FindAllOK {

	return &FindAllOK{}
}

// WithPayload adds the payload to the find all o k response
func (o *FindAllOK) WithPayload(payload *rest_model.BrowserHistoryWithPagination) *FindAllOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find all o k response
func (o *FindAllOK) SetPayload(payload *rest_model.BrowserHistoryWithPagination) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindAllOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*FindAllDefault generic error response

swagger:response findAllDefault
*/
type FindAllDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *rest_model.Error `json:"body,omitempty"`
}

// NewFindAllDefault creates FindAllDefault with default headers values
func NewFindAllDefault(code int) *FindAllDefault {
	if code <= 0 {
		code = 500
	}

	return &FindAllDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the find all default response
func (o *FindAllDefault) WithStatusCode(code int) *FindAllDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the find all default response
func (o *FindAllDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the find all default response
func (o *FindAllDefault) WithPayload(payload *rest_model.Error) *FindAllDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find all default response
func (o *FindAllDefault) SetPayload(payload *rest_model.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindAllDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
