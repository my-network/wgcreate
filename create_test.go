package wgcreate

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/device"
)

func TestCreateUserspace(t *testing.T) {
	for i := 0; i < 3; i++ {
		ifaceName, err := Create("utun256", 1000, true, &device.Logger{
			Verbosef: func(format string, args ...interface{}) {},
			Errorf: func(format string, args ...interface{}) {
				t.Errorf(format, args...)
			},
		})
		require.NoError(t, err)
		require.Equal(t, "utun256", ifaceName)
	}
}
