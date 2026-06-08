// Package lib provides a clean, reusable API for interacting with the I2P control interface.
package lib

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-i2p/go-i2pcontrol"
)

var Usage = `i2p-control
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

// Config holds configuration for the I2P control client.
type Config struct {
	Host     string
	Port     string
	Path     string
	Password string
	Verbose  bool
}

// Client represents a connection to the I2P control interface.
type Client struct {
	config Config
}

// CommandOptions holds options for command execution.
type CommandOptions struct {
	Block bool
	Args  []string
}

// CommandResult represents the result of a command execution.
type CommandResult struct {
	Success      bool
	Output       string
	Error        error
	IsShutdown   bool
	WaitShutdown bool
}

// New creates a new I2P control client with the given configuration.
func New(config Config) *Client {
	return &Client{
		config: config,
	}
}

// Connect initializes the connection to the I2P control interface.
func (c *Client) Connect() error {
	i2pcontrol.Initialize(c.config.Host, c.config.Port, c.config.Path)
	_, err := i2pcontrol.Authenticate(c.config.Password)
	return err
}

// logVerbose logs a message if verbose mode is enabled.
func (c *Client) logVerbose(msg string) {
	if c.config.Verbose {
		fmt.Println(msg)
	}
}

// Execute runs a command and returns the result.
func (c *Client) Execute(command string, opts CommandOptions) CommandResult {
	result := CommandResult{
		Success:      true,
		IsShutdown:   false,
		WaitShutdown: false,
	}

	c.logVerbose(command)

	switch command {
	case "echo":
		message, err := i2pcontrol.Echo(strings.Join(opts.Args, " "))
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = message

	case "ident":
		message, err := i2pcontrol.RouterHash()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = message

	case "restart":
		message, err := i2pcontrol.Restart()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = message
		result.IsShutdown = true
		result.WaitShutdown = opts.Block

	case "graceful-restart":
		message, err := i2pcontrol.RestartGraceful()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = message
		result.IsShutdown = true
		result.WaitShutdown = opts.Block

	case "shutdown":
		message, err := i2pcontrol.Shutdown()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = message
		result.IsShutdown = true
		result.WaitShutdown = opts.Block

	case "graceful-shutdown":
		message, err := i2pcontrol.ShutdownGraceful()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = message
		result.IsShutdown = true
		result.WaitShutdown = opts.Block

	case "update":
		needsUpdate, err := i2pcontrol.FindUpdates()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		if needsUpdate {
			result.Output = "You need an update"
			message, err := i2pcontrol.Update()
			if err != nil {
				result.Success = false
				result.Error = err
				break
			}
			result.Output = message
		} else {
			result.Output = "You don't need an update"
		}

	case "find-update":
		needsUpdate, err := i2pcontrol.FindUpdates()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		if needsUpdate {
			result.Output = "You need an update"
		} else {
			result.Output = "You don't need an update"
		}

	case "stat":
		message, err := i2pcontrol.Status()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = fmt.Sprintf("%v", message)

	case "netstat":
		message, err := i2pcontrol.NetStatus()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = fmt.Sprintf("%v", message)

	case "reseedstat":
		isReseeding, err := i2pcontrol.Reseeding()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		if isReseeding {
			result.Output = "Router is reseeding"
		} else {
			result.Output = "Router is not reseeding"
		}

	case "tunstat":
		message, err := i2pcontrol.ParticipatingTunnels()
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = fmt.Sprintf("%v", message)

	case "ratestat":
		if len(opts.Args) < 2 {
			result.Success = false
			result.Error = fmt.Errorf("ratestat requires 2 arguments")
			break
		}
		interval, err := strconv.Atoi(opts.Args[1])
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		message, err := i2pcontrol.RateStat(opts.Args[0], interval)
		if err != nil {
			result.Success = false
			result.Error = err
			break
		}
		result.Output = fmt.Sprintf("%v", message)

	default:
		result.Success = false
		result.Error = fmt.Errorf("unknown command: %s", command)
	}

	return result
}

// WaitForShutdown blocks until the router is shut down.
// Returns the number of tunnels that were participating at shutdown.
func (c *Client) WaitForShutdown() (int, error) {
	lastParticipatingTunnels, err := i2pcontrol.ParticipatingTunnels()
	if err != nil {
		return 0, err
	}

	baseminutes := time.Duration(time.Minute * 11)
	if lastParticipatingTunnels != 0 {
		fmt.Printf("Waiting for expiration of %d participating tunnels in %v\n", lastParticipatingTunnels, baseminutes)
	}

	oldtime := time.Now()
	for {
		participatingTunnels, err := i2pcontrol.ParticipatingTunnels()
		if err != nil {
			return 0, err
		}

		if participatingTunnels != lastParticipatingTunnels {
			minutes := oldtime.Sub(time.Now())
			fmt.Printf("Waiting for expiration of %d participating tunnels in %v\n", participatingTunnels, baseminutes+minutes)
			lastParticipatingTunnels = participatingTunnels
		}

		time.Sleep(time.Duration(time.Second * 1))
		if participatingTunnels < 1 {
			break
		}
	}

	return lastParticipatingTunnels, nil
}
