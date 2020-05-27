package webkeys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/square/go-jose.v2"
)

type GenerateJwksService struct {
	Private jose.JSONWebKey
	Public  jose.JSONWebKey
}

func NewGenerateJwksService() (*GenerateJwksService, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	priv := jose.JSONWebKey{
		Key:       key,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}

	thumb, err := priv.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, errors.New("Unable to generate fingerprint")
	}

	kid := base64.URLEncoding.EncodeToString(thumb)[0:10]
	priv.KeyID = kid

	pub := jose.JSONWebKey{
		Key:       key.Public(),
		KeyID:     kid,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}

	if priv.IsPublic() || !pub.IsPublic() || !priv.Valid() || !pub.Valid() {
		return nil, errors.New("Generated keys are invalid")
	}

	return &GenerateJwksService{
		Public:  pub,
		Private: priv,
	}, nil
}

func (s *GenerateJwksService) FetchJwks() (*jose.JSONWebKeySet, error) {
	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{s.Public, s.Private},
	}

	return &jwks, nil
}

func (s *GenerateJwksService) Save(path string) error {
	pubJwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{s.Public}}
	pubJson, err := json.Marshal(pubJwks)
	if err != nil {
		return err
	}

	privJwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{s.Private}}
	privJson, err := json.Marshal(privJwks)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s/jwk-sig", path)
	pubFile := fmt.Sprintf("%s-pub.json", name)
	privFile := fmt.Sprintf("%s-priv.json", name)

	if err := ioutil.WriteFile(pubFile, pubJson, 0444); err != nil {
		return err
	}

	if err := ioutil.WriteFile(privFile, privJson, 0400); err != nil {
		return err
	}

	return nil
}
