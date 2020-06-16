package webkeys

import "gopkg.in/square/go-jose.v2"

type staticWebKeysService struct {
	keys *jose.JSONWebKeySet
}

func NewStaticJwksService(keys *jose.JSONWebKeySet) WebKeysService {
	return &staticWebKeysService{keys: keys}
}

func (m *staticWebKeysService) Keys() (*jose.JSONWebKeySet, error) {
	return m.keys, nil
}
