// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

// Package api implements the app's own REST API defined in openapi.yaml.
package api

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"app-name/internal/apiserver"
	"app-name/internal/store"

	"sigs.k8s.io/yaml"
)

// Set at build time with -ldflags "-X app-name/internal/api.BuildTimestamp=... -X app-name/internal/api.GitCommit=...".
var (
	BuildTimestamp = "unknown"
	GitCommit      = "unknown"
)

// Server implements apiserver.StrictServerInterface.
type Server struct {
	store   store.Store
	openAPI map[string]any
}

var _ apiserver.StrictServerInterface = (*Server)(nil)

// NewServer needs the content of openapi.yaml to serve it as JSON.
func NewServer(s store.Store, openAPIYAML []byte) (*Server, error) {
	var spec map[string]any
	if err := yaml.Unmarshal(openAPIYAML, &spec); err != nil {
		return nil, fmt.Errorf("openapi.yaml: %w", err)
	}
	return &Server{store: s, openAPI: spec}, nil
}

func (s *Server) ListConfigurations(ctx context.Context, _ apiserver.ListConfigurationsRequestObject) (apiserver.ListConfigurationsResponseObject, error) {
	configs, err := s.store.ListConfigs(ctx)
	if err != nil {
		return nil, err
	}
	out := make(apiserver.ListConfigurations200JSONResponse, 0, len(configs))
	for _, c := range configs {
		out = append(out, toAPI(c))
	}
	return out, nil
}

func (s *Server) CreateConfiguration(ctx context.Context, req apiserver.CreateConfigurationRequestObject) (apiserver.CreateConfigurationResponseObject, error) {
	c, err := fromAPI(*req.Body, store.Config{})
	if err != nil {
		return apiserver.CreateConfiguration400JSONResponse{BadRequestJSONResponse: apiserver.BadRequestJSONResponse{Message: err.Error()}}, nil
	}
	created, err := s.store.CreateConfig(ctx, c)
	if err != nil {
		return nil, err
	}
	return apiserver.CreateConfiguration201JSONResponse(toAPI(created)), nil
}

func (s *Server) GetConfiguration(ctx context.Context, req apiserver.GetConfigurationRequestObject) (apiserver.GetConfigurationResponseObject, error) {
	c, err := s.store.GetConfig(ctx, req.ConfigId)
	if errors.Is(err, store.ErrNotFound) {
		return apiserver.GetConfiguration404JSONResponse{NotFoundJSONResponse: notFound(req.ConfigId)}, nil
	}
	if err != nil {
		return nil, err
	}
	return apiserver.GetConfiguration200JSONResponse(toAPI(c)), nil
}

func (s *Server) UpdateConfiguration(ctx context.Context, req apiserver.UpdateConfigurationRequestObject) (apiserver.UpdateConfigurationResponseObject, error) {
	existing, err := s.store.GetConfig(ctx, req.ConfigId)
	if errors.Is(err, store.ErrNotFound) {
		return apiserver.UpdateConfiguration404JSONResponse{NotFoundJSONResponse: notFound(req.ConfigId)}, nil
	}
	if err != nil {
		return nil, err
	}
	c, err := fromAPI(*req.Body, existing)
	if err != nil {
		return apiserver.UpdateConfiguration400JSONResponse{BadRequestJSONResponse: apiserver.BadRequestJSONResponse{Message: err.Error()}}, nil
	}
	c.ID = req.ConfigId
	updated, err := s.store.UpdateConfig(ctx, c)
	if errors.Is(err, store.ErrNotFound) {
		return apiserver.UpdateConfiguration404JSONResponse{NotFoundJSONResponse: notFound(req.ConfigId)}, nil
	}
	if err != nil {
		return nil, err
	}
	return apiserver.UpdateConfiguration200JSONResponse(toAPI(updated)), nil
}

func (s *Server) DeleteConfiguration(ctx context.Context, req apiserver.DeleteConfigurationRequestObject) (apiserver.DeleteConfigurationResponseObject, error) {
	err := s.store.DeleteConfig(ctx, req.ConfigId)
	if errors.Is(err, store.ErrNotFound) {
		return apiserver.DeleteConfiguration404JSONResponse{NotFoundJSONResponse: notFound(req.ConfigId)}, nil
	}
	if err != nil {
		return nil, err
	}
	return apiserver.DeleteConfiguration204Response{}, nil
}

func (s *Server) GetVersion(context.Context, apiserver.GetVersionRequestObject) (apiserver.GetVersionResponseObject, error) {
	return apiserver.GetVersion200JSONResponse{Timestamp: BuildTimestamp, Commit: GitCommit}, nil
}

func (s *Server) GetOpenAPI(context.Context, apiserver.GetOpenAPIRequestObject) (apiserver.GetOpenAPIResponseObject, error) {
	return apiserver.GetOpenAPI200JSONResponse(s.openAPI), nil
}

func notFound(id int64) apiserver.NotFoundJSONResponse {
	return apiserver.NotFoundJSONResponse{Message: fmt.Sprintf("configuration %d not found", id)}
}

// toAPI never returns the vendor API key; the schema marks it writeOnly.
func toAPI(c store.Config) apiserver.Configuration {
	return apiserver.Configuration{
		Id:              &c.ID,
		Url:             c.URL,
		Enable:          &c.Enable,
		RefreshInterval: &c.RefreshInterval,
		RequestTimeout:  &c.RequestTimeout,
	}
}

// fromAPI applies the request on top of base; omitted fields keep the base value or the schema default.
func fromAPI(in apiserver.Configuration, base store.Config) (store.Config, error) {
	c := base
	if c.ID == 0 {
		c = store.Config{Enable: true, RefreshInterval: 60, RequestTimeout: 30}
	}

	u, err := url.Parse(in.Url)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return store.Config{}, fmt.Errorf("url must be an absolute http or https URL")
	}
	c.URL = in.Url
	if in.ApiKey != nil {
		c.APIKey = *in.ApiKey
	}
	if in.Enable != nil {
		c.Enable = *in.Enable
	}
	if in.RefreshInterval != nil {
		if *in.RefreshInterval < 10 {
			return store.Config{}, fmt.Errorf("refreshInterval must be at least 10 seconds")
		}
		c.RefreshInterval = *in.RefreshInterval
	}
	if in.RequestTimeout != nil {
		if *in.RequestTimeout < 1 {
			return store.Config{}, fmt.Errorf("requestTimeout must be at least 1 second")
		}
		c.RequestTimeout = *in.RequestTimeout
	}
	return c, nil
}
