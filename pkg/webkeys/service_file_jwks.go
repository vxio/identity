package webkeys

import (
	"io/ioutil"

	"github.com/square/go-jose/json"
	"gopkg.in/square/go-jose.v2"
)

type FileJwksService struct {
	filePath string
}

func NewFileJwksService(filePath string) WebKeysService {
	return &FileJwksService{
		filePath: filePath,
	}
}

func (s *FileJwksService) FetchJwks() (*jose.JSONWebKeySet, error) {

	contents, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	jsonWebKeySet := jose.JSONWebKeySet{}
	if err = json.Unmarshal(contents, &jsonWebKeySet); err != nil {
		return nil, err
	}

	return &jsonWebKeySet, nil
}
