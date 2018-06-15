package conn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHConnRun(t *testing.T) {
	nowValue := int64(123456)
	now = func() int64 {
		return nowValue
	}
	stdout := "data from stdout"
	stderr := "data from stderr"
	fakeSSHSession := &FakeSSHSession{
		stdout: stdout,
		stderr: stderr,
	}
	fakeSSHClient := FakeSSHClient{
		session: fakeSSHSession,
	}
	sshConn := SSHConn{client: fakeSSHClient}

	t.Run("with no command user", func(t *testing.T) {
		fakeSSHSession.Clear()
		command := Command{
			Command: []string{
				"whoami",
			},
		}
		commandResponse, err := sshConn.Run(command)

		if err != nil {
			t.Fatalf("Unexpected error %s", err)
		}
		t.Run("sends the command file to the server", func(t *testing.T) {
			expectedCommand := fmt.Sprintf(
				"echo '%s' > /tmp/azad.%d && chmod +x /tmp/azad.%d",
				command.generateFile(),
				nowValue,
				nowValue,
			)
			if len(fakeSSHSession.commands) == 0 {
				t.Fatalf("expected at least one command but got none")
			}
			assert.Equal(t, fakeSSHSession.commands[0], expectedCommand)
		})
		t.Run("sends runs the command on the server", func(t *testing.T) {
			expectedCommand := fmt.Sprintf("/tmp/azad.%d", nowValue)
			if len(fakeSSHSession.commands) == 0 {
				t.Fatalf("expected at least one command but got none")
			}
			assert.Equal(t, fakeSSHSession.commands[1], expectedCommand)
		})
		t.Run("removes the command from server", func(t *testing.T) {
			expectedCommand := fmt.Sprintf("rm /tmp/azad.%d", nowValue)
			if len(fakeSSHSession.commands) == 0 {
				t.Fatalf("expected at least one command but got none")
			}
			assert.Equal(t, fakeSSHSession.commands[2], expectedCommand)
		})
		t.Run("returns the stdout from the command", func(t *testing.T) {
			assert.Equal(t, commandResponse.Stdout(), stdout)
		})
		t.Run("returns the stderr from the command", func(t *testing.T) {
			assert.Equal(t, commandResponse.Stderr(), stderr)
		})
	})
	t.Run("with user set in command", func(t *testing.T) {
		fakeSSHSession.Clear()
		command := Command{
			User: "root",
			Command: []string{
				"whoami",
			},
		}
		_, err := sshConn.Run(command)

		if err != nil {
			t.Fatalf("Unexpected error %s", err)
		}
		t.Run("runs the task with the specified user", func(t *testing.T) {
			expectedCommand := fmt.Sprintf("sudo su - %s -c '/tmp/azad.%d'", command.User, nowValue)
			if len(fakeSSHSession.commands) == 0 {
				t.Fatalf("expected at least one command but got none")
			}
			assert.Equal(t, fakeSSHSession.commands[1], expectedCommand)
		})
	})
}
