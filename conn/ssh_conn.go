package conn

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pythonandchips/azad/logger"
	"golang.org/x/crypto/ssh"
)

// Config for connection
type Config struct {
	KeyFilePath string
	User        string
}

var dial = func(ip string, config *ssh.ClientConfig) (sshClient, error) {
	conn, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return sshClientWrapper{}, fmt.Errorf("Unable to dial %s: %s", ip, err)
	}
	return sshClientWrapper{conn}, nil
}

// ConnectTo creates new connection to specified address
func (sshConn *SSHConn) ConnectTo(ip string) error {
	key, err := ioutil.ReadFile(sshConn.config.KeyFilePath)
	if err != nil {
		return fmt.Errorf("Unable to read private key %s", sshConn.config.KeyFilePath)
	}
	logger.Info("Using key %s", sshConn.config.KeyFilePath)

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: sshConn.config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := dial(ip, config)
	if err != nil {
		return err
	}
	sshConn.client = client
	return nil
}

var now = func() int64 {
	return time.Now().UnixNano()
}

// Run a command against the sshConn host
//
// Returns the result of stdout and stderr in the command response
// Returns an error if an ssh session cannot be created or
// if the command has a non-zero exit code
func (sshConn *SSHConn) Run(command Command) (Response, error) {
	ref := now()
	writeCommandFile := fmt.Sprintf("echo '%s' > /tmp/azad.%d && chmod +x /tmp/azad.%d", command.generateFile(), ref, ref)
	if _, err := sshConn.runOnClient(writeCommandFile); err != nil {
		return CommandResponse{}, err
	}
	var runCommand string
	if command.User != "" {
		runCommand = fmt.Sprintf("sudo su - %s -c '/tmp/azad.%d'", command.User, ref)
	} else {
		runCommand = fmt.Sprintf("/tmp/azad.%d", ref)
	}
	commandResposne, err := sshConn.runOnClient(runCommand)
	if err != nil {
		return commandResposne, err
	}
	cleanUpCommand := fmt.Sprintf("rm /tmp/azad.%d", ref)
	_, err = sshConn.runOnClient(cleanUpCommand)
	return commandResposne, err
}

func (sshConn SSHConn) runOnClient(command string) (Response, error) {
	commandResposne := CommandResponse{
		stdout: bytes.NewBuffer([]byte{}),
		stderr: bytes.NewBuffer([]byte{}),
	}
	session, err := sshConn.client.NewSession()
	if err != nil {
		return commandResposne, err
	}
	session.setStdout(commandResposne.stdout)
	session.setStderr(commandResposne.stderr)
	defer session.Close()
	err = session.Run(command)
	return commandResposne, err
}

// Close connection to client
func (sshConn *SSHConn) Close() {
	sshConn.client.Close()
}

// SSHConn manage connections
type SSHConn struct {
	client sshClient
	config Config
}
