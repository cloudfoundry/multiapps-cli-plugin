// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetSharedDomainsParams creates a new GetSharedDomainsParams object
// with the default values initialized.
func NewGetSharedDomainsParams() *GetSharedDomainsParams {

	return &GetSharedDomainsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetSharedDomainsParamsWithTimeout creates a new GetSharedDomainsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetSharedDomainsParamsWithTimeout(timeout time.Duration) *GetSharedDomainsParams {

	return &GetSharedDomainsParams{

		timeout: timeout,
	}
}

// NewGetSharedDomainsParamsWithContext creates a new GetSharedDomainsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetSharedDomainsParamsWithContext(ctx context.Context) *GetSharedDomainsParams {

	return &GetSharedDomainsParams{

		Context: ctx,
	}
}

// NewGetSharedDomainsParamsWithHTTPClient creates a new GetSharedDomainsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetSharedDomainsParamsWithHTTPClient(client *http.Client) *GetSharedDomainsParams {

	return &GetSharedDomainsParams{
		HTTPClient: client,
	}
}

/*GetSharedDomainsParams contains all the parameters to send to the API endpoint
for the get shared domains operation typically these are written to a http.Request
*/
type GetSharedDomainsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get shared domains params
func (o *GetSharedDomainsParams) WithTimeout(timeout time.Duration) *GetSharedDomainsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get shared domains params
func (o *GetSharedDomainsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get shared domains params
func (o *GetSharedDomainsParams) WithContext(ctx context.Context) *GetSharedDomainsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get shared domains params
func (o *GetSharedDomainsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get shared domains params
func (o *GetSharedDomainsParams) WithHTTPClient(client *http.Client) *GetSharedDomainsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get shared domains params
func (o *GetSharedDomainsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetSharedDomainsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
