// Code generated by go-swagger; DO NOT EDIT.

package silence

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetSilencesParams creates a new GetSilencesParams object
// with the default values initialized.
func NewGetSilencesParams() *GetSilencesParams {
	var ()
	return &GetSilencesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetSilencesParamsWithTimeout creates a new GetSilencesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetSilencesParamsWithTimeout(timeout time.Duration) *GetSilencesParams {
	var ()
	return &GetSilencesParams{

		timeout: timeout,
	}
}

// NewGetSilencesParamsWithContext creates a new GetSilencesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetSilencesParamsWithContext(ctx context.Context) *GetSilencesParams {
	var ()
	return &GetSilencesParams{

		Context: ctx,
	}
}

// NewGetSilencesParamsWithHTTPClient creates a new GetSilencesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetSilencesParamsWithHTTPClient(client *http.Client) *GetSilencesParams {
	var ()
	return &GetSilencesParams{
		HTTPClient: client,
	}
}

/*GetSilencesParams contains all the parameters to send to the API endpoint
for the get silences operation typically these are written to a http.Request
*/
type GetSilencesParams struct {

	/*Filter
	  A list of matchers to filter silences by

	*/
	Filter []string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get silences params
func (o *GetSilencesParams) WithTimeout(timeout time.Duration) *GetSilencesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get silences params
func (o *GetSilencesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get silences params
func (o *GetSilencesParams) WithContext(ctx context.Context) *GetSilencesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get silences params
func (o *GetSilencesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get silences params
func (o *GetSilencesParams) WithHTTPClient(client *http.Client) *GetSilencesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get silences params
func (o *GetSilencesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFilter adds the filter to the get silences params
func (o *GetSilencesParams) WithFilter(filter []string) *GetSilencesParams {
	o.SetFilter(filter)
	return o
}

// SetFilter adds the filter to the get silences params
func (o *GetSilencesParams) SetFilter(filter []string) {
	o.Filter = filter
}

// WriteToRequest writes these params to a swagger request
func (o *GetSilencesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	valuesFilter := o.Filter

	joinedFilter := swag.JoinByFormat(valuesFilter, "")
	// query array param filter
	if err := r.SetQueryParam("filter", joinedFilter...); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
