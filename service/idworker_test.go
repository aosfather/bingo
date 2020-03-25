package service

import "testing"

func TestIdWorker_NextId(t *testing.T) {
	w := CreateWorker(1, 1)
	for i := 0; i < 100; i++ {
		go t.Log(w.NextId())
	}

	t.Log(w.NextIdWithTime())
}
