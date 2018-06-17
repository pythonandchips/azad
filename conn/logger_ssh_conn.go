package conn

import "github.com/pythonandchips/azad/logger"

// LoggerSSHConn log all ssh command to stdout instead of running agains a server
//
// Used for development and simulating configuration runs
type LoggerSSHConn struct {
	ConnectedTo string
	Commands    []Command
	closed      bool
}

// ConnectTo track the host name that would be used to connect to server
func (loggerSSHConn *LoggerSSHConn) ConnectTo(hostName string) error {
	loggerSSHConn.ConnectedTo = hostName
	return nil
}

// Run output the command to logger.Debug
func (loggerSSHConn *LoggerSSHConn) Run(command Command) (Response, error) {
	loggerSSHConn.Commands = append(loggerSSHConn.Commands, command)
	logger.Debug("Running on %s", loggerSSHConn.ConnectedTo)
	logger.Debug(command.generateFile())
	return CommandResponse{}, nil
}

// Close mark the connection as closed
func (loggerSSHConn *LoggerSSHConn) Close() {
	logger.Debug("Closing connection to %s", loggerSSHConn.ConnectedTo)
	loggerSSHConn.closed = true
}
