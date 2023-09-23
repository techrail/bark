package logManager

import (
	"reflect"
	"testing"
)

func TestGetLogger(t *testing.T) {
	type args struct {
		serverUrl   string
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want *LoggerSession
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLogger(tt.args.serverUrl, tt.args.serviceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
