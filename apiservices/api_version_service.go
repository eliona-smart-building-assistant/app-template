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
	"os"
	"template/apiserver"

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"gopkg.in/yaml.v3"
)

// VersionApiService is a service that implements the logic for the VersionApiServicer
// This service should implement the business logic for every endpoint for the VersionApi API.
// Include any external packages or services that will be required by this service.
type VersionApiService struct {
}

// NewVersionApiService creates a default api service
func NewVersionApiService() apiserver.VersionApiServicer {
	return &VersionApiService{}
}

// GetOpenAPI - OpenAPI specification for this API version
func (s *VersionApiService) GetOpenAPI(ctx context.Context) (apiserver.ImplResponse, error) {
	bytes, err := os.ReadFile("openapi.yaml")
	if err != nil {
		log.Error("services", "%s: %v", "GetOpenAPI", err)
		return apiserver.ImplResponse{Code: http.StatusNotFound}, err
	}
	var body interface{}
	err = yaml.Unmarshal(bytes, &body)
	if err != nil {
		log.Error("services", "%s: %v", "GetOpenAPI", err)
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusOK, body), nil
}

var BuildTimestamp string // injected during linking, see Dockerfile
var GitCommit string      // injected during linking, see Dockerfile

// GetVersion - Version of the API
func (s *VersionApiService) GetVersion(ctx context.Context) (apiserver.ImplResponse, error) {
	return apiserver.Response(http.StatusOK, common.Ptr(version())), nil
}

func version() map[string]any {
	return map[string]any{
		"timestamp": BuildTimestamp,
		"commit":    GitCommit,
	}
}
