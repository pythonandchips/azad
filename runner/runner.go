package runner

import "github.com/pythonandchips/azad/conn"

type runners []runner

func (runners runners) Close() {
	for _, runner := range runners {
		runner.Conn.Close()
	}
}

type runner struct {
	Address   string
	Conn      conn.Conn
	Variables map[string]string
}
