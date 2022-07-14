package csvutil

import (
	"testing"
	"time"
)

func Test_getDateFormat(t *testing.T) {
	type args struct {
		fmtString string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "YYYY-MM-DD",
			args: args{fmtString: "DATE(YYYY-MM-DD)"},
			want: "2006-01-02",
		},
		{
			name: "MM/DD/YYYY",
			args: args{fmtString: "DATE(MM/DD/YYYY)"},
			want: "01/02/2006",
		},
		{
			name: "Empty",
			args: args{fmtString: "DATE"},
			want: "2006-01-02",
		},
		{
			name: "Error",
			args: args{fmtString: "DATE(YMD)"},
			want: "YMD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDateFormat(tt.args.fmtString); got != tt.want {
				t.Errorf("getDateFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertDate(t *testing.T) {
	type args struct {
		field   string
		dateFmt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "YYYY-MM-DD",
			args: args{
				field:   "2021-01-02",
				dateFmt: "2006-01-02",
			},
			want: time.Date(2021, 01, 02, 0, 0, 0, 0, time.Local).String(),
		},
		{
			name: "YYYY/MM/DD",
			args: args{
				field:   "2021/01/02",
				dateFmt: "2006/01/02",
			},
			want: time.Date(2021, 01, 02, 0, 0, 0, 0, time.Local).String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertDate(tt.args.field, tt.args.dateFmt); got != tt.want {
				t.Errorf("convertDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
