// Code generated by "vmconfiger -type BaseConfig,KVMConfig,ContainerConfig"; DO NOT EDIT

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sandia-minimega/minimega/v2/pkg/minicli"
)

var vmconfigerCLIHandlers = []minicli.Handler{
	{
		HelpShort: "configures filesystem",
		HelpLong: `Configure the filesystem to use for launching a container. This should
be a root filesystem for a linux distribution (containing /dev, /proc,
/sys, etc.)

Note: this configuration only applies to containers and must be specified.
`,
		Patterns: []string{
			"vm config filesystem [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.FilesystemPath
				return nil
			}

			v := checkPath(c.StringArgs["value"])

			ns.vmConfig.FilesystemPath = v

			return nil
		}),
	},
	{
		HelpShort: "configures hostname",
		HelpLong: `Set a hostname for a container before launching the init program. If not
set, the hostname will be the VM name. The hostname can also be set by
the init program or other root process in the container.

Note: this configuration only applies to containers.
`,
		Patterns: []string{
			"vm config hostname [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.Hostname
				return nil
			}

			ns.vmConfig.Hostname = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures init",
		HelpLong: `Set the init program and args to exec into upon container launch. This
will be PID 1 in the container.

Note: this configuration only applies to containers.

Default: "/init"
`,
		Patterns: []string{
			"vm config init [value]...",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.ListArgs) == 0 {
				if len(ns.vmConfig.Init) == 0 {
					return nil
				}

				r.Response = fmt.Sprintf("%v", ns.vmConfig.Init)
				return nil
			}

			ns.vmConfig.Init = c.ListArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures preinit",
		HelpLong: `Containers start in a highly restricted environment. vm config preinit
allows running processes before isolation mechanisms are enabled. This
occurs when the vm is launched and before the vm is put in the building
state. preinit processes must finish before the vm will be allowed to
start.

Specifically, the preinit command will be run after entering namespaces,
and mounting dependent filesystems, but before cgroups and root
capabilities are set, and before entering the chroot. This means that
the preinit command is run as root and can control the host.

For example, to run a script that enables ip forwarding, which is not
allowed during runtime because /proc is mounted read-only, add a preinit
script:

	vm config preinit enable_ip_forwarding.sh

Note: this configuration only applies to containers.
`,
		Patterns: []string{
			"vm config preinit [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.Preinit
				return nil
			}

			ns.vmConfig.Preinit = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures fifos",
		HelpLong: `Set the number of named pipes to include in the container for
container-host communication. Named pipes will appear on the host in the
instance directory for the container as fifoN, and on the container as
/dev/fifos/fifoN.

Fifos are created using mkfifo() and have all of the same usage
constraints.

Note: this configuration only applies to containers.
`,
		Patterns: []string{
			"vm config fifos [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = strconv.FormatUint(ns.vmConfig.Fifos, 10)
				return nil
			}

			i, err := strconv.ParseUint(c.StringArgs["value"], 10, 64)
			if err != nil {
				return err
			}

			ns.vmConfig.Fifos = i

			return nil
		}),
	},
	{
		HelpShort: "configures volume",
		HelpLong: `Attach one or more volumes to a container. These directories will be
mounted inside the container at the specified location.

For example, to mount /scratch/data to /data inside the container:

 vm config volume /data /scratch/data

Commands with the same <key> will overwrite previous volumes:

 vm config volume /data /scratch/data2
 vm config volume /data
 /scratch/data2

Note: this configuration only applies to containers.

Default: empty map
`,
		Patterns: []string{
			"vm config volume",
			"vm config volume <key> [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if c.StringArgs["key"] == "" {
				var b bytes.Buffer

				for k, v := range ns.vmConfig.VolumePaths {
					fmt.Fprintf(&b, "%v -> %v\n", k, v)
				}

				r.Response = b.String()
				return nil
			}

			if c.StringArgs["value"] == "" {
				if ns.vmConfig.VolumePaths != nil {
					r.Response = ns.vmConfig.VolumePaths[c.StringArgs["value"]]
				}
				return nil
			}

			if ns.vmConfig.VolumePaths == nil {
				ns.vmConfig.VolumePaths = make(map[string]string)
			}

			v := checkPath(c.StringArgs["value"])

			ns.vmConfig.VolumePaths[c.StringArgs["key"]] = v

			return nil
		}),
	},
	{
		HelpShort: "configures qemu",
		HelpLong: `Set the QEMU binary name to invoke. Relative paths are ok.

Note: this configuration only applies to KVM-based VMs.

Default: "kvm"
`,
		Patterns: []string{
			"vm config qemu [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.QemuPath
				return nil
			}

			v := checkPath(c.StringArgs["value"])

			ns.vmConfig.QemuPath = v

			return nil
		}),
	},
	{
		HelpShort: "configures kernel",
		HelpLong: `Attach a kernel image to a VM. If set, QEMU will boot from this image
instead of any disk image.

Note: this configuration only applies to KVM-based VMs.
`,
		Patterns: []string{
			"vm config kernel [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.KernelPath
				return nil
			}

			v := checkPath(c.StringArgs["value"])

			ns.vmConfig.KernelPath = v

			return nil
		}),
	},
	{
		HelpShort: "configures initrd",
		HelpLong: `Attach an initrd image to a VM. Passed along with the kernel image at
boot time.

Note: this configuration only applies to KVM-based VMs.
`,
		Patterns: []string{
			"vm config initrd [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.InitrdPath
				return nil
			}

			v := checkPath(c.StringArgs["value"])

			ns.vmConfig.InitrdPath = v

			return nil
		}),
	},
	{
		HelpShort: "configures cdrom",
		HelpLong: `Attach a cdrom to a VM. When using a cdrom, it will automatically be set
to be the boot device.

Note: this configuration only applies to KVM-based VMs.
`,
		Patterns: []string{
			"vm config cdrom [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.CdromPath
				return nil
			}

			v := checkPath(c.StringArgs["value"])

			ns.vmConfig.CdromPath = v

			return nil
		}),
	},
	{
		HelpShort: "configures migrate",
		HelpLong: `Assign a migration image, generated by a previously saved VM to boot
with. By default, images are read from the files directory as specified
with -filepath. This can be overridden by using an absolute path.
Migration images should be booted with a kernel/initrd, disk, or cdrom.
Use 'vm migrate' to generate migration images from running VMs.

Note: this configuration only applies to KVM-based VMs.
`,
		Patterns: []string{
			"vm config migrate [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.MigratePath
				return nil
			}

			v := checkPath(c.StringArgs["value"])

			ns.vmConfig.MigratePath = v

			return nil
		}),
	},
	{
		HelpShort: "configures cpu",
		HelpLong: `Set the virtual CPU architecture.

By default, set to 'host' which matches the host CPU. See 'qemu -cpu
help' for a list of supported CPUs.

The accepted values for this configuration depend on the QEMU binary
name specified by 'vm config qemu'.

Note: this configuration only applies to KVM-based VMs.

Default: "host"
`,
		Patterns: []string{
			"vm config cpu [value]",
		},

		Suggest: wrapSuggest(suggestCPU),

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.CPU
				return nil
			}

			if err := validCPU(ns.vmConfig, c.StringArgs["value"]); err != nil {
				return err
			}

			ns.vmConfig.CPU = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures sockets",
		HelpLong: `Set the number of CPU sockets. If unspecified, QEMU will calculate
missing values based on vCPUs, cores, and threads.
`,
		Patterns: []string{
			"vm config sockets [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = strconv.FormatUint(ns.vmConfig.Sockets, 10)
				return nil
			}

			i, err := strconv.ParseUint(c.StringArgs["value"], 10, 64)
			if err != nil {
				return err
			}

			ns.vmConfig.Sockets = i

			return nil
		}),
	},
	{
		HelpShort: "configures cores",
		HelpLong: `Set the number of CPU cores per socket. If unspecified, QEMU will
calculate missing values based on vCPUs, sockets, and threads.
`,
		Patterns: []string{
			"vm config cores [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = strconv.FormatUint(ns.vmConfig.Cores, 10)
				return nil
			}

			i, err := strconv.ParseUint(c.StringArgs["value"], 10, 64)
			if err != nil {
				return err
			}

			ns.vmConfig.Cores = i

			return nil
		}),
	},
	{
		HelpShort: "configures threads",
		HelpLong: `Set the number of CPU threads per core. If unspecified, QEMU will
calculate missing values based on vCPUs, sockets, and cores.
`,
		Patterns: []string{
			"vm config threads [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = strconv.FormatUint(ns.vmConfig.Threads, 10)
				return nil
			}

			i, err := strconv.ParseUint(c.StringArgs["value"], 10, 64)
			if err != nil {
				return err
			}

			ns.vmConfig.Threads = i

			return nil
		}),
	},
	{
		HelpShort: "configures machine",
		HelpLong: `Specify the machine type. See 'qemu -M help' for a list supported
machine types.

The accepted values for this configuration depend on the QEMU binary
name specified by 'vm config qemu'.

Note: this configuration only applies to KVM-based VMs.
`,
		Patterns: []string{
			"vm config machine [value]",
		},

		Suggest: wrapSuggest(suggestMachine),

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.Machine
				return nil
			}

			if err := validMachine(ns.vmConfig, c.StringArgs["value"]); err != nil {
				return err
			}

			ns.vmConfig.Machine = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures serial-ports",
		HelpLong: `Specify the serial ports that will be created for the VM to use. Serial
ports specified will be mapped to the VM's /dev/ttySX device, where X
refers to the connected unix socket on the host at
$minimega_runtime/<vm_id>/serialX.

Examples:

To display current serial ports:
  vm config serial-ports

To create three serial ports:
  vm config serial-ports 3

Note: Whereas modern versions of Windows support up to 256 COM ports,
Linux typically only supports up to four serial devices. To use more,
make sure to pass "8250.n_uarts = 4" to the guest Linux kernel at boot.
Replace 4 with another number.
`,
		Patterns: []string{
			"vm config serial-ports [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = strconv.FormatUint(ns.vmConfig.SerialPorts, 10)
				return nil
			}

			i, err := strconv.ParseUint(c.StringArgs["value"], 10, 64)
			if err != nil {
				return err
			}

			ns.vmConfig.SerialPorts = i

			return nil
		}),
	},
	{
		HelpShort: "configures virtio-ports",
		HelpLong: `Specify the virtio-serial ports that will be created for the VM to use.
Virtio-serial ports specified will be mapped to the VM's
/dev/virtio-port/<portname> device, where <portname> refers to the
connected unix socket on the host at
$minimega_runtime/<vm_id>/virtio-serialX.

Examples:

To display current virtio-serial ports:
  vm config virtio-ports

To create three virtio-serial ports:
  vm config virtio-ports 3

To explicitly name the virtio-ports, pass a comma-separated list of names:

  vm config virtio-ports foo,bar

The ports (on the guest) will then be mapped to /dev/virtio-port/foo and
/dev/virtio-port/bar.
`,
		Patterns: []string{
			"vm config virtio-ports [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.VirtioPorts
				return nil
			}

			ns.vmConfig.VirtioPorts = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures vga",
		HelpLong: `Specify the graphics card to emulate. "cirrus" or "std" should work with
most operating systems.

Default: "std"
`,
		Patterns: []string{
			"vm config vga [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.Vga
				return nil
			}

			ns.vmConfig.Vga = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures append",
		HelpLong: `Add an append string to a kernel set with vm kernel. Setting vm append
without using vm kernel will result in an error.

For example, to set a static IP for a linux VM:

	vm config append ip=10.0.0.5 gateway=10.0.0.1 netmask=255.255.255.0 dns=10.10.10.10

Note: this configuration only applies to KVM-based VMs.
`,
		Patterns: []string{
			"vm config append [value]...",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.ListArgs) == 0 {
				if len(ns.vmConfig.Append) == 0 {
					return nil
				}

				r.Response = fmt.Sprintf("%v", ns.vmConfig.Append)
				return nil
			}

			ns.vmConfig.Append = c.ListArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures usb-use-xhci",
		HelpLong: `If true will use xHCI USB controller. Otherwise will use EHCI.
EHCI does not support USB 3.0, but may be used for backwards compatibility.

Default: true
`,
		Patterns: []string{
			"vm config usb-use-xhci [true,false]",
		},
		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.BoolArgs) == 0 {
				r.Response = strconv.FormatBool(ns.vmConfig.UsbUseXHCI)
				return nil
			}

			ns.vmConfig.UsbUseXHCI = c.BoolArgs["true"]

			return nil
		}),
	},
	{
		HelpShort: "configures tpm-socket",
		HelpLong: `If specified, will configure VM to use virtual Trusted Platform Module (TPM)
socket at the path provided
`,
		Patterns: []string{
			"vm config tpm-socket [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.TpmSocketPath
				return nil
			}

			v := checkPath(c.StringArgs["value"])

			ns.vmConfig.TpmSocketPath = v

			return nil
		}),
	},
	{
		HelpShort: "configures bidirectional-copy-paste",
		HelpLong: `Enables bidirectional copy paste instead of basic pasting into VM.
Requires QEMU 6.1+ compiled with qemu-vdagent chardev and for spice-vdagent to be installed on VM.

Default: false
`,
		Patterns: []string{
			"vm config bidirectional-copy-paste [true,false]",
		},
		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.BoolArgs) == 0 {
				r.Response = strconv.FormatBool(ns.vmConfig.BidirectionalCopyPaste)
				return nil
			}

			ns.vmConfig.BidirectionalCopyPaste = c.BoolArgs["true"]

			return nil
		}),
	},
	{
		HelpShort: "configures qemu-append",
		HelpLong: `Add additional arguments to be passed to the QEMU instance. For example:

	vm config qemu-append -serial tcp:localhost:4001

Note: this configuration only applies to KVM-based VMs.
`,
		Patterns: []string{
			"vm config qemu-append [value]...",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.ListArgs) == 0 {
				if len(ns.vmConfig.QemuAppend) == 0 {
					return nil
				}

				r.Response = fmt.Sprintf("%v", ns.vmConfig.QemuAppend)
				return nil
			}

			ns.vmConfig.QemuAppend = c.ListArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures uuid",
		HelpLong: `Configures the UUID for a virtual machine. If not set, the VM will be
given a random one when it is launched.
`,
		Patterns: []string{
			"vm config uuid [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.UUID
				return nil
			}

			ns.vmConfig.UUID = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures vcpus",
		HelpLong: `Configures the number of virtual CPUs to allocate for a VM.

Default: 1
`,
		Patterns: []string{
			"vm config vcpus [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = strconv.FormatUint(ns.vmConfig.VCPUs, 10)
				return nil
			}

			i, err := strconv.ParseUint(c.StringArgs["value"], 10, 64)
			if err != nil {
				return err
			}

			ns.vmConfig.VCPUs = i

			return nil
		}),
	},
	{
		HelpShort: "configures memory",
		HelpLong: `Configures the amount of physical memory to allocate (in megabytes).

Default: 2048
`,
		Patterns: []string{
			"vm config memory [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = strconv.FormatUint(ns.vmConfig.Memory, 10)
				return nil
			}

			i, err := strconv.ParseUint(c.StringArgs["value"], 10, 64)
			if err != nil {
				return err
			}

			ns.vmConfig.Memory = i

			return nil
		}),
	},
	{
		HelpShort: "configures snapshot",
		HelpLong: `Enable or disable snapshot mode for disk images and container
filesystems. When enabled, disks/filesystems will have temporary snapshots created
when run and changes will not be saved. This allows a single
disk/filesystem to be used for many VMs.

Default: true
`,
		Patterns: []string{
			"vm config snapshot [true,false]",
		},
		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.BoolArgs) == 0 {
				r.Response = strconv.FormatBool(ns.vmConfig.Snapshot)
				return nil
			}

			ns.vmConfig.Snapshot = c.BoolArgs["true"]

			return nil
		}),
	},
	{
		HelpShort: "configures schedule",
		HelpLong: `Set a host where the VM should be scheduled.

Note: Cannot specify Schedule and Colocate in the same config.
`,
		Patterns: []string{
			"vm config schedule [value]",
		},

		Suggest: wrapHostnameSuggest(true, false, false),

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.Schedule
				return nil
			}

			if err := validSchedule(ns.vmConfig, c.StringArgs["value"]); err != nil {
				return err
			}

			ns.vmConfig.Schedule = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures colocate",
		HelpLong: `Colocate this VM with another VM that has already been launched or is
queued for launching.

Note: Cannot specify Colocate and Schedule in the same
`,
		Patterns: []string{
			"vm config colocate [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = ns.vmConfig.Colocate
				return nil
			}

			if err := validColocate(ns.vmConfig, c.StringArgs["value"]); err != nil {
				return err
			}

			ns.vmConfig.Colocate = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "configures coschedule",
		HelpLong: `Set a limit on the number of VMs that should be scheduled on the same
host as the VM. A limit of zero means that the VM should be scheduled by
itself. A limit of -1 means that there is no limit. This is only used
when launching VMs in a namespace.

Default: -1
`,
		Patterns: []string{
			"vm config coschedule [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.StringArgs) == 0 {
				r.Response = strconv.FormatInt(ns.vmConfig.Coschedule, 10)
				return nil
			}

			i, err := strconv.ParseInt(c.StringArgs["value"], 10, 64)
			if err != nil {
				return err
			}

			ns.vmConfig.Coschedule = i

			return nil
		}),
	},
	{
		HelpShort: "configures backchannel",
		HelpLong: `Enable/disable serial command and control layer for this VM.

Default: true
`,
		Patterns: []string{
			"vm config backchannel [true,false]",
		},
		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if len(c.BoolArgs) == 0 {
				r.Response = strconv.FormatBool(ns.vmConfig.Backchannel)
				return nil
			}

			ns.vmConfig.Backchannel = c.BoolArgs["true"]

			return nil
		}),
	},
	{
		HelpShort: "configures tags",
		HelpLong: `Set tags in the same manner as "vm tag". These tags will apply to all
newly launched VMs.

Default: empty map
`,
		Patterns: []string{
			"vm config tags",
			"vm config tags <key> [value]",
		},

		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			if c.StringArgs["key"] == "" {
				var b bytes.Buffer

				for k, v := range ns.vmConfig.Tags {
					fmt.Fprintf(&b, "%v -> %v\n", k, v)
				}

				r.Response = b.String()
				return nil
			}

			if c.StringArgs["value"] == "" {
				if ns.vmConfig.Tags != nil {
					r.Response = ns.vmConfig.Tags[c.StringArgs["value"]]
				}
				return nil
			}

			if ns.vmConfig.Tags == nil {
				ns.vmConfig.Tags = make(map[string]string)
			}

			ns.vmConfig.Tags[c.StringArgs["key"]] = c.StringArgs["value"]

			return nil
		}),
	},
	{
		HelpShort: "reset one or more configurations to default value",
		Patterns: []string{
			"clear vm config",
			"clear vm config <append,>",
			"clear vm config <backchannel,>",
			"clear vm config <bidirectional-copy-paste,>",
			"clear vm config <bonds,>",
			"clear vm config <cpu,>",
			"clear vm config <cdrom,>",
			"clear vm config <colocate,>",
			"clear vm config <cores,>",
			"clear vm config <coschedule,>",
			"clear vm config <disks,>",
			"clear vm config <fifos,>",
			"clear vm config <filesystem,>",
			"clear vm config <hostname,>",
			"clear vm config <init,>",
			"clear vm config <initrd,>",
			"clear vm config <kernel,>",
			"clear vm config <machine,>",
			"clear vm config <memory,>",
			"clear vm config <migrate,>",
			"clear vm config <networks,>",
			"clear vm config <preinit,>",
			"clear vm config <qemu-append,>",
			"clear vm config <qemu-override,>",
			"clear vm config <qemu,>",
			"clear vm config <schedule,>",
			"clear vm config <serial-ports,>",
			"clear vm config <snapshot,>",
			"clear vm config <sockets,>",
			"clear vm config <tags,>",
			"clear vm config <threads,>",
			"clear vm config <tpm-socket,>",
			"clear vm config <uuid,>",
			"clear vm config <usb-use-xhci,>",
			"clear vm config <vcpus,>",
			"clear vm config <vga,>",
			"clear vm config <virtio-ports,>",
			"clear vm config <volume,>",
		},
		Call: wrapSimpleCLI(func(ns *Namespace, c *minicli.Command, r *minicli.Response) error {
			// at most one key will be set in BoolArgs but we don't know what it
			// will be so we have to loop through the args and set whatever key we
			// see.
			mask := Wildcard
			for k := range c.BoolArgs {
				mask = k
			}

			ns.vmConfig.Clear(mask)

			return nil
		}),
	},
}

func (v *BaseConfig) Info(field string) (string, error) {
	if field == "uuid" {
		return v.UUID, nil
	}
	if field == "vcpus" {
		return strconv.FormatUint(v.VCPUs, 10), nil
	}
	if field == "memory" {
		return strconv.FormatUint(v.Memory, 10), nil
	}
	if field == "snapshot" {
		return strconv.FormatBool(v.Snapshot), nil
	}
	if field == "schedule" {
		return v.Schedule, nil
	}
	if field == "colocate" {
		return v.Colocate, nil
	}
	if field == "coschedule" {
		return fmt.Sprintf("%v", v.Coschedule), nil
	}
	if field == "backchannel" {
		return strconv.FormatBool(v.Backchannel), nil
	}
	if field == "networks" {
		return fmt.Sprintf("%v", v.Networks), nil
	}
	if field == "bonds" {
		return fmt.Sprintf("%v", v.Bonds), nil
	}
	if field == "tags" {
		return fmt.Sprintf("%v", v.Tags), nil
	}

	return "", fmt.Errorf("invalid info field: %v", field)
}

func (v *BaseConfig) Clear(mask string) {
	if mask == Wildcard || mask == "uuid" {
		v.UUID = ""
	}
	if mask == Wildcard || mask == "vcpus" {
		v.VCPUs = 1
	}
	if mask == Wildcard || mask == "memory" {
		v.Memory = 2048
	}
	if mask == Wildcard || mask == "snapshot" {
		v.Snapshot = true
	}
	if mask == Wildcard || mask == "schedule" {
		v.Schedule = ""
	}
	if mask == Wildcard || mask == "colocate" {
		v.Colocate = ""
	}
	if mask == Wildcard || mask == "coschedule" {
		v.Coschedule = -1
	}
	if mask == Wildcard || mask == "backchannel" {
		v.Backchannel = true
	}
	if mask == Wildcard || mask == "networks" {
		v.Networks = NetConfigs{}
	}
	if mask == Wildcard || mask == "bonds" {
		v.Bonds = BondConfigs{}
	}
	if mask == Wildcard || mask == "tags" {
		v.Tags = make(map[string]string)
	}
}

func (v *BaseConfig) WriteConfig(w io.Writer) error {
	if v.UUID != "" {
		fmt.Fprintf(w, "vm config uuid %v\n", v.UUID)
	}
	if v.VCPUs != 1 {
		fmt.Fprintf(w, "vm config vcpus %v\n", v.VCPUs)
	}
	if v.Memory != 2048 {
		fmt.Fprintf(w, "vm config memory %v\n", v.Memory)
	}
	if v.Snapshot != true {
		fmt.Fprintf(w, "vm config snapshot %t\n", v.Snapshot)
	}
	if v.Schedule != "" {
		fmt.Fprintf(w, "vm config schedule %v\n", v.Schedule)
	}
	if v.Colocate != "" {
		fmt.Fprintf(w, "vm config colocate %v\n", v.Colocate)
	}
	if v.Coschedule != -1 {
		fmt.Fprintf(w, "vm config coschedule %v\n", v.Coschedule)
	}
	if v.Backchannel != true {
		fmt.Fprintf(w, "vm config backchannel %t\n", v.Backchannel)
	}
	if err := v.Networks.WriteConfig(w); err != nil {
		return err
	}
	if err := v.Bonds.WriteConfig(w); err != nil {
		return err
	}
	for k, v := range v.Tags {
		fmt.Fprintf(w, "vm config tags %v %v\n", k, v)
	}

	return nil
}

func (v *BaseConfig) ReadConfig(r io.Reader, ns string) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "vm config") {
			continue
		}

		config := strings.Fields(line)[2:]
		field := config[0]

		switch field {
		case "uuid":
			v.UUID = config[1]
		case "vcpus":
			v.VCPUs, _ = strconv.ParseUint(config[1], 10, 64)
		case "memory":
			v.Memory, _ = strconv.ParseUint(config[1], 10, 64)
		case "snapshot":
			v.Snapshot, _ = strconv.ParseBool(config[1])
		case "schedule":
			v.Schedule = config[1]
		case "colocate":
			v.Colocate = config[1]
		case "coschedule":
			v.Coschedule, _ = strconv.ParseInt(config[1], 10, 64)
		case "backchannel":
			v.Backchannel, _ = strconv.ParseBool(config[1])
		case "networks":
			v.ReadFieldConfig(strings.NewReader(line), "networks", ns)
		case "bonds":
			v.ReadFieldConfig(strings.NewReader(line), "bonds", ns)
		case "tags":
			v.Tags[config[1]] = config[2]
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (v *ContainerConfig) Info(field string) (string, error) {
	if field == "filesystem" {
		return v.FilesystemPath, nil
	}
	if field == "hostname" {
		return v.Hostname, nil
	}
	if field == "init" {
		return fmt.Sprintf("%v", v.Init), nil
	}
	if field == "preinit" {
		return v.Preinit, nil
	}
	if field == "fifos" {
		return strconv.FormatUint(v.Fifos, 10), nil
	}
	if field == "volume" {
		return fmt.Sprintf("%v", v.VolumePaths), nil
	}

	return "", fmt.Errorf("invalid info field: %v", field)
}

func (v *ContainerConfig) Clear(mask string) {
	if mask == Wildcard || mask == "filesystem" {
		v.FilesystemPath = ""
	}
	if mask == Wildcard || mask == "hostname" {
		v.Hostname = ""
	}
	if mask == Wildcard || mask == "init" {
		v.Init = []string{"/init"}
	}
	if mask == Wildcard || mask == "preinit" {
		v.Preinit = ""
	}
	if mask == Wildcard || mask == "fifos" {
		v.Fifos = 0
	}
	if mask == Wildcard || mask == "volume" {
		v.VolumePaths = make(map[string]string)
	}
}

func (v *ContainerConfig) WriteConfig(w io.Writer) error {
	if v.FilesystemPath != "" {
		fmt.Fprintf(w, "vm config filesystem %v\n", v.FilesystemPath)
	}
	if v.Hostname != "" {
		fmt.Fprintf(w, "vm config hostname %v\n", v.Hostname)
	}
	if len(v.Init) > 0 {
		fmt.Fprintf(w, "vm config init %v\n", quoteJoin(v.Init, " "))
	}
	if v.Preinit != "" {
		fmt.Fprintf(w, "vm config preinit %v\n", v.Preinit)
	}
	if v.Fifos != 0 {
		fmt.Fprintf(w, "vm config fifos %v\n", v.Fifos)
	}
	for k, v := range v.VolumePaths {
		fmt.Fprintf(w, "vm config volume %v %v\n", k, v)
	}

	return nil
}

func (v *ContainerConfig) ReadConfig(r io.Reader, ns string) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "vm config") {
			continue
		}

		config := strings.Fields(line)[2:]
		field := config[0]

		switch field {
		case "filesystem":
			v.FilesystemPath = config[1]
		case "hostname":
			v.Hostname = config[1]
		case "init":
			v.Init = strings.Fields(config[1])
		case "preinit":
			v.Preinit = config[1]
		case "fifos":
			v.Fifos, _ = strconv.ParseUint(config[1], 10, 64)
		case "volume":
			v.VolumePaths[config[1]] = config[2]
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (v *KVMConfig) Info(field string) (string, error) {
	if field == "qemu" {
		return v.QemuPath, nil
	}
	if field == "kernel" {
		return v.KernelPath, nil
	}
	if field == "initrd" {
		return v.InitrdPath, nil
	}
	if field == "cdrom" {
		return v.CdromPath, nil
	}
	if field == "migrate" {
		return v.MigratePath, nil
	}
	if field == "cpu" {
		return v.CPU, nil
	}
	if field == "sockets" {
		return strconv.FormatUint(v.Sockets, 10), nil
	}
	if field == "cores" {
		return strconv.FormatUint(v.Cores, 10), nil
	}
	if field == "threads" {
		return strconv.FormatUint(v.Threads, 10), nil
	}
	if field == "machine" {
		return v.Machine, nil
	}
	if field == "serial-ports" {
		return strconv.FormatUint(v.SerialPorts, 10), nil
	}
	if field == "virtio-ports" {
		return v.VirtioPorts, nil
	}
	if field == "vga" {
		return v.Vga, nil
	}
	if field == "append" {
		return fmt.Sprintf("%v", v.Append), nil
	}
	if field == "disks" {
		return fmt.Sprintf("%v", v.Disks), nil
	}
	if field == "usb-use-xhci" {
		return strconv.FormatBool(v.UsbUseXHCI), nil
	}
	if field == "tpm-socket" {
		return v.TpmSocketPath, nil
	}
	if field == "bidirectional-copy-paste" {
		return strconv.FormatBool(v.BidirectionalCopyPaste), nil
	}
	if field == "qemu-append" {
		return fmt.Sprintf("%v", v.QemuAppend), nil
	}
	if field == "qemu-override" {
		return fmt.Sprintf("%v", v.QemuOverride), nil
	}

	return "", fmt.Errorf("invalid info field: %v", field)
}

func (v *KVMConfig) Clear(mask string) {
	if mask == Wildcard || mask == "qemu" {
		v.QemuPath = "kvm"
	}
	if mask == Wildcard || mask == "kernel" {
		v.KernelPath = ""
	}
	if mask == Wildcard || mask == "initrd" {
		v.InitrdPath = ""
	}
	if mask == Wildcard || mask == "cdrom" {
		v.CdromPath = ""
	}
	if mask == Wildcard || mask == "migrate" {
		v.MigratePath = ""
	}
	if mask == Wildcard || mask == "cpu" {
		v.CPU = "host"
	}
	if mask == Wildcard || mask == "sockets" {
		v.Sockets = 0
	}
	if mask == Wildcard || mask == "cores" {
		v.Cores = 0
	}
	if mask == Wildcard || mask == "threads" {
		v.Threads = 0
	}
	if mask == Wildcard || mask == "machine" {
		v.Machine = ""
	}
	if mask == Wildcard || mask == "serial-ports" {
		v.SerialPorts = 0
	}
	if mask == Wildcard || mask == "virtio-ports" {
		v.VirtioPorts = ""
	}
	if mask == Wildcard || mask == "vga" {
		v.Vga = "std"
	}
	if mask == Wildcard || mask == "append" {
		v.Append = nil
	}
	if mask == Wildcard || mask == "disks" {
		v.Disks = DiskConfigs{}
	}
	if mask == Wildcard || mask == "usb-use-xhci" {
		v.UsbUseXHCI = true
	}
	if mask == Wildcard || mask == "tpm-socket" {
		v.TpmSocketPath = ""
	}
	if mask == Wildcard || mask == "bidirectional-copy-paste" {
		v.BidirectionalCopyPaste = false
	}
	if mask == Wildcard || mask == "qemu-append" {
		v.QemuAppend = nil
	}
	if mask == Wildcard || mask == "qemu-override" {
		v.QemuOverride = QemuOverrides{}
	}
}

func (v *KVMConfig) WriteConfig(w io.Writer) error {
	if v.QemuPath != "kvm" {
		fmt.Fprintf(w, "vm config qemu %v\n", v.QemuPath)
	}
	if v.KernelPath != "" {
		fmt.Fprintf(w, "vm config kernel %v\n", v.KernelPath)
	}
	if v.InitrdPath != "" {
		fmt.Fprintf(w, "vm config initrd %v\n", v.InitrdPath)
	}
	if v.CdromPath != "" {
		fmt.Fprintf(w, "vm config cdrom %v\n", v.CdromPath)
	}
	if v.MigratePath != "" {
		fmt.Fprintf(w, "vm config migrate %v\n", v.MigratePath)
	}
	if v.CPU != "host" {
		fmt.Fprintf(w, "vm config cpu %v\n", v.CPU)
	}
	if v.Sockets != 0 {
		fmt.Fprintf(w, "vm config sockets %v\n", v.Sockets)
	}
	if v.Cores != 0 {
		fmt.Fprintf(w, "vm config cores %v\n", v.Cores)
	}
	if v.Threads != 0 {
		fmt.Fprintf(w, "vm config threads %v\n", v.Threads)
	}
	if v.Machine != "" {
		fmt.Fprintf(w, "vm config machine %v\n", v.Machine)
	}
	if v.SerialPorts != 0 {
		fmt.Fprintf(w, "vm config serial-ports %v\n", v.SerialPorts)
	}
	if v.VirtioPorts != "" {
		fmt.Fprintf(w, "vm config virtio-ports %v\n", v.VirtioPorts)
	}
	if v.Vga != "std" {
		fmt.Fprintf(w, "vm config vga %v\n", v.Vga)
	}
	if len(v.Append) > 0 {
		fmt.Fprintf(w, "vm config append %v\n", quoteJoin(v.Append, " "))
	}
	if err := v.Disks.WriteConfig(w); err != nil {
		return err
	}
	if v.UsbUseXHCI != true {
		fmt.Fprintf(w, "vm config usb-use-xhci %t\n", v.UsbUseXHCI)
	}
	if v.TpmSocketPath != "" {
		fmt.Fprintf(w, "vm config tpm-socket %v\n", v.TpmSocketPath)
	}
	if v.BidirectionalCopyPaste != false {
		fmt.Fprintf(w, "vm config bidirectional-copy-paste %t\n", v.BidirectionalCopyPaste)
	}
	if len(v.QemuAppend) > 0 {
		fmt.Fprintf(w, "vm config qemu-append %v\n", quoteJoin(v.QemuAppend, " "))
	}
	if err := v.QemuOverride.WriteConfig(w); err != nil {
		return err
	}

	return nil
}

func (v *KVMConfig) ReadConfig(r io.Reader, ns string) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "vm config") {
			continue
		}

		config := strings.Fields(line)[2:]
		field := config[0]

		switch field {
		case "qemu":
			v.QemuPath = config[1]
		case "kernel":
			v.KernelPath = config[1]
		case "initrd":
			v.InitrdPath = config[1]
		case "cdrom":
			v.CdromPath = config[1]
		case "migrate":
			v.MigratePath = config[1]
		case "cpu":
			v.CPU = config[1]
		case "sockets":
			v.Sockets, _ = strconv.ParseUint(config[1], 10, 64)
		case "cores":
			v.Cores, _ = strconv.ParseUint(config[1], 10, 64)
		case "threads":
			v.Threads, _ = strconv.ParseUint(config[1], 10, 64)
		case "machine":
			v.Machine = config[1]
		case "serial-ports":
			v.SerialPorts, _ = strconv.ParseUint(config[1], 10, 64)
		case "virtio-ports":
			v.VirtioPorts = config[1]
		case "vga":
			v.Vga = config[1]
		case "append":
			v.Append = strings.Fields(config[1])
		case "disks":
			v.ReadFieldConfig(strings.NewReader(line), "disks", ns)
		case "usb-use-xhci":
			v.UsbUseXHCI, _ = strconv.ParseBool(config[1])
		case "tpm-socket":
			v.TpmSocketPath = config[1]
		case "bidirectional-copy-paste":
			v.BidirectionalCopyPaste, _ = strconv.ParseBool(config[1])
		case "qemu-append":
			v.QemuAppend = strings.Fields(config[1])
		case "qemu-override":
			v.ReadFieldConfig(strings.NewReader(line), "qemu-override", ns)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
