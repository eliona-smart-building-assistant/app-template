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
	"net/http"
	"template/apiserver"
)

// CustomizationAPIService is a service that implements the logic for the CustomizationAPIServicer
// This service should implement the business logic for every endpoint for the CustomizationAPI API.
// Include any external packages or services that will be required by this service.
type CustomizationAPIService struct {
}

// NewCustomizationAPIService creates a default api service
func NewCustomizationAPIService() apiserver.CustomizationAPIServicer {
	return &CustomizationAPIService{}
}

// GetDashboardTemplateByName - Get a full dashboard template
func (s *CustomizationAPIService) GetDashboardTemplateByName(ctx context.Context, dashboardTemplateName string, projectId string) (apiserver.ImplResponse, error) {
	if dashboardTemplateName == "Template" {
		return apiserver.ImplResponse{Code: http.StatusNotImplemented}, nil
	} else {
		return apiserver.ImplResponse{Code: http.StatusNotFound}, nil
	}
}
