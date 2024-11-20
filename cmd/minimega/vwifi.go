// Copyright 2024 Colorado School of Mines CSCI 370 FA24 NREL 2 Group

package main

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/sandia-minimega/minimega/v2/internal/vwifi"
	log "github.com/sandia-minimega/minimega/v2/pkg/minilog"
)

// vwifiServerBinary is the vwifi server binary
const vwifiServerBinary = "vwifi-server"

// vwifiControllerBinary is the vwifi controller binary
const vwifiControllerBinary = "vwifi-ctrl"

// vwifiServer is the vwifi server instance
var vwifiServer *vwifi.Server = nil

// vwifiController is the vwifi controller instance
var vwifiController *vwifi.Controller = nil

// loopbackIP is the loopback IPv4 address
var loopbackIP = net.IPv4(127, 0, 0, 1)

// startVwifiServer starts the vwifi server
func startVwifiServer() error {
	if !*f_enableVwifi {
		return nil
	}

	if vwifiServer != nil || vwifiController != nil {
		return fmt.Errorf("vwifi server is already running")
	}

	// Locate the binaries
	vwifiServerPath, err := exec.LookPath(vwifiServerBinary)

	if err != nil {
		return fmt.Errorf("could not find vwifi server binary (%s): %e", vwifiServerBinary, err)
	}

	vwifiControllerPath, err := exec.LookPath(vwifiControllerBinary)

	if err != nil {
		return fmt.Errorf("could not find vwifi controller binary (%s): %e", vwifiControllerBinary, err)
	}

	log.Info("starting vwifi with server binary %s and controller binary %s", vwifiServerPath, vwifiControllerPath)

	// Start the server
	vwifiServer = vwifi.NewServer(vwifiServerPath, vwifi.ServerOptions{
		PacketLoss: *f_vwifiPacketLoss,
	})

	serverStatus, err := vwifiServer.Start()

	if err != nil {
		return err
	}

	// Start the controller
	vwifiController = vwifi.NewController(vwifiControllerPath, vwifi.ControllerOptions{
		IP:   loopbackIP,
		Port: serverStatus.ControlPort,
	})

	return nil
}

// stopVwifiServer stops the vwifi server
func stopVwifiServer() error {
	if vwifiServer == nil {
		return nil
	}

	log.Info("stopping vwifi")

	err := vwifiServer.Stop()

	if err != nil {
		return err
	}

	return nil
}
