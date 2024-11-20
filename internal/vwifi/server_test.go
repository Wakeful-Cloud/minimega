// Copyright 2024 Colorado School of Mines CSCI 370 FA24 NREL 2 Group

package vwifi

import (
	"reflect"
	"testing"
)

func TestServer_parseStatus(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		want    *ServerStatus
		wantErr bool
	}{
		{
			name: "parses status with packet loss disabled",
			args: args{
				raw: "CLIENT VHOST : Listener on port : 8211\nCLIENT TCP : Listener on port : 8212\nSPY : Listener on port : 8213\nCTRL : Listener on port : 8214\nSize of disconnected : 15\nPacket loss : disable\nScale : 1",
			},
			want: &ServerStatus{
				VsockPort:          8211,
				TcpPort:            8212,
				SpyPort:            8213,
				ControlPort:        8214,
				SizeOfDisconnected: 15,
				PacketLoss:         false,
				Scale:              1,
			},
			wantErr: false,
		},
		{
			name: "parses status with packet loss enabled",
			args: args{
				raw: "CLIENT VHOST : Listener on port : 8211\nCLIENT TCP : Listener on port : 8212\nSPY : Listener on port : 8213\nCTRL : Listener on port : 8214\nSize of disconnected : 15\nPacket loss : enable\nScale : 1",
			},
			want: &ServerStatus{
				VsockPort:          8211,
				TcpPort:            8212,
				SpyPort:            8213,
				ControlPort:        8214,
				SizeOfDisconnected: 15,
				PacketLoss:         true,
				Scale:              1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{}
			got, err := server.parseStatus(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.parseStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.parseStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
