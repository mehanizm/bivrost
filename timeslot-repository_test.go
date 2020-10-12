package bivrost

import (
	"reflect"
	"testing"
)

func TestInMemorySlotStorage(t *testing.T) {
	s := NewInMemoryTimeSlotRepository(DefaultLogger(INFO, "timeslot-repo"))
	s.Add(1, 1)
	s.Add(1, 2)
	s.Add(1, 3)
	s.Add(2, 1)
	s.Add(2, 2)
	s.Delete(1)
	if !reflect.DeepEqual(s.Get(2), []interface{}{1, 2}) {
		t.Fatal("not expected result", s.Get(2))
	}
	if s.Get(1) != nil {
		t.Fatal("not expected result", s.Get(1))
	}
	for i := 0; i < 1_000; i++ {
		go func(i int) {
			s.Add(timeInSeconds(i+10)%100, uint64(i))
		}(i)
		go func(i int) {
			if i%100 == 0 {
				s.Delete(timeInSeconds(i+10) % 100)
			}
		}(i)
	}
}
