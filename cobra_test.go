package ps1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const exampleCLIName = "confluent"

var exampleTokens = []Token{
	{
		Name: 'A',
		Desc: "Description of token A.",
		Func: func() string { return "Value A" },
	},
	{
		Name: 'B',
		Desc: "Description of token B.",
		Func: func() string { return "Value B" },
	},
}

func TestShort(t *testing.T) {
	p := &ps1{cliName: exampleCLIName}
	require.Equal(t, "Add confluent context to your terminal prompt.", p.short())
}

func TestLong(t *testing.T) {
	expected := `Use this command to add ` + "`" + `confluent` + "`" + ` information to your terminal prompt.

Bash:

$ export PS1='$(confluent prompt) '$PS1

ZSH:

$ setopt prompt_subst
$ export PS1='$(confluent prompt) '$PS1

You can customize the prompt by calling passing the ` + "`" + `--format` + "`" + ` flag, for example ` + "`" + `-f '(confluent|%C)'` + "`" + `.
To make this permanent, you must add the above lines to your Bash or ZSH profile.

Formatting Tokens
~~~~~~~~~~~~~~~~~

This command comes with a number of formatting tokens. What follows is a list of all tokens:

* %A - Description of token A.
* %B - Description of token B.

Style
~~~~~

The style of the text can be changed with a combination of functions, colors, and attributes.

Functions:
* fgcolor - Change the foreground color.
* bgcolor - Change the background color.
* attr    - Change a text attribute.

Colors:
* black
* blue
* cyan
* green
* magenta
* red
* white
* yellow

Text Attributes:
* bold
* invert
* italicize
* underline

Examples
~~~~~~~~

* {{fgcolor "blue" "this text is blue"}}
* {{bgcolor "blue" "this text has a blue background"}}
* {{attr "bold" "this text is bold"}}

Use a vertical bar to separate further attributes:
* {{fgcolor "red" "this text is red, has a blue background, and is bold"|bgcolor "blue"|attr "bold"}}

We can use tokens and colors in the same format string:
* ({{fgcolor "blue" "confluent"}}|{{fgcolor "red" "%C"}})`

	p := New(exampleCLIName, exampleTokens)
	require.Equal(t, expected, p.long())
}

func TestPrompt(t *testing.T) {
	p := new(ps1)
	out := p.prompt("format", 1000)
	require.Equal(t, "format", out)
}

func TestPrompt_ErrUndefinedFunction(t *testing.T) {
	p := &ps1{cliName: exampleCLIName}
	out := p.prompt(`{{func}}`, 1000)
	require.Equal(t, `(confluent|function "func" not defined)`, out)
}

func TestPrompt_ErrColorNotFound(t *testing.T) {
	p := &ps1{cliName: exampleCLIName}
	out := p.prompt(`{{fgcolor "gray"}}`, 1000)
	require.Equal(t, `(confluent|fgcolor "gray" not found)`, out)
}

func TestPrompt_Timeout(t *testing.T) {
	p := &ps1{cliName: exampleCLIName}
	out := p.prompt("", 0)
	require.Equal(t, "(confluent)", out)
}

func TestPrompt_TimeoutOnToken(t *testing.T) {
	p := &ps1{
		cliName: exampleCLIName,
		tokens: []Token{
			{
				Name: 'A',
				Func: func() string {
					time.Sleep(time.Second)
					return ""
				},
			},
		},
	}
	out := p.prompt("%A", 0)
	require.Equal(t, "(confluent)", out)
}

func TestWrite_Basic(t *testing.T) {
	p := new(ps1)
	out, err := p.write("format")
	require.NoError(t, err)
	require.Equal(t, "format", out)
}

func TestWrite_Color(t *testing.T) {
	p := new(ps1)
	out, err := p.write(`{{fgcolor "black" "format"}}`)
	require.NoError(t, err)
	require.Equal(t, "format", out)
}

func TestWrite_ErrColorNotFound(t *testing.T) {
	p := new(ps1)
	_, err := p.write(`{{fgcolor "gray"}}`)
	require.Error(t, err)
}

func TestWrite_ColorAndAttributes(t *testing.T) {
	p := new(ps1)
	out, err := p.write(`{{fgcolor "black" "format" | attr "bold"}}`)
	require.NoError(t, err)
	require.Equal(t, "format", out)
}

func TestWrite_Token(t *testing.T) {
	p := &ps1{tokens: exampleTokens}
	out, err := p.write("%A")
	require.NoError(t, err)
	require.Equal(t, "Value A", out)
}

func TestWrite_ColorAndToken(t *testing.T) {
	p := &ps1{tokens: exampleTokens}
	out, err := p.write(`{{fgcolor "black" "%A"}}`)
	require.NoError(t, err)
	require.Equal(t, "Value A", out)
}
