package ps1

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
)

func TestPrintWithAttr(t *testing.T) {
	attrs := map[string]color.Attribute{
		"blue": color.FgBlue,
	}
	out, err := printWithAttr("", attrs, "blue", "text")
	require.NoError(t, err)
	require.Equal(t, "text", out)
}

func TestPrintWithAttr_ErrNotFound(t *testing.T) {
	_, err := printWithAttr("attrType", map[string]color.Attribute{}, "attr")
	require.Error(t, err)
}
