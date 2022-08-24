package puretcp

import (
	"net"

	"puretcp/pow"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/puretcp", new(TCP))
}

type solver interface {
	Solve(conn net.Conn) error
}

type TCP struct {
	solv solver
}

func (tcp *TCP) Connect(addr string) (net.Conn, error) {

	solv, err := pow.NewSolverPlugin()
	if err != nil {
		return nil, err
	}

	tcp.solv = solv

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (tcp *TCP) GetQuote(conn net.Conn) string {
	err := tcp.solv.Solve(conn)
	if err != nil {
		return ""
	}

	quote := tcp.Read(conn)
	//log.Printf("QUOTE IS: %s", quote)
	return quote
}

func (tcp *TCP) Write(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (tcp *TCP) WriteLn(conn net.Conn, data []byte) error {
	return tcp.Write(conn, append(data, []byte("\n")...))
}

func (tcp *TCP) Read(conn net.Conn) string {
	reply := make([]byte, 1024)

	bytes_readed, err := conn.Read(reply)
	if err != nil {
		return ""
	}
	return string(reply[0:bytes_readed])
}
