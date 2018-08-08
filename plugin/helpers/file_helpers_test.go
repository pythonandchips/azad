package helpers

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	"github.com/stretchr/testify/assert"
)

func TestChecksum(t *testing.T) {
	data := []byte("CHECKSUM THIS")
	path := "$HOME/file.conf"

	commands := Checksum(data, path)

	expect.EqualFatal(t, len(commands), 2)

	assert.Equal(t, commands[0], `echo "9081ec21087deb3803c02119e8be3372c2fe5c0c $HOME/file.conf" | sha1sum -c -`)
	assert.Equal(t, commands[1], `if [ $? = 0 ]; then exit 40; fi`)
}

func TestWriteEncodedFile(t *testing.T) {
	data := []byte("CHECKSUM THIS")
	path := "$HOME/file.conf"

	commands := WriteEncodedFile(data, path)

	expect.EqualFatal(t, len(commands), 2)
	assert.Equal(t, commands[0], "filebase64encoded=Q0hFQ0tTVU0gVEhJUw==")
	assert.Equal(t, commands[1], "echo $filebase64encoded | base64 -d > $HOME/file.conf")
}
