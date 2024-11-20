// Copyright 2024 Colorado School of Mines CSCI 370 FA24 NREL 2 Group

package vwifi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	log "github.com/sandia-minimega/minimega/v2/pkg/minilog"
)

// ServerOptions are the vwifi-server command line options
type ServerOptions struct {
	// Whether or not to enable packet loss
	PacketLoss bool

	// Whether or not to use the port number in computing the client hash ID
	UsePortInHash bool

	// The listening port for the VSOCK socket
	VsockPort uint16

	// The listening port for the TCP socket
	TcpPort uint16

	// The listening port for the network spy socket
	SpyPort uint16

	// The listening port for the control socket
	ControlPort uint16
}

// Server is a vwifi-server instance
type Server struct {
	// The server binary path
	binary string

	// The server options
	options ServerOptions

	// The server command
	cmd *exec.Cmd
}

// NewServer creates a new vwifi-server instance
func NewServer(binary string, options ServerOptions) *Server {
	return &Server{
		binary:  binary,
		options: options,
	}
}

// maxServerStatusChecks is the number of times to check if the server is ready
const maxServerStatusChecks = 10

// serverStatusWait is the time to wait between server status checks
var serverStatusWait = 1 * time.Second

// ServerStatus is the status of the server as reported by the server itself
type ServerStatus struct {
	// The listening port for the VSOCK socket
	VsockPort uint16

	// The listening port for the TCP socket
	TcpPort uint16

	// The listening port for the network spy socket
	SpyPort uint16

	// The listening port of the server's control socket
	ControlPort uint16

	// No clue what this is, but it's in the output
	SizeOfDisconnected uint32

	// Whether or not the server is simulating packet loss
	PacketLoss bool

	// The simulation's physical scale
	Scale float64
}

// serverStatusPattern is the pattern to match the status output as reported by the server
var serverStatusPattern = regexp.MustCompile(`^\s*CLIENT VHOST : Listener on port : (\d+)\nCLIENT TCP : Listener on port : (\d+)\nSPY : Listener on port : (\d+)\nCTRL : Listener on port : (\d+)\nSize of disconnected : (\d+)\nPacket loss : (enable|disable)\nScale : ([+-]?\d+(?:\.\d+)?)\s*$`)

// errServerStatusNotMatched is the error when the server status does not match the pattern
var errServerStatusNotMatched = errors.New("server status not matched")

// parseStatus parses the server status
func (server *Server) parseStatus(raw string) (*ServerStatus, error) {
	// Parse the status
	matches := serverStatusPattern.FindStringSubmatch(raw)

	if matches == nil {
		return nil, errServerStatusNotMatched
	}

	vsockPort, err := strconv.ParseUint(matches[1], 10, 16)

	if err != nil {
		return nil, err
	}

	tcpPort, err := strconv.ParseUint(matches[2], 10, 16)

	if err != nil {
		return nil, err
	}

	spyPort, err := strconv.ParseUint(matches[3], 10, 16)

	if err != nil {
		return nil, err
	}

	controlPort, err := strconv.ParseUint(matches[4], 10, 16)

	if err != nil {
		return nil, err
	}

	sizeOfDisconnected, err := strconv.ParseUint(matches[5], 10, 32)

	if err != nil {
		return nil, err
	}

	scale, err := strconv.ParseFloat(matches[7], 64)

	if err != nil {
		return nil, err
	}

	status := ServerStatus{
		VsockPort:          uint16(vsockPort),
		TcpPort:            uint16(tcpPort),
		SpyPort:            uint16(spyPort),
		ControlPort:        uint16(controlPort),
		SizeOfDisconnected: uint32(sizeOfDisconnected),
		PacketLoss:         matches[6] == "enable",
		Scale:              scale,
	}

	return &status, nil
}

// Start the server
func (server *Server) Start() (*ServerStatus, error) {
	// Generate the arguments
	args := []string{}

	if server.options.PacketLoss {
		args = append(args, "--lost-packets")
	}

	if server.options.UsePortInHash {
		args = append(args, "--use-port-in-hash")
	}

	if server.options.VsockPort != 0 {
		args = append(args, "--port-vhost", fmt.Sprintf("%d", server.options.VsockPort))
	}

	if server.options.TcpPort != 0 {
		args = append(args, "--port-tcp", fmt.Sprintf("%d", server.options.TcpPort))
	}

	if server.options.SpyPort != 0 {
		args = append(args, "--port-spy", fmt.Sprintf("%d", server.options.SpyPort))
	}

	if server.options.ControlPort != 0 {
		args = append(args, "--port-ctrl", fmt.Sprintf("%d", server.options.ControlPort))
	}

	// Start the server
	log.Debug("starting vwifi server at %s with args %v", server.binary, args)

	server.cmd = exec.Command(server.binary, args...)

	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	server.cmd.Stdout = io.MultiWriter(
		minilogWriter{
			level:  log.INFO,
			prefix: "vwifi-server",
		},
		stdoutBuffer,
	)
	server.cmd.Stderr = io.MultiWriter(
		minilogWriter{
			level:  log.ERROR,
			prefix: "vwifi-server",
		},
		stderrBuffer,
	)

	err := server.cmd.Start()

	if err != nil {
		return nil, err
	}

	// Wait for the server to be ready
	var status *ServerStatus = nil
	for i := 0; i < maxServerStatusChecks; i++ {
		// Parse the status
		status, err = server.parseStatus(stdoutBuffer.String())

		if err == nil && !errors.Is(err, errServerStatusNotMatched) {
			return status, nil
		} else if err != nil {
			time.Sleep(serverStatusWait)
		} else {
			goto serverReady
		}
	}

	return nil, errors.New("server did not start in time")

serverReady:
	return status, nil
}

// Stop the server
func (server *Server) Stop() error {
	log.Debug("stopping vwifi server")

	if server.cmd == nil {
		return errors.New("server is not running")
	}

	err := server.cmd.Process.Kill()

	if err != nil {
		return err
	}

	server.cmd = nil

	return nil
}
