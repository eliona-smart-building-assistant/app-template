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
	"encoding/json"
	"io"
	"net/http"
	"os"
	"template/apiserver"

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"gopkg.in/yaml.v3"
)

// VersionAPIService is a service that implements the logic for the VersionAPIServicer
// This service should implement the business logic for every endpoint for the VersionAPI API.
// Include any external packages or services that will be required by this service.
type VersionAPIService struct {
}

// NewVersionAPIService creates a default api service
func NewVersionAPIService() apiserver.VersionAPIServicer {
	return &VersionAPIService{}
}

// GetOpenAPI - OpenAPI specification for this API version
func (s *VersionAPIService) GetOpenAPI(ctx context.Context) (apiserver.ImplResponse, error) {
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
	if err := tryJSONEncoding(body); err != nil {
		log.Error("services", "openapi file not encodable to JSON: %v", err)
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, nil
	}
	return apiserver.Response(http.StatusOK, body), nil
}

func tryJSONEncoding(i interface{}) error {
	return json.NewEncoder(io.Discard).Encode(i)
}

var BuildTimestamp string // injected during linking, see Dockerfile
var GitCommit string      // injected during linking, see Dockerfile

// GetVersion - Version of the API
func (s *VersionAPIService) GetVersion(ctx context.Context) (apiserver.ImplResponse, error) {
	return apiserver.Response(http.StatusOK, common.Ptr(version())), nil
}

func version() map[string]any {
	return map[string]any{
		"timestamp": BuildTimestamp,
		"commit":    GitCommit,
	}
}
