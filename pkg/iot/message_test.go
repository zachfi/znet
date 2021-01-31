package iot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZigbee(t *testing.T) {

	cases := []struct {
		Message []byte
		Obj     ZigbeeBridgeLog
	}{
		{
			[]byte(`{"type":"device_announced","message":"announce","meta":{"friendly_name":"0x0017880104650857"}}`),
			ZigbeeBridgeLog{
				Type:    "device_announced",
				Message: "announce",
				Meta: map[string]interface{}{
					"friendly_name": "0x0017880104650857",
				},
			},
		},
	}

	for _, tc := range cases {
		obj := ZigbeeBridgeLog{}
		err := json.Unmarshal(tc.Message, &obj)
		require.NoError(t, err)
		require.Equal(t, tc.Obj, obj)
	}

}
