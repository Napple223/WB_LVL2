package main

import (
	"reflect"
	"testing"
)

func Test_flags_parseFFlag(t *testing.T) {
	type args struct {
		strF string
	}
	tests := []struct {
		name       string
		f          *flags
		args       args
		wantErr    bool
		wantFField []int
	}{
		{
			name: "no flags",
			f:    newFlags(),
			args: args{
				"",
			},
			wantErr:    false,
			wantFField: []int{},
		},
		{
			name: "1 arg",
			f:    newFlags(),
			args: args{
				"1",
			},
			wantErr:    false,
			wantFField: []int{1},
		},
		{
			name: "fail parsing",
			f:    newFlags(),
			args: args{
				"b",
			},
			wantErr:    true,
			wantFField: []int{0},
		},
		{
			name: "2 args through ,",
			f:    newFlags(),
			args: args{
				"1,2",
			},
			wantErr:    false,
			wantFField: []int{1, 2},
		},
		{
			name: "args through -",
			f:    newFlags(),
			args: args{
				"0-3",
			},
			wantErr:    false,
			wantFField: []int{0, 1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.parseFFlag(tt.args.strF); (err != nil) != tt.wantErr {
				t.Errorf("flags.parseFFlag() error = %v, wantErr %v", err, tt.wantErr)
			}
			for i, v := range tt.f.fields {
				if v != tt.wantFField[i] {
					t.Errorf("got f fields: %v, want: %v", tt.f.fields, tt.wantFField)
				}
			}
		})
	}
}

func Test_flags_cut(t *testing.T) {
	type args struct {
		data []string
	}
	tests := []struct {
		name string
		f    *flags
		args args
		want [][]string
	}{
		{
			name: "default delimeter",
			f: &flags{
				fields:    []int{},
				delimiter: "\t",
			},
			args: args{
				[]string{"test\ttest", "test1\ttest1"},
			},
			want: [][]string{{"test", "test"}, {"test1", "test1"}},
		},
		{
			name: "delimeter :",
			f: &flags{
				fields:    []int{},
				delimiter: ":",
			},
			args: args{
				[]string{"test:test", "test1:test1"},
			},
			want: [][]string{{"test", "test"}, {"test1", "test1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.cut(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flags.cut() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flags_print(t *testing.T) {
	type args struct {
		rows [][]string
	}
	tests := []struct {
		name string
		f    *flags
		args args
		want string
	}{
		{
			name: "separated",
			f: &flags{
				fields:    []int{},
				separated: true,
				delimiter: "\t",
			},
			args: args{
				[][]string{{"test", "test"}, {"test\ttest"}},
			},
			want: "\ntest\ttest",
		},
		{
			name: "no fields flag",
			f: &flags{
				fields:    []int{},
				delimiter: "\t",
			},
			args: args{
				[][]string{{"test", "test"}, {"test1"}},
			},
			want: "\ntest\ttest\ntest1",
		},
		{
			name: "fields flag",
			f: &flags{
				fields:    []int{0, 2},
				delimiter: "\t",
			},
			args: args{
				[][]string{{"test", "wrong", "test"}, {"test1", "wrong", "test1"}},
			},
			want: "\ntest\ttest\t\ntest1\ttest1\t",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.print(tt.args.rows); got != tt.want {
				t.Errorf("flags.print() = %v\nwant %v", got, tt.want)
			}
		})
	}
}
