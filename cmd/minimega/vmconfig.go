// Copyright (2012) Sandia Corporation.
// Under the terms of Contract DE-AC04-94AL85000 with Sandia Corporation,
// the U.S. Government retains certain rights in this software.
//
//go:generate ../../bin/vmconfiger -type BaseConfig,KVMConfig,ContainerConfig

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/sandia-minimega/minimega/v2/internal/bridge"
	"github.com/sandia-minimega/minimega/v2/internal/qemu"
	log "github.com/sandia-minimega/minimega/v2/pkg/minilog"
)

// VMConfig contains all the configs possible for a VM. When a VM of a
// particular kind is launched, only the pertinent configuration is copied so
// fields from other configs will have the zero value for the field type.
type VMConfig struct {
	BaseConfig
	KVMConfig
	ContainerConfig
}

func NewVMConfig() VMConfig {
	c := VMConfig{}
	c.Clear(Wildcard)
	return c
}

type ConfigWriter interface {
	WriteConfig(io.Writer) error
}

type ConfigReader interface {
	ReadConfig(io.ReadSeeker, string) error
}

// BaseConfig contains all fields common to all VM types.
type BaseConfig struct {
	// Configures the UUID for a virtual machine. If not set, the VM will be
	// given a random one when it is launched.
	UUID string

	// Configures the number of virtual CPUs to allocate for a VM.
	//
	// Default: 1
	VCPUs uint64

	// Configures the amount of physical memory to allocate (in megabytes).
	//
	// Default: 2048
	Memory uint64

	// Enable or disable snapshot mode for disk images and container
	// filesystems. When enabled, disks/filesystems will have temporary snapshots created
	// when run and changes will not be saved. This allows a single
	// disk/filesystem to be used for many VMs.
	//
	// Default: true
	Snapshot bool

	// Set a host where the VM should be scheduled.
	//
	// Note: Cannot specify Schedule and Colocate in the same config.
	Schedule string `validate:"validSchedule" suggest:"wrapHostnameSuggest(true, false, false)"`

	// Colocate this VM with another VM that has already been launched or is
	// queued for launching.
	//
	// Note: Cannot specify Colocate and Schedule in the same
	Colocate string `validate:"validColocate"`

	// Set a limit on the number of VMs that should be scheduled on the same
	// host as the VM. A limit of zero means that the VM should be scheduled by
	// itself. A limit of -1 means that there is no limit. This is only used
	// when launching VMs in a namespace.
	//
	// Default: -1
	Coschedule int64

	// Enable/disable serial command and control layer for this VM.
	//
	// Default: true
	Backchannel bool

	// Networks for the VM, handler is not generated by vmconfiger.
	Networks NetConfigs

	// Bonds for the VM. Key is bond name, value is BondConfig.
	Bonds BondConfigs

	// Set tags in the same manner as "vm tag". These tags will apply to all
	// newly launched VMs.
	//
	// Default: empty map
	Tags map[string]string
}

func (old VMConfig) Copy() VMConfig {
	return VMConfig{
		BaseConfig:      old.BaseConfig.Copy(),
		KVMConfig:       old.KVMConfig.Copy(),
		ContainerConfig: old.ContainerConfig.Copy(),
	}
}

func (vm VMConfig) String(namespace string) string {
	return vm.BaseConfig.String(namespace) +
		vm.KVMConfig.String() +
		vm.ContainerConfig.String()
}

func (vm *VMConfig) Clear(mask string) {
	vm.BaseConfig.Clear(mask)
	vm.KVMConfig.Clear(mask)
	vm.ContainerConfig.Clear(mask)
}

func (vm *VMConfig) WriteConfig(w io.Writer) error {
	funcs := []func(io.Writer) error{
		vm.BaseConfig.WriteConfig,
		vm.KVMConfig.WriteConfig,
		vm.ContainerConfig.WriteConfig,
	}

	for _, fn := range funcs {
		if err := fn(w); err != nil {
			return err
		}
	}

	return nil
}

func (vm *VMConfig) ReadConfig(r io.ReadSeeker, ns string) error {
	vm.BaseConfig.ReadConfig(r, ns)
	r.Seek(0, io.SeekStart)
	vm.KVMConfig.ReadConfig(r, ns)
	r.Seek(0, io.SeekStart)
	vm.ContainerConfig.ReadConfig(r, ns)

	return nil
}

func (old BaseConfig) Copy() BaseConfig {
	// Copy all fields
	res := old

	// Make deep copy of slices
	res.Networks = make(NetConfigs, len(old.Networks))
	copy(res.Networks, old.Networks)

	// Make deep copy of bonds
	res.Bonds = make(BondConfigs, len(old.Bonds))
	for i, b := range old.Bonds {
		res.Bonds[i] = b.Copy()
	}

	// Make deep copy of tags
	res.Tags = map[string]string{}
	for k, v := range old.Tags {
		res.Tags[k] = v
	}

	return res
}

func (vm *BaseConfig) String(namespace string) string {
	// create output
	var o bytes.Buffer
	fmt.Fprintln(&o, "VM configuration:")
	w := new(tabwriter.Writer)
	w.Init(&o, 5, 0, 1, ' ', 0)
	fmt.Fprintf(w, "Memory:\t%v\n", vm.Memory)
	fmt.Fprintf(w, "VCPUs:\t%v\n", vm.VCPUs)
	fmt.Fprintf(w, "Networks:\t%v\n", vm.NetworkString(namespace))
	fmt.Fprintf(w, "Bonds:\t%v\n", vm.BondString(namespace))
	fmt.Fprintf(w, "Snapshot:\t%v\n", vm.Snapshot)
	fmt.Fprintf(w, "UUID:\t%v\n", vm.UUID)
	fmt.Fprintf(w, "Schedule host:\t%v\n", vm.Schedule)
	fmt.Fprintf(w, "Coschedule limit:\t%v\n", vm.Coschedule)
	fmt.Fprintf(w, "Colocate:\t%v\n", vm.Colocate)
	fmt.Fprintf(w, "Backchannel:\t%v\n", vm.Backchannel)
	if vm.Tags != nil {
		fmt.Fprintf(w, "Tags:\t%v\n", marshal(vm.Tags))
	} else {
		fmt.Fprint(w, "Tags:\t{}\n")
	}
	w.Flush()
	fmt.Fprintln(&o)
	return o.String()
}

func (vm *BaseConfig) NetworkString(namespace string) string {
	return fmt.Sprintf("[%s]", vm.Networks.String())
}

func (vm *BaseConfig) BondString(namespace string) string {
	return fmt.Sprintf("[%s]", vm.Bonds.String())
}

func (vm *BaseConfig) QosString(b, t, i string) string {
	var val string
	br, err := getBridge(b)
	if err != nil {
		return val
	}

	ops := br.GetQos(t)
	if ops == nil {
		return ""
	}

	val += fmt.Sprintf("%s: ", i)
	for _, op := range ops {
		if op.Type == bridge.Delay {
			val += fmt.Sprintf("delay %s ", op.Value)
		}
		if op.Type == bridge.Loss {
			val += fmt.Sprintf("loss %s ", op.Value)
		}
		if op.Type == bridge.Rate {
			val += fmt.Sprintf("rate %s ", op.Value)
		}
	}
	return strings.Trim(val, " ")
}

func (vm *BaseConfig) ReadFieldConfig(r io.Reader, field, namespace string) error {
	switch field {
	case "networks":
		ns := GetOrCreateNamespace(namespace)

		// get valid NIC drivers for current qemu/machine
		nics, err := qemu.NICs(ns.vmConfig.QemuPath, ns.vmConfig.Machine)
		if err != nil {
			if strings.Contains(err.Error(), "executable file not found in $PATH") {
				// warn on not finding kvm because we may just be using containers,
				// otherwise throw a regular error
				log.Warnln(err)
			} else {
				return err
			}
		}

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "vm config networks") {
				continue
			}

			specs := strings.Fields(line)[3:]

			for _, spec := range specs {
				nic, err := ParseNetConfig(spec, nics)
				if err != nil {
					log.Warnln(err) // ??
					continue
				}

				vlan, err := lookupVLAN(namespace, nic.Alias)
				if err != nil {
					log.Warnln(err) // ??
					continue
				}

				nic.VLAN = vlan
				nic.Raw = spec

				vm.Networks = append(vm.Networks, *nic)
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	case "bonds":
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "vm config bonds") {
				continue
			}

			specs := strings.Fields(line)[3:]

			for _, spec := range specs {
				bond, err := ParseBondConfig(spec)
				if err != nil {
					log.Warnln(err) // ??
					continue
				}

				bond.VLAN = -1

				for _, idx := range bond.Interfaces {
					if len(vm.Networks) <= idx {
						log.Warn("no such interface %v for vm %v", idx, vm.UUID) // ??
						continue
					}

					cfg := vm.Networks[idx]

					if cfg.Wifi {
						return fmt.Errorf("cannot bond wifi interfaces")
					}

					if bond.Bridge == "" {
						bond.Bridge = cfg.Bridge
					} else if cfg.Bridge != bond.Bridge {
						return fmt.Errorf("interfaces being bonded are not on the same bridge")
					}

					if bond.VLAN < 0 {
						bond.VLAN = cfg.VLAN
					} else if cfg.VLAN != bond.VLAN {
						log.Warn("interface %d on vm %s is not on VLAN %d -- still defaulting to %d for bond", idx, vm.UUID, bond.VLAN, bond.VLAN)
					}

					bond.QinQ = bond.QinQ || cfg.QinQ
				}

				bond.Raw = spec

				vm.Bonds = append(vm.Bonds, *bond)
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func validSchedule(vmConfig VMConfig, s string) error {
	if vmConfig.Colocate != "" && s != "" {
		return errors.New("cannot specify schedule and colocate in the same config")
	}

	if s == "localhost" {
		s = hostname
	}

	// check if s is in the namespace
	ns := GetNamespace()

	if !ns.Hosts[s] {
		return fmt.Errorf("host is not in namespace: %v", s)
	}

	return nil
}

func validColocate(vmConfig VMConfig, s string) error {
	if vmConfig.Schedule != "" && s != "" {
		return errors.New("cannot specify colocate and schedule in the same config")
	}

	// TODO: could check if s is a known VM
	return nil
}
