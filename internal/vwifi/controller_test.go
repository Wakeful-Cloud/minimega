// Copyright 2024 Colorado School of Mines CSCI 370 FA24 NREL 2 Group

package vwifi

import (
	"net"
	"reflect"
	"testing"
)

func TestController_parseClients(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		want    []ControllerClient
		wantErr bool
	}{
		{
			name: "ignores junk lines",
			args: args{
				raw: "wjfoiwejfpowejofpiwjeoif\nokjsopjweopijwef\nowijefowiejfowjeoifj\n",
			},
			want:    []ControllerClient{},
			wantErr: false,
		},
		{
			name: "parses anonymous spy client",
			args: args{
				raw: "S:1",
			},
			want: []ControllerClient{
				{
					Spy:  true,
					CID:  1,
					Name: "",
					X:    0,
					Y:    0,
					Z:    0,
				},
			},
			wantErr: false,
		},
		{
			name: "parses anonymous regular client",
			args: args{
				raw: "2 5 6 7",
			},
			want: []ControllerClient{
				{
					Spy:  false,
					CID:  2,
					Name: "",
					X:    5,
					Y:    6,
					Z:    7,
				},
			},
			wantErr: false,
		},
		{
			name: "parses named spy client",
			args: args{
				raw: "S:1 (my spy name)",
			},
			want: []ControllerClient{
				{
					Spy:  true,
					CID:  1,
					Name: "my spy name",
					X:    0,
					Y:    0,
					Z:    0,
				},
			},
			wantErr: false,
		},
		{
			name: "parses named regular client",
			args: args{
				raw: "2 (my regular name) 5 6 7",
			},
			want: []ControllerClient{
				{
					Spy:  false,
					CID:  2,
					Name: "my regular name",
					X:    5,
					Y:    6,
					Z:    7,
				},
			},
			wantErr: false,
		},
		{
			name: "parses mixed clients",
			args: args{
				raw: "S:1 (my spy name)\n2 (my regular name) 5 6 7",
			},
			want: []ControllerClient{
				{
					Spy:  true,
					CID:  1,
					Name: "my spy name",
					X:    0,
					Y:    0,
					Z:    0,
				},
				{
					Spy:  false,
					CID:  2,
					Name: "my regular name",
					X:    5,
					Y:    6,
					Z:    7,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &Controller{}
			got, err := controller.parseClients(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("Controller.parseClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Controller.parseClients() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_parseStatus(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		want    *ControllerServerStatus
		wantErr bool
	}{
		{
			name: "parses status with packet loss disabled and disconnected spy",
			args: args{
				raw: "CTRL : IP : 127.0.0.1\nCTRL : Port : 8214\nSRV : PacketLoss : Disable\nSRV : Scale : 1\nSRV VHOST : Port : 8211\nSRV INET : Port : 8212\nSRV : SizeOfDisconnected : 15\nSPY : Disconnected",
			},
			want: &ControllerServerStatus{
				ControlIP:          net.IPv4(127, 0, 0, 1),
				ControlPort:        8214,
				PacketLoss:         false,
				Scale:              1,
				VsockPort:          8211,
				TcpPort:            8212,
				SizeOfDisconnected: 15,
				SpyConnected:       false,
			},
			wantErr: false,
		},
		{
			name: "parses status with packet loss enabled and disconnected spy",
			args: args{
				raw: "CTRL : IP : 127.0.0.1\nCTRL : Port : 8214\nSRV : PacketLoss : Enable\nSRV : Scale : 1\nSRV VHOST : Port : 8211\nSRV INET : Port : 8212\nSRV : SizeOfDisconnected : 15\nSPY : Disconnected",
			},
			want: &ControllerServerStatus{
				ControlIP:          net.IPv4(127, 0, 0, 1),
				ControlPort:        8214,
				PacketLoss:         true,
				Scale:              1,
				VsockPort:          8211,
				TcpPort:            8212,
				SizeOfDisconnected: 15,
				SpyConnected:       false,
			},
			wantErr: false,
		},
		{
			name: "parses status with packet loss disabled and connected spy",
			args: args{
				raw: "CTRL : IP : 127.0.0.1\nCTRL : Port : 8214\nSRV : PacketLoss : Disable\nSRV : Scale : 1\nSRV VHOST : Port : 8211\nSRV INET : Port : 8212\nSRV : SizeOfDisconnected : 15\nSPY : Connected",
			},
			want: &ControllerServerStatus{
				ControlIP:          net.IPv4(127, 0, 0, 1),
				ControlPort:        8214,
				PacketLoss:         false,
				Scale:              1,
				VsockPort:          8211,
				TcpPort:            8212,
				SizeOfDisconnected: 15,
				SpyConnected:       true,
			},
			wantErr: false,
		},
		{
			name: "parses status with packet loss enabled and connected spy",
			args: args{
				raw: "CTRL : IP : 127.0.0.1\nCTRL : Port : 8214\nSRV : PacketLoss : Enable\nSRV : Scale : 1\nSRV VHOST : Port : 8211\nSRV INET : Port : 8212\nSRV : SizeOfDisconnected : 15\nSPY : Connected",
			},
			want: &ControllerServerStatus{
				ControlIP:          net.IPv4(127, 0, 0, 1),
				ControlPort:        8214,
				PacketLoss:         true,
				Scale:              1,
				VsockPort:          8211,
				TcpPort:            8212,
				SizeOfDisconnected: 15,
				SpyConnected:       true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &Controller{}
			got, err := controller.parseStatus(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("Controller.parseStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Controller.parseStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
