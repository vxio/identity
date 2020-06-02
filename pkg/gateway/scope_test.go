package gateway_test

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	. "github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/webkeys"
	"github.com/stretchr/testify/assert"
)

type Scope struct {
	a          *assert.Assertions
	time       stime.StaticTimeService
	keys       webkeys.GenerateJwksService
	mw         *Middleware
	identityID IdentityID
	tenantID   TenantID
}

func NewScope(t *testing.T) Scope {
	l := logging.NewDefaultLogger()
	a := assert.New(t)

	s := NewRandomSession()

	stime := stime.NewStaticTimeService()

	keys, err := webkeys.NewGenerateJwksService()
	a.Nil(err)

	gatewayMiddleware, err := NewMiddleware(l, stime, keys)
	a.Nil(err)

	return Scope{
		a:          a,
		time:       stime,
		keys:       *keys,
		mw:         gatewayMiddleware,
		identityID: s.CallerID,
		tenantID:   s.TenantID,
	}
}

func (s *Scope) SignedString(sessionJwt SessionClaims) string {
	signingMethod := jwt.GetSigningMethod(s.keys.Private.Algorithm)

	token := jwt.NewWithClaims(signingMethod, sessionJwt)
	token.Header["kid"] = s.keys.Private.KeyID

	tokenString, err := token.SignedString(s.keys.Private.Key)
	s.a.Nil(err)

	return tokenString
}

func (s *Scope) NewSessionJwt() SessionClaims {
	return SessionClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: s.time.Now().Add(time.Hour).Unix(),
			NotBefore: s.time.Now().Add(time.Minute * -1).Unix(),
			IssuedAt:  s.time.Now().Unix(),
			Id:        uuid.New().String(),
			Subject:   s.identityID.String(),

			Audience: "moov",
			Issuer:   "moov",
		},
		CallerID: uuid.UUID(s.identityID),
		TenantID: uuid.UUID(s.tenantID),
	}
}
