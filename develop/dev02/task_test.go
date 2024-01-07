package main

import "testing"

func Test_unzipString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "len == 0",
			args: args{
				input: "",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "only digits in string",
			args: args{
				input: "45",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "common case",
			args: args{
				input: "a4bc2d5e",
			},
			want:    "aaaabccddddde",
			wantErr: false,
		},
		{
			name: "string without changes",
			args: args{
				input: "abcd",
			},
			want:    "abcd",
			wantErr: false,
		},
		{
			name: "simple escape sequence",
			args: args{
				input: `qwe\4\5`,
			},
			want:    "qwe45",
			wantErr: false,
		},
		{
			name: "escape sequence digit as rune",
			args: args{
				input: `qwe\45`,
			},
			want:    "qwe44444",
			wantErr: false,
		},
		{
			name: "escape sequence \\ as string",
			args: args{
				input: `qwe\\5`,
			},
			want:    "qwe\\\\\\\\\\",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unzipString(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("unzipString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unzipString() = %v, want %v", got, tt.want)
			}
		})
	}
}
