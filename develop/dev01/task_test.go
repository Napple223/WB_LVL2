package main

import (
	"fmt"
	"testing"
)

func Test_getNTPTime(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "wrong host",
			args: args{
				address: "someRandomHostName",
			},
			wantErr: true,
		},
		{
			name: "correct host",
			args: args{
				address: "0.beevik-ntp.pool.ntp.org",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNTPTime(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNTPTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
