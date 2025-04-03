//  This file is part of the Eliona project.
//  Copyright © 2025 IoTEC AG. All Rights Reserved.
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
	apiserver "app-name/api/generated"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

func (s *VersionAPIService) GetOpenAPI(ctx context.Context) (apiserver.ImplResponse, error) {
	bytes, err := os.ReadFile("openapi.yaml")
	if err != nil {
		log.Error("services", "GetOpenAPI - Error reading openapi.yaml: %v", err)
		return apiserver.ImplResponse{Code: http.StatusNotFound}, err
	}

	var body interface{}
	err = yaml.Unmarshal(bytes, &body)
	if err != nil {
		log.Error("services", "GetOpenAPI - Error unmarshalling YAML: %v", err)
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}

	if err := tryJSONEncoding(body); err != nil {
		log.Error("services", "GetOpenAPI - openapi file not encodable to JSON: %v", err)
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, nil
	}

	return apiserver.Response(http.StatusOK, body), nil
}

func tryJSONEncoding(i interface{}) error {
	if err := json.NewEncoder(io.Discard).Encode(i); err != nil {
		unsupportedTypeErr := checkForUnsupportedTypes(i, "")
		if unsupportedTypeErr != nil {
			return fmt.Errorf("encoding json: %v, unsupportedTypeErr: %v", err, unsupportedTypeErr)
		}
	}
	return nil
}

func checkForUnsupportedTypes(i interface{}, path string) error {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		// Log the error immediately because this is not JSON-compatible
		err := fmt.Errorf("unsupported type map[interface{}]interface{} at path '%s'", path)
		log.Error("services", "%v", err)
		return err
	case []interface{}:
		for idx, v := range x {
			newPath := fmt.Sprintf("%s[%d]", path, idx)
			if err := checkForUnsupportedTypes(v, newPath); err != nil {
				return err
			}
		}
	case map[string]interface{}:
		for k, v := range x {
			newPath := fmt.Sprintf("%s.%s", path, k)
			if err := checkForUnsupportedTypes(v, newPath); err != nil {
				return err
			}
		}
	case string, int, float64, bool, nil:
		// These are valid JSON types, no action needed
		return nil
	default:
		// Log any other unsupported types
		err := fmt.Errorf("unsupported type at path '%s': %T", path, x)
		log.Error("services", "%v", err)
		return err
	}
	return nil
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
