package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckStr(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int64
		err  error
	}{
		{
			name: "Not valid number",
			args: args{
				str: "4561261212345464",
			},
			want: -1,
			err:  ErrNotValid,
		},
		{
			name: "Valid number #1",
			args: args{
				str: "79927398713",
			},
			want: 79927398713,
		},
		{
			name: "Valid number #2",
			args: args{
				str: "4561261212345467",
			},
			want: 4561261212345467,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckStr(tt.args.str)

			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
