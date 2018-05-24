package bot

import (
"testing"
)

func TestGetUUID(t *testing.T) {
	wx, err := NewWecat()
	if err != nil {
		panic(err)
	}

	wx.Start()
}

