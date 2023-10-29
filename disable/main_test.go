package disable

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestIamUser(t *testing.T) {
	type args struct {
		username string
		log      *zerolog.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "fake", args: args{username: "notarealuser", log: &zerolog.Logger{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IamUser(tt.args.username, tt.args.log)
		})
	}
}
