package conn

import "fmt"

type StdOutConn struct {
	hostName string
}

func (stdOutConn *StdOutConn) ConnectTo(hostName string) error {
	stdOutConn.hostName = hostName
	fmt.Printf("Connected to %s\n", hostName)
	return nil
}

func (stdOutConn *StdOutConn) Run(command Command) error {
	fmt.Printf("================= BEGIN %s =================\n\n", stdOutConn.hostName)
	fmt.Println(command.generateFile())
	fmt.Printf("================= END %s =================\n\n", stdOutConn.hostName)
	return nil
}

func (stdOutConn *StdOutConn) Close() {
	fmt.Printf("Closed Connected to %s\n", stdOutConn.hostName)
}
