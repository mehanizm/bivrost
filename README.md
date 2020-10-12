# `BIVROST`

> Send messages to channel with specified time to receive

[![GoDoc](https://godoc.org/github.com/mehanizm/bivrost?status.svg)](https://pkg.go.dev/github.com/mehanizm/bivrost)
![Go](https://github.com/mehanizm/bivrost/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/mehanizm/bivrost/branch/main/graph/badge.svg)](https://codecov.io/gh/mehanizm/bivrost)
[![Go Report](https://goreportcard.com/badge/github.com/mehanizm/bivrost)](https://goreportcard.com/report/github.com/mehanizm/bivrost)

See an example in cmd/main.go

```go
func main() {

	log.Println("start")

	s, in, out := bivrost.Init()
	go s.Serve()

	var events = []event{
		{1, "one second"},
		{2, "two seconds"},
		{3, "three seconds"},
	}

	go func() {
		for _, event := range events {
			in <- &bivrost.Event{
				When:   time.Now().Add(time.Duration(event.after) * time.Second),
				Entity: interface{}(event.message),
			}
		}
	}()

	go func() {
		for event := range out {
			log.Println(event)
		}
	}()

	time.Sleep(4 * time.Second)
	s.Cancel()
	time.Sleep(1 * time.Millisecond)
}
```

OUTPUT:
```shell
>> 2020/10/12 13:08:31 start
>> 2020/10/12 13:08:32 one second
>> 2020/10/12 13:08:33 two seconds
>> 2020/10/12 13:08:34 three seconds
>> 2020/10/12 13:08:35 bivrost [bivrost] INFO: cancel serving
```