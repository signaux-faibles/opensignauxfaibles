package sfregexp

import (
	"testing"
)

func TestValidSiret(t *testing.T) {
	type args struct {
		siret string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"cas normal", args{"33516816700021"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidSiret(tt.args.siret); got != tt.want {
				t.Errorf("ValidSiret() = %v, want %v", got, tt.want)
			}
		})
	}
}
