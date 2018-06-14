package conn

type Conn interface {
	ConnectTo(string) error
	Run(Command) error
	Close()
}

var NewConn = newConn

func newConn() Conn {
	return &SSHConn{}
}

func newFakeConn() Conn {
	return &FakeSSHConn{}
}

func Simulate() {
	NewConn = newFakeConn
}
