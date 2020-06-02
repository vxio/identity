package webkeys_test

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/moov-io/identity/pkg/logging"
	. "github.com/moov-io/identity/pkg/webkeys"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateKeys(t *testing.T) {
	a, s := Setup(t)

	config := WebKeysConfig{}
	service, err := NewWebKeysService(s.logger, config)
	a.Nil(err)

	keys, err := service.Keys()
	a.Nil(err)
	a.Len(keys.Keys, 2)

	generator, ok := service.(*GenerateJwksService)
	a.True(ok)

	dir, err := ioutil.TempDir("", "test")
	a.Nil(err)

	generator.Save(dir)
	t.Cleanup(func() {
		syscall.Rmdir(dir)
	})
}

func Test_HttpKeys(t *testing.T) {
	a, s := Setup(t)

	generator, err := NewGenerateJwksService()
	a.Nil(err)

	controller := NewJWKSController(generator)

	router := mux.NewRouter()
	controller.AppendRoutes(router)

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	// testing...
	config := WebKeysConfig{
		HTTP: &HttpConfig{
			URLs: []string{
				server.URL + controller.WellKnownJwksPath(),
			},
		},
	}

	service, err := NewWebKeysService(s.logger, config)
	a.Nil(err)

	keys, err := service.Keys()
	a.Nil(err)
	a.Len(keys.Keys, 1)

	// Fetching from the server will only produce public keys
	for _, v := range keys.Keys {
		a.True(v.IsPublic())
	}
}

func Test_HttpKeys_NotFound(t *testing.T) {
	a, s := Setup(t)

	router := mux.NewRouter()
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	// testing...
	config := WebKeysConfig{
		HTTP: &HttpConfig{
			URLs: []string{
				server.URL + "/.well-known/jwks.json",
			},
		},
	}

	_, err := NewWebKeysService(s.logger, config)
	a.NotNil(err)
}

func Test_FileKeys(t *testing.T) {
	a, s := Setup(t)

	generator, err := NewGenerateJwksService()
	a.Nil(err)

	dir, err := ioutil.TempDir("", "webkeystesting-"+uuid.New().String())
	a.Nil(err)

	generator.Save(dir)
	t.Cleanup(func() {
		syscall.Rmdir(dir)
	})

	// testing...
	config := WebKeysConfig{
		File: &FileConfig{
			Paths: []string{
				dir + "/jwk-sig-pub.json",
				dir + "/jwk-sig-priv.json",
			},
		},
	}

	service, err := NewWebKeysService(s.logger, config)
	a.Nil(err)

	keys, err := service.Keys()
	a.Nil(err)

	a.Len(keys.Keys, 2)
}

func Test_FileKeys_NotFound(t *testing.T) {
	a, s := Setup(t)

	// testing...
	config := WebKeysConfig{
		File: &FileConfig{
			Paths: []string{
				os.TempDir() + "/jwk-sig-pub.json",
				os.TempDir() + "/jwk-sig-priv.json",
			},
		},
	}

	_, err := NewWebKeysService(s.logger, config)
	a.NotNil(err)
}

type Scope struct {
	logger log.Logger
}

func Setup(t *testing.T) (*assert.Assertions, Scope) {
	return assert.New(t), Scope{
		logger: log.NewNopLogger(),
	}
}