package webkeys

import (
	"io/ioutil"

	"github.com/square/go-jose/json"
	"gopkg.in/square/go-jose.v2"
)

type FileJwksService struct {
	config FileConfig
	keys   jose.JSONWebKeySet
}

func NewFileJwksService(config FileConfig) (WebKeysService, error) {
	service := &FileJwksService{config, jose.JSONWebKeySet{}}

	keys, err := service.Load()
	if err != nil {
		return nil, err
	}

	service.keys = *keys

	return service, nil
}

func (s *FileJwksService) Load() (*jose.JSONWebKeySet, error) {

	allKeys := []jose.JSONWebKey{}

	for _, path := range s.config.Paths {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		jsonWebKeySet := jose.JSONWebKeySet{}
		if err = json.Unmarshal(contents, &jsonWebKeySet); err != nil {
			return nil, err
		}

		allKeys = append(allKeys, jsonWebKeySet.Keys...)
	}

	allKeySet := jose.JSONWebKeySet{
		Keys: allKeys,
	}

	return &allKeySet, nil
}

func (s *FileJwksService) Keys() (*jose.JSONWebKeySet, error) {
	return &s.keys, nil
}
