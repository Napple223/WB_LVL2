package main

import (
	"reflect"
	"testing"
)

func Test_anagramSearch(t *testing.T) {
	type args struct {
		inputArray *[]string
	}
	tests := []struct {
		name string
		args args
		want *map[string]*[]string
	}{
		{
			name: "empty slice",
			args: args{
				inputArray: &[]string{},
			},
			want: &map[string]*[]string{},
		},
		{
			name: "1 element set",
			args: args{
				inputArray: &[]string{"булка", "мужик"},
			},
			want: &map[string]*[]string{},
		},
		{
			name: "Upper register",
			args: args{
				inputArray: &[]string{"ПЯТКА", "ТяПкА"},
			},
			want: &map[string]*[]string{
				"пятка": {"тяпка"},
			},
		},
		{
			name: "Should pass",
			args: args{
				inputArray: &[]string{"пятка", "пятак", "ТяПкА", "булка", "листок", "слиток", "столик"},
			},
			want: &map[string]*[]string{
				"пятка":  {"пятак", "тяпка"},
				"листок": {"слиток", "столик"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := anagramSearch(tt.args.inputArray); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("anagramSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
