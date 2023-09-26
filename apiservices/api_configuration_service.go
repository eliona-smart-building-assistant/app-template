//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package apiservices

import (
	"context"
	"errors"
	"net/http"
	"template/apiserver"
)

// ConfigurationApiService is a service that implements the logic for the ConfigurationApiServicer
// This service should implement the business logic for every endpoint for the ConfigurationApi API.
// Include any external packages or services that will be required by this service.
type ConfigurationApiService struct {
}

// NewConfigurationApiService creates a default api service
func NewConfigurationApiService() apiserver.ConfigurationAPIServicer {
	return &ConfigurationApiService{}
}

// GetConfigurations - Get example configurations
func (s *ConfigurationApiService) GetConfigurations(ctx context.Context) (apiserver.ImplResponse, error) {
	// TODO - update GetConfigurations with the required logic for this service method.
	// Add api_configuration_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Configuration{}) or use other options such as http.Ok ...
	//return Response(200, []Configuration{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("GetConfigurations method not implemented")
}

// PostConfiguration - Creates an example configuration
func (s *ConfigurationApiService) PostConfiguration(ctx context.Context, configuration apiserver.Configuration) (apiserver.ImplResponse, error) {
	// TODO - update PostConfiguration with the required logic for this service method.
	// Add api_configuration_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, Configuration{}) or use other options such as http.Ok ...
	//return Response(201, Configuration{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("PostConfiguration method not implemented")
}
