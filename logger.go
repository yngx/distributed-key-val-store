package main

type EventType byte

// iota is a predefined value that can be used in a constant declaration.
// When used, it represents successive untyped integer constants that can be used
// to construct a set of related constants. Its value restarts
// at zero in each constant declaration and increments with each constant assignment.
// It also allows implicit reptition.
const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type Event struct {
	Sequence  uint64    // A unique record ID
	EventType EventType // The action taken
	Key       string    // The key affected by this transaction
	Value     string    // The value of a PUT the transaction
}

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}
