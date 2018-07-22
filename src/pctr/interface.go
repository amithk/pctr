package pctr

// Uses uint64 counters

type CounterInterface interface {

	// Increment the counter value, and return the resultant value
	IncrementValue(incr uint64) (uint64, error)

	// Delete the counter. Once deleted, calls to other APIs will
	// return error.
	DeleteCounter() error

	// Get current value of counter.
	GetValue() (uint64, error)

	// Return true if counter is deleted
	IsDeleted() bool
}
