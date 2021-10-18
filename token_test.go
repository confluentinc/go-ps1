package ps1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	token := &Token{Name: 'A'}
	require.Equal(t, "%A", token.String())
}
