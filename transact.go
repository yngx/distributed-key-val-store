package main

import (
	"bufio"
	"fmt"
	"os"
)

type EventType byte

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

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

type FileTransactionLogger struct {
	events       chan<- Event // Write-only channel for sending events
	errors       <-chan error // Read-only channel for receiving errors
	lastSequence uint64       // The last used event sequence number
	file         *os.File     // The location of the transaction log
}

// this is how you inplement an interface
func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileTransactionLogger{file: file}, nil
}

func (l *FileTransactionLogger) Run() {
	events := make(chan Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	// creates a go routine
	go func() {
		for e := range events { // Retrieve the next Event

			l.lastSequence++ // Increment sequence number

			_, err := fmt.Fprintf( // Write the event to the log
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file) // Create a Scanner for l.file
	outEvent := make(chan Event)        // An unbuffered Event channel
	outError := make(chan error, 1)     // A buffered error channel

	go func() {
		var e Event

		defer close(outEvent) // Close the channels when the
		defer close(outError) // goroutine ends

		for scanner.Scan() {
			line := scanner.Text()

			fmt.Sscanf(
				line, "%d\t%d\t%s\t%s",
				&e.Sequence, &e.EventType, &e.Key, &e.Value)

			// Sanity check! Are the sequence numbers in increasing order?
			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastSequence = e.Sequence // Update last used sequence #

			outEvent <- e // Send the event along
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}
