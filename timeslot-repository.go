package bivrost

import (
	"sync"
)

// TimeSlotRepository is the interface
// to store timeslots of the events
type TimeSlotRepository interface {
	Add(when timeInSeconds, what interface{})
	Get(when timeInSeconds) []interface{}
	Delete(when timeInSeconds)
}

// InMemoryTimeSlotRepository in-memory implementation of the TimeSlotRepository interface
type InMemoryTimeSlotRepository struct {
	mu      sync.RWMutex
	storage map[timeInSeconds][]interface{}
	logger  Logger
}

// NewInMemoryTimeSlotRepository constructor
func NewInMemoryTimeSlotRepository(logger Logger) *InMemoryTimeSlotRepository {
	return &InMemoryTimeSlotRepository{
		mu:      sync.RWMutex{},
		storage: make(map[timeInSeconds][]interface{}),
		logger:  logger,
	}
}

// Add an event data ID to event timeslot
func (i *InMemoryTimeSlotRepository) Add(when timeInSeconds, what interface{}) {
	i.mu.Lock()
	defer i.mu.Unlock()
	_, ok := i.storage[when]
	if !ok {
		i.storage[when] = make([]interface{}, 0)
	}
	i.storage[when] = append(i.storage[when], what)
	i.logger.Debugf("add event to timeslot %v, now there is %v IDs", when, len(i.storage[when]))
}

// Get the events IDs for the special timeslot
func (i *InMemoryTimeSlotRepository) Get(when timeInSeconds) []interface{} {
	i.mu.RLock()
	defer i.mu.RUnlock()
	entities, ok := i.storage[when]
	if !ok {
		i.logger.Debugf("tried to get event for %v time slot, but nothing", when)
		return nil
	}
	i.logger.Debugf("found %v events for the timeslot %v", len(entities), when)
	return entities
}

// Delete timeslot and all the data
func (i *InMemoryTimeSlotRepository) Delete(when timeInSeconds) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if entities, ok := i.storage[when]; ok {
		i.logger.Debugf("delete %v events for the timeslot %v", len(entities), when)
		delete(i.storage, when)
		return
	}
	i.logger.Debugf("tried delete events for the %v timeslot, but nothing", when)
}
