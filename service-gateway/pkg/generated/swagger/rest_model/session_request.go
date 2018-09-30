// Code generated by go-swagger; DO NOT EDIT.

package rest_model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SessionRequest session request
// swagger:model SessionRequest
type SessionRequest struct {

	// session ID
	// Required: true
	SessionID *string `json:"sessionID"`
}

// Validate validates this session request
func (m *SessionRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSessionID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SessionRequest) validateSessionID(formats strfmt.Registry) error {

	if err := validate.Required("sessionID", "body", m.SessionID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SessionRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SessionRequest) UnmarshalBinary(b []byte) error {
	var res SessionRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
