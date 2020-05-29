package webkeys

import (
	"fmt"
	"net/http"

	"github.com/square/go-jose/json"
	"gopkg.in/square/go-jose.v2"
)

type HTTPJwksService struct {
	config HttpConfig
	client *http.Client
	keys   jose.JSONWebKeySet
}

func NewHTTPJwksService(config HttpConfig, client *http.Client) (WebKeysService, error) {
	if client == nil {
		client = &http.Client{}
	}

	service := &HTTPJwksService{
		client: &http.Client{},
		config: config,
		keys:   jose.JSONWebKeySet{},
	}

	keys, err := service.Load()
	if err != nil {
		return nil, err
	}

	service.keys = *keys

	return service, nil
}

func (s *HTTPJwksService) Load() (*jose.JSONWebKeySet, error) {

	allKeys := []jose.JSONWebKey{}

	for _, url := range s.config.URLs {

		resp, err := s.client.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Failed to load JWKS due to response code: %d", resp.StatusCode)
		}

		jsonWebKeySet := jose.JSONWebKeySet{}
		if err = json.NewDecoder(resp.Body).Decode(&jsonWebKeySet); err != nil {
			return nil, err
		}

		allKeys = append(allKeys, jsonWebKeySet.Keys...)
	}

	allKeySet := jose.JSONWebKeySet{
		Keys: allKeys,
	}

	return &allKeySet, nil
}

func (s *HTTPJwksService) Keys() (*jose.JSONWebKeySet, error) {
	return &s.keys, nil
}
