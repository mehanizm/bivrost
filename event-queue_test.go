package bivrost

import (
	"sync"
	"testing"
)

func TestTimeOrderedQueue(t *testing.T) {
	initSl := []timeInSeconds{3, 5}
	toq := NewHeapQueue(initSl, DefaultLogger(DEBUG, "event-queue"))
	toq.Add(0)
	toq.Add(6)
	toq.Add(0)
	toq.Add(onceUponATime)
	m := toq.Next()
	if m != 0 {
		t.Fatalf("wrong max %v, %v", m, 0)
	}
	var m1, m2 timeInSeconds
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		m1 = toq.Next()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		m2 = toq.Next()
	}()
	wg.Wait()
	if m1 != 3 && m1 != 5 {
		t.Fatalf("wrong min %v, %v or %v", m1, 3, 5)
	}
	if m2 != 3 && m2 != 5 {
		t.Fatalf("wrong min %v, %v or %v", m2, 3, 5)
	}
	m = toq.Next()
	if m != 6 {
		t.Fatalf("wrong min %v, %v", m, 6)
	}
	m = toq.Next()
	if m != onceUponATime {
		t.Fatalf("should be onceUponATime")
	}
}
