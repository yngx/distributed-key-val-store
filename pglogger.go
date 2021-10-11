import (
	"database/sql"
	_ "github.com/lib/pq" // Anonymously import the driver package
)

type PostgresDBParams struct {
    dbName   string
    host     string
    user     string
    password string
}

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

func NewPostgresTransactionLogger(config PostgresDBParams) (TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
        config.host, config.dbName, config.user, config.password)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open db: %w", err)
    }

	err = db.Ping()                 // Test the database connection
    if err != nil {
        return nil, fmt.Errorf("failed to open db connection: %w", err)
    }

    logger := &PostgresTransactionLogger{db: db}

	// implement verify table exists
    exists, err := logger.verifyTableExists()
    if err != nil {
        return nil, fmt.Errorf("failed to verify table exists: %w", err)
    }

	if !exists {
		// implement create Table
        if err = logger.createTable(); err != nil {
            return nil, fmt.Errorf("failed to create table: %w", err)
        }
    }

    return logger, nil
}

func (l *PostgresTransactionLogger) Run() {
	// fill out
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan Error) {
	// fill out
}