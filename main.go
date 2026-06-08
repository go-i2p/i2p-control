package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-i2p/i2p-control/lib"
)

var usage = `i2p-control
===========

Terminal interface to monitor and manage I2P router service. Basically, an
terminal i2pcontrol client.

        -host default:"127.0.0.1"
        -port default:"7657"
        -path default:"jsonrpc"
        -password default:"itoopie"
        -method default:"echo"
        -block default:false
        -verbose default:false

Installation with go get

        go get -u github.com/go-i2p/i2p-control

The methods that have been implemented are

        echo              : i2pcontrol:Echo
        stat              : i2pcontrol:RouterInfo:i2p.router.status
        netstat           : i2pcontrol:RouterInfo:i2p.router.net.router.status
        tunstat           : i2pcontrol:RouterInfo:i2p.router.net.tunnels.participating
        restart           : i2pcontrol:Restart
        graceful-restart  : i2pcontrol:RestartGraceful
        shutdown          : i2pcontrol:Shutdown
        graceful-shutdown : i2pcontrol:ShutdownGraceful
        update            : i2pcontrol:Update
        find-update       : i2pcontrol:FindUpdate
		ident             : i2pcontrol:RouterInfo:hash(go-i2p only)

So, for instance, to initiate a graceful shutdown and block until the router is
shut down, use the command:

        i2p-control -block -method=graceful-shutdown

`

var (
	host     = flag.String("host", "localhost", "Host of the i2pcontrol interface")
	port     = flag.String("port", "7657", "Port of the i2pcontrol interface")
	path     = flag.String("path", "jsonrpc", "Path to the i2pcontrol interface")
	password = flag.String("password", "itoopie", "Password for the i2pcontrol interface")
	command  = flag.String("method", "echo", "Method call to invoke")
	shelp    = flag.Bool("h", false, "Show the help message")
	sverbose = flag.Bool("v", false, "Verbosely update participating tunnel count while running.")
	lverbose = flag.Bool("verbose", false, "Verbosely update participating tunnel count while running.")
	lhelp    = flag.Bool("help", false, "Show the help message")
	block    = flag.Bool("block", false, "Block the terminal until the router is completely shut down")
)

func main() {
	flag.Parse()
	if *shelp || *lhelp {
		fmt.Printf(usage)
		return
	}

	args := flag.Args()
	if len(args) < 2 {
		args = append(args, "bw.sendBps")
		args = append(args, "300000")
	}

	// Create client configuration
	config := lib.Config{
		Host:     *host,
		Port:     *port,
		Path:     *path,
		Password: *password,
		Verbose:  *sverbose || *lverbose,
	}

	// Create and connect client
	client := lib.New(config)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	// Execute command
	result := client.Execute(*command, lib.CommandOptions{
		Block: *block,
		Args:  args,
	})

	if !result.Success {
		log.Fatal(result.Error)
	}

	if result.Output != "" {
		fmt.Println(result.Output)
	}

	// Handle shutdown if needed
	if result.IsShutdown && result.WaitShutdown {
		if _, err := client.WaitForShutdown(); err != nil {
			log.Fatal(err)
		}
	}
}
