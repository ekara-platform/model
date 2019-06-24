package model

import (
	"reflect"
	"testing"
)

func TestEnvVars_inherit(t *testing.T) {
	type args struct {
		parent map[string]string
	}
	tests := []struct {
		name    string
		r       EnvVars
		args    args
		want    EnvVars
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.inherit(tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvVars.inherit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvVars.inherit() = %v, want %v", got, tt.want)
			}
		})
	}
}
