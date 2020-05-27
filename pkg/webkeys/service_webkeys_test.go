package webkeys_test

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	. "github.com/moov-io/identity/pkg/webkeys"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateKeys(t *testing.T) {
	a, s := Setup(t)

	config := WebKeysConfig{}
	service, err := NewWebKeysService(s.logger, config)
	a.Nil(err)

	keys, err := service.FetchJwks()
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
			URL: server.URL + controller.WellKnownJwksPath(),
		},
	}

	service, err := NewWebKeysService(s.logger, config)
	a.Nil(err)

	keys, err := service.FetchJwks()
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
			URL: server.URL + "/.well-known/jwks.json",
		},
	}

	service, err := NewWebKeysService(s.logger, config)
	a.Nil(err)

	_, err = service.FetchJwks()
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
			Path: dir + "/jwk-sig-pub.json",
		},
	}

	service, err := NewWebKeysService(s.logger, config)
	a.Nil(err)

	keys, err := service.FetchJwks()
	a.Nil(err)

	a.Len(keys.Keys, 1)
}

func Test_FileKeys_NotFound(t *testing.T) {
	a, s := Setup(t)

	// testing...
	config := WebKeysConfig{
		File: &FileConfig{
			Path: os.TempDir() + "/jwk-sig-pub.json",
		},
	}

	service, err := NewWebKeysService(s.logger, config)
	a.Nil(err)

	_, err = service.FetchJwks()
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
