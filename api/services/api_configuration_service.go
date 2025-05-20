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
	appmodel "app-name/app/model"
	"app-name/broker"
	dbhelper "app-name/db/helper"
	"context"
	"errors"
	"fmt"
	"net/http"
)

// ConfigurationAPIService is a service that implements the logic for the ConfigurationAPIServicer
// This service should implement the business logic for every endpoint for the ConfigurationAPI API.
// Include any external packages or services that will be required by this service.
type ConfigurationAPIService struct {
}

// NewConfigurationAPIService creates a default api service
func NewConfigurationAPIService() apiserver.ConfigurationAPIServicer {
	return &ConfigurationAPIService{}
}

func (s *ConfigurationAPIService) GetConfigurations(ctx context.Context) (apiserver.ImplResponse, error) {
	appConfigs, err := dbhelper.GetConfigs(ctx)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	var configs []apiserver.Configuration
	for _, appConfig := range appConfigs {
		configs = append(configs, toAPIConfig(appConfig))
	}
	return apiserver.Response(http.StatusOK, configs), nil
}

func (s *ConfigurationAPIService) PostConfiguration(ctx context.Context, config apiserver.Configuration) (apiserver.ImplResponse, error) {
	appConfig := toAppConfig(config)
	if err := broker.TestAuthentication(appConfig); err != nil {
		return apiserver.ImplResponse{Code: http.StatusBadRequest}, fmt.Errorf("testing authentication: %v", err)
	}
	insertedConfig, err := dbhelper.UpsertConfig(ctx, appConfig)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusCreated, toAPIConfig(insertedConfig)), nil
}

func (s *ConfigurationAPIService) GetConfigurationById(ctx context.Context, configId int64) (apiserver.ImplResponse, error) {
	config, err := dbhelper.GetConfig(ctx, configId)
	if errors.Is(err, dbhelper.ErrNotFound) {
		return apiserver.ImplResponse{Code: http.StatusNotFound}, nil
	}
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusOK, toAPIConfig(config)), nil
}

func (s *ConfigurationAPIService) PutConfigurationById(ctx context.Context, configId int64, config apiserver.Configuration) (apiserver.ImplResponse, error) {
	config.Id = &configId
	appConfig := toAppConfig(config)
	if err := broker.TestAuthentication(appConfig); err != nil {
		return apiserver.ImplResponse{Code: http.StatusBadRequest}, fmt.Errorf("testing authentication: %v", err)
	}
	upsertedConfig, err := dbhelper.UpsertConfig(ctx, appConfig)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusCreated, toAPIConfig(upsertedConfig)), nil
}

func (s *ConfigurationAPIService) DeleteConfigurationById(ctx context.Context, configId int64) (apiserver.ImplResponse, error) {
	err := dbhelper.DeleteConfig(ctx, configId)
	if errors.Is(err, dbhelper.ErrNotFound) {
		return apiserver.ImplResponse{Code: http.StatusNotFound}, nil
	}
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.ImplResponse{Code: http.StatusNoContent}, nil
}

func toAPIConfig(appConfig appmodel.Configuration) apiserver.Configuration {
	return apiserver.Configuration{
		Id:                &appConfig.Id,
		ApiAccessChangeMe: appConfig.ApiAccessChangeMe,
		Enable:            &appConfig.Enable,
		RefreshInterval:   appConfig.RefreshInterval,
		RequestTimeout:    &appConfig.RequestTimeout,
		AssetFilter:       toAPIAssetFilter(appConfig.AssetFilter),
		Active:            &appConfig.Active,
		ProjectIDs:        &appConfig.ProjectIDs,
		UserId:            &appConfig.UserId,
	}
}

func toAPIAssetFilter(appAF [][]appmodel.FilterRule) (result [][]apiserver.FilterRule) {
	for _, outer := range appAF {
		var innerResult []apiserver.FilterRule
		for _, fr := range outer {
			innerResult = append(innerResult, apiserver.FilterRule{
				Parameter: fr.Parameter,
				Regex:     fr.Regex,
			})
		}
		result = append(result, innerResult)
	}
	return result
}

func toAppConfig(apiConfig apiserver.Configuration) (appConfig appmodel.Configuration) {
	appConfig.ApiAccessChangeMe = apiConfig.ApiAccessChangeMe

	if apiConfig.Id != nil {
		appConfig.Id = *apiConfig.Id
	}
	appConfig.RefreshInterval = apiConfig.RefreshInterval
	if apiConfig.RequestTimeout != nil {
		appConfig.RequestTimeout = *apiConfig.RequestTimeout
	}
	if apiConfig.AssetFilter != nil {
		appConfig.AssetFilter = toAppAssetFilter(apiConfig.AssetFilter)
	}
	if apiConfig.Active != nil {
		appConfig.Active = *apiConfig.Active
	}
	if apiConfig.Enable != nil {
		appConfig.Enable = *apiConfig.Enable
	}
	if apiConfig.ProjectIDs != nil {
		appConfig.ProjectIDs = *apiConfig.ProjectIDs
	}
	return appConfig
}

func toAppAssetFilter(apiAF [][]apiserver.FilterRule) (result [][]appmodel.FilterRule) {
	for _, outer := range apiAF {
		var innerResult []appmodel.FilterRule
		for _, fr := range outer {
			innerResult = append(innerResult, appmodel.FilterRule{
				Parameter: fr.Parameter,
				Regex:     fr.Regex,
			})
		}
		result = append(result, innerResult)
	}
	return result
}
