package main

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_trimSuffix(t *testing.T) {
	type args struct {
		column string
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantSuffix string
		wantErr    bool
	}{
		{
			name: "OK",
			args: args{
				"1G",
			},
			want:       1,
			wantSuffix: "G",
			wantErr:    false,
		},
		{
			name: "!OK",
			args: args{
				"g",
			},
			want:       -1,
			wantSuffix: "",
			wantErr:    true,
		},
		{
			name: "digits only",
			args: args{
				"28",
			},
			want:       -1,
			wantSuffix: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := trimSuffix(tt.args.column)
			if (err != nil) != tt.wantErr {
				t.Errorf("trimSuffix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("trimSuffix() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantSuffix {
				t.Errorf("trimSuffix() got1 = %v, want %v", got1, tt.wantSuffix)
			}
		})
	}
}

func Test_readFile(t *testing.T) {
	type args struct {
		inputFileName string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				"./input.txt",
			},
			want:    []string{"hello there", "i am somebody"},
			wantErr: false,
		},
		{
			name: "!OK",
			args: args{
				"",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readFile(tt.args.inputFileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("readFile() error\n = %v, \nwantErr %v\n", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Println()
				t.Errorf("\nreadFile() = %v,\nwant %v\n", got, tt.want)
			}
		})
	}
}

func Test_splitData(t *testing.T) {
	type args struct {
		data []string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "OK",
			args: args{
				[]string{"1 2 3", "4 5 6"},
			},
			want: [][]string{{"1", "2", "3"}, {"4", "5", "6"}},
		},
		{
			name: "OK empty",
			args: args{
				[]string{},
			},
			want: [][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flags_columnNumericSort(t *testing.T) {
	type args struct {
		data [][]string
	}
	tests := []struct {
		name    string
		f       *flags
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "column out of range data",
			args: args{
				data: [][]string{{"1"}, {"1"}},
			},
			f: &flags{
				column: 1,
			},
			want:    []string{"1", "1"},
			wantErr: true,
		},
		{
			name: "human read",
			f: &flags{
				humanRead: true,
			},
			args: args{
				data: [][]string{{"1c", "1b"}, {"1a", "3b"}},
			},
			want:    []string{"1a 3b", "1c 1b"},
			wantErr: false,
		},
		{
			name: "reverse",
			f: &flags{
				reverse: true,
			},
			args: args{
				data: [][]string{{"1", "1b"}, {"3", "3b"}},
			},
			want:    []string{"3 3b", "1 1b"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.columnNumericSort(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("flags.columnNumericSort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flags.columnNumericSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
