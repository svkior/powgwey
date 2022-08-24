package main

import (
	"flag"
	"net"
	"os"

	"puretcp/pow"

	"gopkg.in/dailymuse/gzap.v1"
)

func PrintAllEnvVariables() {
	allEnvs := os.Environ()

	for _, value := range allEnvs {
		gzap.Logger.Info("Client",
			gzap.String("env", value),
		)
	}
}

func main() {
	//PrintAllEnvVariables()

	netPath := flag.String("server", "server:8000", "connect to server string")

	if err := gzap.InitLogger(); err != nil {
		panic(err)
	}

	solv, err := pow.NewSolverPlugin()
	if err != nil {
		gzap.Logger.Fatal("unable to create PowSolver", gzap.Error(err))
	}

	gzap.Logger.Info("Connecting to server", gzap.String("net_path", *netPath))

	conn, err := net.Dial("tcp", *netPath)
	if err != nil {
		gzap.Logger.Fatal("unable to connect to server ",
			gzap.Error(err),
			gzap.String("net_path", *netPath),
		)
	}

	err = solv.Solve(conn)
	if err != nil {
		gzap.Logger.Fatal("unable to solve puzzle ",
			gzap.Error(err),
			gzap.String("net_path", *netPath),
		)
	}

	reply := make([]byte, 1024)
	bytes_readed, err := conn.Read(reply)
	if err != nil {
		gzap.Logger.Fatal("error read quote from server",
			gzap.Error(err),
			gzap.String("net_path", *netPath),
		)
	}

	quote := string(reply[0:bytes_readed])

	gzap.Logger.Info("got quote from server",
		gzap.String("quote", quote),
	)
}
