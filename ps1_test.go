package ps1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	New(exampleCLIName, exampleTokens)
}

func TestNew_PanicNoCLIName(t *testing.T) {
	require.Panics(t, func() {
		New("", exampleTokens)
	})
}

func TestNew_PanicNoTokens(t *testing.T) {
	for _, tokens := range [][]Token{nil, {}} {
		require.Panics(t, func() {
			New(exampleCLIName, tokens)
		})
	}
}
