// Copyright 2015-2023 National Technology & Engineering Solutions of Sandia, LLC (NTESS).
// Under the terms of Contract DE-NA0003525 with NTESS, the U.S. Government retains certain
// rights in this software.

package main

import (
	"reflect"
	"testing"
)

func TestParseNetConfig(t *testing.T) {
	type args struct {
		spec string
		nics map[string]bool
	}
	tests := []struct {
		name    string
		args    args
		want    NetConfig
		wantErr bool
	}{
		{
			name: "vlan alias",
			args: args{
				spec: "xyz",
				nics: map[string]bool{},
			},
			want: NetConfig{
				Alias:                  "xyz",
				VLAN:                   0,
				Bridge:                 DefaultBridge,
				Tap:                    "",
				MAC:                    "",
				Driver:                 DefaultKVMDriver,
				IP4:                    "",
				IP6:                    "",
				QinQ:                   false,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "bridge,vlan alias",
			args: args{
				spec: "my_bridge,foo",
				nics: map[string]bool{},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 "my_bridge",
				Tap:                    "",
				MAC:                    "",
				Driver:                 DefaultKVMDriver,
				IP4:                    "",
				IP6:                    "",
				QinQ:                   false,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "vlan alias,mac",
			args: args{
				spec: "foo,de:ad:be:ef:ca:fe",
				nics: map[string]bool{},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 DefaultBridge,
				Tap:                    "",
				MAC:                    "de:ad:be:ef:ca:fe",
				Driver:                 DefaultKVMDriver,
				IP4:                    "",
				IP6:                    "",
				QinQ:                   false,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "vlan alias,driver",
			args: args{
				spec: "foo,virtio-net-pci",
				nics: map[string]bool{
					"virtio-net-pci": true,
				},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 DefaultBridge,
				Tap:                    "",
				MAC:                    "",
				Driver:                 "virtio-net-pci",
				IP4:                    "",
				IP6:                    "",
				QinQ:                   false,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "bridge,vlan alias,mac",
			args: args{
				spec: "my_bridge,foo,de:ad:be:ef:ca:fe",
				nics: map[string]bool{},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 "my_bridge",
				Tap:                    "",
				MAC:                    "de:ad:be:ef:ca:fe",
				Driver:                 DefaultKVMDriver,
				IP4:                    "",
				IP6:                    "",
				QinQ:                   false,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "bridge,vlan alias,driver",
			args: args{
				spec: "my_bridge,foo,virtio-net-pci",
				nics: map[string]bool{
					"virtio-net-pci": true,
				},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 "my_bridge",
				Tap:                    "",
				MAC:                    "",
				Driver:                 "virtio-net-pci",
				IP4:                    "",
				IP6:                    "",
				QinQ:                   false,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "bridge,vlan alias,qinq",
			args: args{
				spec: "my_bridge,foo,qinq",
				nics: map[string]bool{},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 "my_bridge",
				Tap:                    "",
				MAC:                    "",
				Driver:                 DefaultKVMDriver,
				IP4:                    "",
				IP6:                    "",
				QinQ:                   true,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "vlan alias,mac,driver",
			args: args{
				spec: "foo,de:ad:be:ef:ca:fe,virtio-net-pci",
				nics: map[string]bool{
					"virtio-net-pci": true,
				},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 DefaultBridge,
				Tap:                    "",
				MAC:                    "de:ad:be:ef:ca:fe",
				Driver:                 "virtio-net-pci",
				IP4:                    "",
				IP6:                    "",
				QinQ:                   false,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "vlan alias,mac,qinq",
			args: args{
				spec: "foo,de:ad:be:ef:ca:fe,qinq",
				nics: map[string]bool{},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 DefaultBridge,
				Tap:                    "",
				MAC:                    "de:ad:be:ef:ca:fe",
				Driver:                 DefaultKVMDriver,
				IP4:                    "",
				IP6:                    "",
				QinQ:                   true,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "vlan alias,driver,qinq",
			args: args{
				spec: "foo,virtio-net-pci,qinq",
				nics: map[string]bool{
					"virtio-net-pci": true,
				},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 DefaultBridge,
				Tap:                    "",
				MAC:                    "",
				Driver:                 "virtio-net-pci",
				IP4:                    "",
				IP6:                    "",
				QinQ:                   true,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "bridge,vlan alias,mac,driver",
			args: args{
				spec: "my_bridge,foo,de:ad:be:ef:ca:fe,virtio-net-pci",
				nics: map[string]bool{
					"virtio-net-pci": true,
				},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 "my_bridge",
				Tap:                    "",
				MAC:                    "de:ad:be:ef:ca:fe",
				Driver:                 "virtio-net-pci",
				IP4:                    "",
				IP6:                    "",
				QinQ:                   false,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "bridge,vlan alias,mac,qinq",
			args: args{
				spec: "my_bridge,foo,de:ad:be:ef:ca:fe,qinq",
				nics: map[string]bool{},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 "my_bridge",
				Tap:                    "",
				MAC:                    "de:ad:be:ef:ca:fe",
				Driver:                 DefaultKVMDriver,
				IP4:                    "",
				IP6:                    "",
				QinQ:                   true,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "bridge,vlan alias,driver,qinq",
			args: args{
				spec: "my_bridge,foo,virtio-net-pci,qinq",
				nics: map[string]bool{
					"virtio-net-pci": true,
				},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 "my_bridge",
				Tap:                    "",
				MAC:                    "",
				Driver:                 "virtio-net-pci",
				IP4:                    "",
				IP6:                    "",
				QinQ:                   true,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "vlan alias,mac,driver,qinq",
			args: args{
				spec: "foo,de:ad:be:ef:ca:fe,virtio-net-pci,qinq",
				nics: map[string]bool{
					"virtio-net-pci": true,
				},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 DefaultBridge,
				Tap:                    "",
				MAC:                    "de:ad:be:ef:ca:fe",
				Driver:                 "virtio-net-pci",
				IP4:                    "",
				IP6:                    "",
				QinQ:                   true,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "wifi,x coordinate,y coordinate,z coordinate",
			args: args{
				spec: "wifi,1,2,3",
				nics: map[string]bool{},
			},
			want: NetConfig{
				Alias:                  "",
				VLAN:                   0,
				Bridge:                 DefaultBridge,
				Tap:                    "",
				MAC:                    "",
				Driver:                 DefaultKVMDriver,
				IP4:                    "",
				IP6:                    "",
				QinQ:                   false,
				Wifi:                   true,
				WifiStationCoordinates: WifiStationCoordinates{X: 1, Y: 2, Z: 3},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
		{
			name: "bridge,vlan alias,mac,driver,qinq",
			args: args{
				spec: "my_bridge,foo,de:ad:be:ef:ca:fe,virtio-net-pci,qinq",
				nics: map[string]bool{
					"virtio-net-pci": true,
				},
			},
			want: NetConfig{
				Alias:                  "foo",
				VLAN:                   0,
				Bridge:                 "my_bridge",
				Tap:                    "",
				MAC:                    "de:ad:be:ef:ca:fe",
				Driver:                 "virtio-net-pci",
				IP4:                    "",
				IP6:                    "",
				QinQ:                   true,
				Wifi:                   false,
				WifiStationCoordinates: WifiStationCoordinates{X: 0, Y: 0, Z: 0},
				RxRate:                 0,
				TxRate:                 0,
				Raw:                    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseNetConfig(tt.args.spec, tt.args.nics)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(*parsed, tt.want) {
				t.Errorf("ParseNetConfig() = %v, want %v", parsed, tt.want)
			}

			serialized := parsed.String()

			if serialized != tt.args.spec {
				t.Errorf("NetConfig.String() = %v, want %v", serialized, tt.args.spec)
			}
		})
	}
}

func TestParseBondConfig(t *testing.T) {
	examples := []string{
		"0,1,active-backup",
		"0,1,active-backup,foo-bond",
		"0,1,active-backup,qinq",
		"1,3,balance-tcp,qinq,foo-bond",
		"1,3,balance-tcp,active,no-lacp-fallback",
		"1,3,balance-tcp,active,no-lacp-fallback,qinq",
		"1,3,balance-tcp,active,no-lacp-fallback,foo-bond",
		"1,3,balance-tcp,active,no-lacp-fallback,qinq,foo-bond",
	}

	for _, s := range examples {
		r, err := ParseBondConfig(s)
		if err != nil {
			t.Fatalf("unable to parse `%v`: %v", s, err)
		}

		got := r.String()
		if got != s {
			t.Fatalf("unequal: `%v` != `%v`", s, got)
		}
	}
}
