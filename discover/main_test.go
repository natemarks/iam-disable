package discover

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestUsers(t *testing.T) {
	type args struct {
		log *zerolog.Logger
	}
	tests := []struct {
		name         string
		args         args
		minUserCount int
	}{
		{name: "sdf", args: args{log: &zerolog.Logger{}}, minUserCount: 60},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUsers := Users(tt.args.log); len(gotUsers) < tt.minUserCount {
				t.Errorf("Users() = %v, want %v", gotUsers, tt.minUserCount)
			}
		})
	}
}
