// Copyright 2024 Colorado School of Mines CSCI 370 FA24 NREL 2 Group

package vwifi

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sandia-minimega/minimega/v2/pkg/minilog"
)

// ControllerOptions are the vwifi-ctrl command line options
type ControllerOptions struct {
	// The IP address of the server to control
	IP net.IP

	// The port of the server to control
	Port uint16
}

// Controller is a vwifi-ctrl instance
type Controller struct {
	// The controller binary path
	binary string

	// The controller options
	options ControllerOptions
}

// NewController creates a new vwifi-ctrl instance
func NewController(binary string, options ControllerOptions) *Controller {
	return &Controller{
		binary:  binary,
		options: options,
	}
}

// run runs a command against the controller, returning the stdout and stderr
func (controller *Controller) run(additionalArgs []string) (*bytes.Buffer, *bytes.Buffer, error) {
	// Generate the arguments
	args := []string{}

	if controller.options.IP != nil {
		args = append(args, "--ip", controller.options.IP.String())
	}

	if controller.options.Port != 0 {
		args = append(args, "--port", fmt.Sprintf("%d", controller.options.Port))
	}

	args = append(args, additionalArgs...)

	// Start the client
	log.Debug("running vwifi-ctrl at %s with args %v", controller.binary, args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, controller.binary, args...)

	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	cmd.Stdout = io.MultiWriter(
		minilogWriter{
			level:  log.INFO,
			prefix: "vwifi-ctrl",
		},
		stdoutBuffer,
	)
	cmd.Stderr = io.MultiWriter(
		minilogWriter{
			level:  log.ERROR,
			prefix: "vwifi-ctrl",
		},
		stderrBuffer,
	)

	err := cmd.Run()

	if err != nil {
		return nil, nil, err
	}

	return stdoutBuffer, stderrBuffer, nil
}

// ControllerClient is a client connected to the server as reported by the controller
type ControllerClient struct {
	// Whether or not the client is the spy
	Spy bool

	// The client CID
	CID uint32

	// The client's name (if any)
	Name string

	// The client's X position (in meters; meaningless if the client is the spy)
	X int32

	// The client's Y position (in meters; meaningless if the client is the spy)
	Y int32

	// The client's Z position (in meters; meaningless if the client is the spy)
	Z int32
}

// spyClientLinePattern is the pattern to match a line containing a spy client in the list clients output
var spyClientLinePattern = regexp.MustCompile(`^S:(\d+)(?: \(([^)]+)\))?$`)

// regularClientLinePattern is the pattern to match a line containing a regular client in the list clients output
var regularClientLinePattern = regexp.MustCompile(`^(\d+)(?: \(([^)]+)\))? (\d+) (\d+) (\d+)$`)

// parseClients parses the clients output as reported by the controller
func (controller *Controller) parseClients(raw string) ([]ControllerClient, error) {
	// Parse the output
	clients := []ControllerClient{}

	// Iterate over line by line
	scanner := bufio.NewScanner(strings.NewReader(raw))

	for scanner.Scan() {
		line := scanner.Text()

		// Parse the client
		if spyMatches := spyClientLinePattern.FindStringSubmatch(line); spyMatches != nil {
			cid, err := strconv.ParseUint(spyMatches[1], 10, 32)

			if err != nil {
				return nil, err
			}

			clients = append(clients, ControllerClient{
				Spy:  true,
				CID:  uint32(cid),
				Name: spyMatches[2],
			})
		} else if regularMatches := regularClientLinePattern.FindStringSubmatch(line); regularMatches != nil {
			cid, err := strconv.ParseUint(regularMatches[1], 10, 32)

			if err != nil {
				return nil, err
			}

			x, err := strconv.ParseInt(regularMatches[3], 10, 32)

			if err != nil {
				return nil, err
			}

			y, err := strconv.ParseInt(regularMatches[4], 10, 32)

			if err != nil {
				return nil, err
			}

			z, err := strconv.ParseInt(regularMatches[5], 10, 32)

			if err != nil {
				return nil, err
			}

			clients = append(clients, ControllerClient{
				Spy:  false,
				CID:  uint32(cid),
				Name: regularMatches[2],
				X:    int32(x),
				Y:    int32(y),
				Z:    int32(z),
			})
		}
	}

	return clients, nil
}

// ListClients lists the clients connected to the server
func (controller *Controller) ListClients() ([]ControllerClient, error) {
	// Run the command
	stdout, _, err := controller.run([]string{"ls"})

	if err != nil {
		return nil, err
	}

	// Parse the clients
	clients, err := controller.parseClients(stdout.String())

	if err != nil {
		return nil, err
	}

	return clients, nil
}

// MoveClient moves a client to a new position
func (controller *Controller) MoveClient(cid uint32, x int32, y int32, z int32) error {
	// Run the command
	args := []string{"set", fmt.Sprintf("%d", cid), fmt.Sprintf("%d", x), fmt.Sprintf("%d", y), fmt.Sprintf("%d", z)}

	_, _, err := controller.run(args)

	return err
}

// SetName sets the name of a client
func (controller *Controller) SetName(cid uint32, name string) error {
	// Run the command
	args := []string{"setname", fmt.Sprintf("%d", cid), name}

	_, _, err := controller.run(args)

	return err
}

// SetPacketLoss sets the packet loss
func (controller *Controller) SetPacketLoss(loosePackets bool) error {
	// Run the command
	args := []string{"loss"}

	if loosePackets {
		args = append(args, "yes")
	} else {
		args = append(args, "no")
	}

	_, _, err := controller.run(args)

	return err
}

// ControllerServerStatus is the status of the server as reported by the controller
type ControllerServerStatus struct {
	// The listening IP address of the server's control socket
	ControlIP net.IP

	// The listening port of the server's control socket
	ControlPort uint16

	// Whether or not the server is simulating packet loss
	PacketLoss bool

	// The simulation's physical scale
	Scale float64

	// The listening port for the VSOCK socket
	VsockPort uint16

	// The listening port for the TCP socket
	TcpPort uint16

	// No clue what this is, but it's in the output
	SizeOfDisconnected uint32

	// Whether or not the server is connected to a spy
	SpyConnected bool
}

// controllerServerStatusPattern is the pattern to match the status output as reported by the controller
var controllerServerStatusPattern = regexp.MustCompile(`^\s*CTRL : IP : (.+)\nCTRL : Port : (\d+)\nSRV : PacketLoss : (Enable|Disable)\nSRV : Scale : ([+-]?\d+(?:\.\d+)?)\nSRV VHOST : Port : (\d+)\nSRV INET : Port : (\d+)\nSRV : SizeOfDisconnected : (\d+)\nSPY : (Connected|Disconnected)\s*$`)

// parseStatus parses the status output as reported by the controller
func (controller *Controller) parseStatus(raw string) (*ControllerServerStatus, error) {
	matches := controllerServerStatusPattern.FindStringSubmatch(raw)

	if matches == nil {
		return nil, fmt.Errorf("failed to parse status output: %s", raw)
	}

	controlPort, err := strconv.ParseUint(matches[2], 10, 16)

	if err != nil {
		return nil, err
	}

	scale, err := strconv.ParseFloat(matches[4], 64)

	if err != nil {
		return nil, err
	}

	vsockPort, err := strconv.ParseUint(matches[5], 10, 16)

	if err != nil {
		return nil, err
	}

	tcpPort, err := strconv.ParseUint(matches[6], 10, 16)

	if err != nil {
		return nil, err
	}

	sizeOfDisconnected, err := strconv.ParseUint(matches[7], 10, 32)

	if err != nil {
		return nil, err
	}

	status := ControllerServerStatus{
		ControlIP:          net.ParseIP(matches[1]),
		ControlPort:        uint16(controlPort),
		PacketLoss:         matches[3] == "Enable",
		Scale:              scale,
		VsockPort:          uint16(vsockPort),
		TcpPort:            uint16(tcpPort),
		SizeOfDisconnected: uint32(sizeOfDisconnected),
		SpyConnected:       matches[8] == "Connected",
	}

	return &status, nil
}

// Status gets the status of the server
func (controller *Controller) Status() (*ControllerServerStatus, error) {
	// Run the command
	stdout, _, err := controller.run([]string{"status"})

	if err != nil {
		return nil, err
	}

	// Parse the status
	status, err := controller.parseStatus(stdout.String())

	if err != nil {
		return nil, err
	}

	return status, nil
}

// SetScale sets the scale of the simulation
func (controller *Controller) SetScale(scale float64) error {
	// Run the command
	args := []string{"scale", fmt.Sprintf("%f", scale)}

	_, _, err := controller.run(args)

	return err
}
