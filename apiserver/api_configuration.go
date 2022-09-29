/*
 * App template API
 *
 * API to access and configure the app template
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ConfigurationApiController binds http requests to an api service and writes the service results to the http response
type ConfigurationApiController struct {
	service      ConfigurationApiServicer
	errorHandler ErrorHandler
}

// ConfigurationApiOption for how the controller is set up.
type ConfigurationApiOption func(*ConfigurationApiController)

// WithConfigurationApiErrorHandler inject ErrorHandler into controller
func WithConfigurationApiErrorHandler(h ErrorHandler) ConfigurationApiOption {
	return func(c *ConfigurationApiController) {
		c.errorHandler = h
	}
}

// NewConfigurationApiController creates a default api controller
func NewConfigurationApiController(s ConfigurationApiServicer, opts ...ConfigurationApiOption) Router {
	controller := &ConfigurationApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ConfigurationApiController
func (c *ConfigurationApiController) Routes() Routes {
	return Routes{
		{
			"GetExamples",
			strings.ToUpper("Get"),
			"/v1/examples",
			c.GetExamples,
		},
		{
			"PostExample",
			strings.ToUpper("Post"),
			"/v1/examples",
			c.PostExample,
		},
	}
}

// GetExamples - Get example configuration
func (c *ConfigurationApiController) GetExamples(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetExamples(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// PostExample - Creates an example configuration
func (c *ConfigurationApiController) PostExample(w http.ResponseWriter, r *http.Request) {
	exampleParam := Example{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&exampleParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertExampleRequired(exampleParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.PostExample(r.Context(), exampleParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}