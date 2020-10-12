package bivrost

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestService_Serve(t *testing.T) {
	type testEvent struct {
		after int64
		desc  string
	}
	events := []testEvent{
		{
			after: 1,
			desc:  "add1simple",
		},
		{
			after: 2,
			desc:  "add2simple",
		},
		{
			after: 4,
			desc:  "add4simple",
		},
		{
			after: 1,
			desc:  "add1oneMoreTime",
		},
		{
			after: 4,
			desc:  "add4oneMoreTime",
		},
		{
			after: 3,
			desc:  "sleep1secondAndAdd3",
		},
		{
			after: 1,
			desc:  "sleep1secondAndAdd1",
		},
	}

	s, in, out := Init()

	go s.Serve()
	defer s.Cancel()

	go func() {
		for i, event := range events {
			if i > 4 {
				time.Sleep(1 * time.Second)
			}
			in <- &Event{
				When:   time.Now().Add(time.Duration(event.after) * time.Second),
				Entity: fmt.Sprintf("%v_%v_", time.Now().Add(time.Duration(event.after)*time.Second).Unix(), time.Now().Unix()) + event.desc,
			}
		}
	}()
	go func() {
		time.Sleep(5 * time.Second)
		s.Cancel()
	}()
	for entity := range out {
		eventTime := strings.Split(entity.(string), "_")[0]
		currentTime := fmt.Sprintf("%v", time.Now().Unix())
		if eventTime != currentTime {
			t.Errorf("incorrect result, event time %v current time %v, for event data %v", eventTime, currentTime, entity.(string))
		}
	}

}
