package conn

import (
	"fmt"
	"testing"
)

func TestSSHConnection(t *testing.T) {
	sshConn := &SSHConn{}
	defer sshConn.Close()
	err := sshConn.ConnectTo([]string{"52.91.156.19"})
	if err != nil {
		fmt.Println(err)
	}
	command := Command{
		Command: []string{
			"apt-get update",
			"apt-get install -yq mysql-server",
		},
		Interpreter: "bash",
		User:        "root",
	}
	err = sshConn.Run(command)
	if err != nil {
		fmt.Println(err)
	}
}
