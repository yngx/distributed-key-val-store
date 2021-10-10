import (
	"database/sql"
	_ "github.com/lib/pq" // Anonymously import the driver package
)

type PostgresTransactionLogger struct {
	events chan<- Event
	errors <-chan error
	db     *sql.DB // The database access interface
}

func (l *PostgresTransactionLogger) WritePut(key, value: string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func NewPostgresTransactionLogger(host, dbName, user, password string ) (TransactionLogger, error) {
	// fill out
}

func (l *PostgresTransactionLogger) Run() {
	// fill out
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan Error) {
	// fill out
}