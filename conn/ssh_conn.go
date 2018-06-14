package conn

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

// SSHConn manage connections
type SSHConn struct {
	client *ssh.Client
}

func (sshConn *SSHConn) ConnectTo(ip string) error {
	home, _ := homedir.Dir()
	keyPath := filepath.Join(home, ".ssh", "id_rsa")
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("Unable to read private key %s", keyPath)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("Unable to parse private key: %v\n", err)
	}

	config := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return fmt.Errorf("Unable to dial %s: %s", ip, err)
	}
	sshConn.client = conn
	return nil
}

func (sshConn *SSHConn) Run(command Command) error {
	ref := time.Now().UnixNano()
	commands := []string{
		fmt.Sprintf("echo '%s' > /tmp/azad.%d && chmod +x /tmp/azad.%d",
			command.generateFile(), ref, ref),
	}
	if command.User != "" {
		commands = append(commands, fmt.Sprintf("sudo su - %s -c '/tmp/azad.%d'",
			command.User, ref))
	} else {
		commands = append(commands, fmt.Sprintf("/tmp/azad.%d", ref))
	}
	commands = append(commands, fmt.Sprintf("rm /tmp/azad.%d", ref))
	return sshConn.runOnClient(commands)
}

func (conn SSHConn) runOnClient(commands []string) error {
	for _, c := range commands {
		session, err := conn.client.NewSession()
		defer session.Close()
		if err != nil {
			return err
		}
		if err := session.Run(c); err != nil {
			return err
		}
	}
	return nil
}

func (sshConn *SSHConn) Close() {
	sshConn.client.Close()
}
