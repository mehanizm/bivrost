package bivrost

import (
	"sync"
	"time"
)

const (
	inChanBuffer  = 100
	outChanBuffer = 100
)

// Event struct holds info about event
// when and which data it contains
type Event struct {
	When   time.Time
	Entity interface{}
}

// Service main bivrost service struct
type Service struct {
	// InChan is exported to client and used
	// to send events: entity with when in unix timestamp in seconds
	InChan chan *Event
	// OutChan is exported to client and used
	// to receive entities when the time has come
	OutChan chan interface{}

	// mutex to sync isServing flag
	mu sync.RWMutex
	// UTC location info
	loc *time.Location
	// queue to store sorted timestamps and get the min
	eventQueue EventQueue
	// timeslot repo to get entity IDs from entity repo by timestamp
	timeSlotRepo TimeSlotRepository
	// newEvent is the signal channel to process new event
	newEvent chan timeInSeconds
	// cancel is the signal channel to process cancel event
	cancel chan struct{}
	// isServing – flag that service is serving now
	isServing bool
	// logger interface to use external loggers
	logger Logger
}

// Init default bivrost service
func Init() (*Service, chan<- *Event, <-chan interface{}) {
	logLevel := INFO
	s := &Service{
		InChan:       make(chan *Event, inChanBuffer),
		OutChan:      make(chan interface{}, outChanBuffer),
		mu:           sync.RWMutex{},
		eventQueue:   NewHeapQueue([]timeInSeconds{}, DefaultLogger(logLevel, "event-queue")),
		timeSlotRepo: NewInMemoryTimeSlotRepository(DefaultLogger(logLevel, "timeslot-repo")),
		newEvent:     make(chan timeInSeconds),
		cancel:       make(chan struct{}),
		isServing:    false,
		logger:       DefaultLogger(logLevel, "bivrost"),
	}
	return s, s.InChan, s.OutChan
}

func (s *Service) WithEventQueue(eq EventQueue) *Service {
	s.eventQueue = eq
	return s
}

func (s *Service) WithTimeSlotRepository(ts TimeSlotRepository) *Service {
	s.timeSlotRepo = ts
	return s
}

func (s *Service) WithLogger(logger Logger) *Service {
	s.logger = logger
	return s
}

// addEvent internal method to add new event
// we add event in two step transaction:
// * add the event data to timeslot repo
// * add timestamp to the eventQueue
func (s *Service) addEvent(event *Event) {
	when := timeInSeconds(event.When.UTC().Unix())
	s.logger.Debugf("adding event after %v", when)
	if when-timeInSeconds(time.Now().UTC().Unix()) < 0 {
		s.logger.Errorf("new event %v is in the past, skip", event.When)
		return
	}
	s.timeSlotRepo.Add(when, event.Entity)
	s.eventQueue.Add(when)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isServing {
		s.newEvent <- when
	}
}

// Cancel serving
func (s *Service) Cancel() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isServing {
		s.cancel <- struct{}{}
	}
}

// Serve main function
func (s *Service) Serve() {

	var nextEvent timeInSeconds
	var timer *time.Timer

	// get new events from chan and add it
	go func() {
		for event := range s.InChan {
			s.addEvent(event)
		}
	}()

	nextEvent = s.eventQueue.Next()

	// start serving
	s.mu.Lock()
	s.isServing = true
	s.mu.Unlock()
	timer = time.NewTimer(time.Duration(int64(nextEvent)-time.Now().UTC().Unix()) * time.Second)
	for {
		select {

		// if the timer went off we work on corresponding entities
		case <-timer.C:
			entities := s.timeSlotRepo.Get(nextEvent)
			go func() {
				for _, entity := range entities {
					s.mu.Lock()
					s.logger.Debugf("send to out chan entity %v", entity)
					s.OutChan <- entity
					s.mu.Unlock()
				}
			}()
			nextEvent = s.eventQueue.Next()
			timer = time.NewTimer(time.Duration(int64(nextEvent)-time.Now().UTC().Unix()) * time.Second)

		// if there is new event we check for the needed changing of the next event
		case newEvent := <-s.newEvent:
			if newEvent <= nextEvent {
				s.logger.Debugf("add new event to queue with changing next %v", newEvent)
				s.eventQueue.Add(nextEvent)
				nextEvent = s.eventQueue.Next()
				if !timer.Stop() {
					<-timer.C
				}
				timer = time.NewTimer(time.Duration(int64(nextEvent)-time.Now().UTC().Unix()) * time.Second)
			} else {
				s.logger.Debugf("new event added without changing next %v", newEvent)
			}

		// if there is a cancel event we clean up and go away
		case <-s.cancel:
			s.logger.Infof("cancel serving")
			if nextEvent != onceUponATime {
				s.eventQueue.Add(nextEvent)
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			s.isServing = false
			if !timer.Stop() {
				<-timer.C
			}
			close(s.OutChan)
			return
		}
	}
}
