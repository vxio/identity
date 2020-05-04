package jwks

import (
	"fmt"
	"net/http"

	"github.com/square/go-jose/json"
	"gopkg.in/square/go-jose.v2"
)

type WebJwksService struct {
	client  *http.Client
	jwksURI string
}

func NewWebJwksService(jwksURI string) JwksService {
	return &WebJwksService{
		client:  &http.Client{},
		jwksURI: jwksURI,
	}
}

func (s *WebJwksService) FetchJwks() (*jose.JSONWebKeySet, error) {
	resp, err := s.client.Get(s.jwksURI)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to load JWKS due to response code: %d", resp.StatusCode)
	}

	jsonWebKeySet := new(jose.JSONWebKeySet)
	if err = json.NewDecoder(resp.Body).Decode(jsonWebKeySet); err != nil {
		return nil, err
	}

	return jsonWebKeySet, nil
}
