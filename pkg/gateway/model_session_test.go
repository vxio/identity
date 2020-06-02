package gateway

import (
	"testing"

	"github.com/moov-io/identity/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func Test_SessionLogAttributes(t *testing.T) {
	session := NewRandomSession()

	buffer, logger := logging.NewBufferLogger()

	logger.With(&session).Info().Log("test")

	output := buffer.String()
	assert.Contains(t, output, "identity_id="+session.CallerID.String())
	assert.Contains(t, output, "tenant_id="+session.TenantID.String())
}
