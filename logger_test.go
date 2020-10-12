package bivrost

import (
	"testing"
)

func TestDefaultLog_Errorf(t *testing.T) {
	tests := []struct {
		name   string
		logger Logger
	}{
		{
			name:   "debug",
			logger: DefaultLogger(DEBUG, "test debug"),
		},
		{
			name:   "info",
			logger: DefaultLogger(INFO, "test info"),
		},
		{
			name:   "warning",
			logger: DefaultLogger(WARNING, "test warning"),
		},
		{
			name:   "error",
			logger: DefaultLogger(ERROR, "test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logger.Debugf("message %v", 1)
			tt.logger.Infof("message %v", 2)
			tt.logger.Warningf("message %v", 3)
			tt.logger.Errorf("message %v", 4)
		})
	}
}
