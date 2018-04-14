// Code generated by go-swagger; DO NOT EDIT.

package session

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ValidateOrGenerateHandlerFunc turns a function with the right signature into a validate or generate handler
type ValidateOrGenerateHandlerFunc func(ValidateOrGenerateParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ValidateOrGenerateHandlerFunc) Handle(params ValidateOrGenerateParams) middleware.Responder {
	return fn(params)
}

// ValidateOrGenerateHandler interface for that can handle valid validate or generate params
type ValidateOrGenerateHandler interface {
	Handle(ValidateOrGenerateParams) middleware.Responder
}

// NewValidateOrGenerate creates a new http.Handler for the validate or generate operation
func NewValidateOrGenerate(ctx *middleware.Context, handler ValidateOrGenerateHandler) *ValidateOrGenerate {
	return &ValidateOrGenerate{Context: ctx, Handler: handler}
}

/*ValidateOrGenerate swagger:route POST /session session validateOrGenerate

ValidateOrGenerate validate or generate API

*/
type ValidateOrGenerate struct {
	Context *middleware.Context
	Handler ValidateOrGenerateHandler
}

func (o *ValidateOrGenerate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewValidateOrGenerateParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}