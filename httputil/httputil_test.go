package httputil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExcludeSensitiveRequestBody(t *testing.T) {
	sourceText := `{"log_message":"user: Wrong username or password","method":"POST","request_body":"{\n\t\"username\": \"teta.kibites@gmail.com\",\n\t\"password\": \"k1tab1saa\"\n}","response_body":"{\"api_code\":101200,\"errors\":[{\"details\":{\"id\":\"Username atau password salah\",\"en\":\"Wrong username or password\"}}]}","stack":[{"message":"user: Wrong username or password"}]`
	match := passRemover.MatchString(sourceText)
	assert.True(t, match)

	ExcludeSensitiveRequestBody(&sourceText)
	t.Log(sourceText)

	match = passRemover.MatchString(sourceText)
	assert.False(t, match)
}