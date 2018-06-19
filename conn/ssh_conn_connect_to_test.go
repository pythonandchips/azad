package conn

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pythonandchips/azad/logger"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestSSHConnConnectTo(t *testing.T) {
	logger.StubLogger()
	client := &FakeSSHClient{}
	connectIP := "10.0.0.1"
	var connectedToIP string
	var sshConfig *ssh.ClientConfig
	dial = func(ip string, config *ssh.ClientConfig) (sshClient, error) {
		connectedToIP = ip
		sshConfig = config
		return client, nil
	}

	t.Run("completes connection to the server", func(t *testing.T) {
		fakeKey, _ := filepath.Abs("./fixtures/fake_key")
		sshUser := "ssh_user"
		config := Config{KeyFilePath: fakeKey, User: sshUser}
		conn := SSHConn{config: config}
		err := conn.ConnectTo(connectIP)
		if err != nil {
			t.Fatalf("expected no error but got %s", err)
		}

		assert.Equal(t, conn.client, client)
		assert.Equal(t, connectIP, connectedToIP)
	})
	t.Run("specifies the configuration of the client", func(t *testing.T) {
		fakeKey, _ := filepath.Abs("./fixtures/fake_key")
		sshUser := "ssh_user"
		config := Config{KeyFilePath: fakeKey, User: sshUser}
		conn := SSHConn{config: config}
		err := conn.ConnectTo(connectIP)
		if err != nil {
			t.Fatalf("expected no error but got %s", err)
		}

		assert.Equal(t, sshConfig.User, sshUser)
		assert.Equal(t, len(sshConfig.Auth), 1)
		assert.NotNil(t, sshConfig.HostKeyCallback)
	})
	t.Run("returns an error if the key does not exist", func(t *testing.T) {
		fakeKey, _ := filepath.Abs("./fixtures/key_does_not_exist")
		sshUser := "ssh_user"
		config := Config{KeyFilePath: fakeKey, User: sshUser}
		conn := SSHConn{config: config}
		err := conn.ConnectTo(connectIP)
		if err.Error() == fmt.Sprintf("unable to read private key %s", fakeKey) {
			t.Fatalf("expected no error but got %s", err)
		}
	})
	t.Run("returns an error if the key cannot be read", func(t *testing.T) {
		fakeKey, _ := filepath.Abs("./fixtures/invalid_key")
		sshUser := "ssh_user"
		config := Config{KeyFilePath: fakeKey, User: sshUser}
		conn := SSHConn{config: config}
		err := conn.ConnectTo(connectIP)
		if err.Error() == fmt.Sprintf("unable to parse private key %s", fakeKey) {
			t.Fatalf("expected no error but got %s", err)
		}
	})
	t.Run("returns an error if it cannont connect to server", func(t *testing.T) {
		connectionError := fmt.Errorf("Cannot connect to server")
		dial = func(string, *ssh.ClientConfig) (sshClient, error) {
			return nil, connectionError
		}
		fakeKey, _ := filepath.Abs("./fixtures/fake_key")
		sshUser := "ssh_user"
		config := Config{KeyFilePath: fakeKey, User: sshUser}
		conn := SSHConn{config: config}
		err := conn.ConnectTo(connectIP)
		if err != connectionError {
			t.Fatalf("expected connection error but got %s", err)
		}
	})
}
