package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro"
	"time"

	// Blank import for postgresql driver only
	_ "github.com/lib/pq"
	"sync"
)

const (
	createTableInterval string = `CREATE TABLE IF NOT EXISTS "interval" (
"id" SERIAL,
"start_time" TIMESTAMP WITH TIME ZONE NOT NULL,
"planned_duration" BIGINT DEFAULT 0,
"actual_duration" BIGINT DEFAULT 0,
"category" TEXT NOT NULL,
"state" INTEGER DEFAULT 1,
PRIMARY KEY("id")
);`
)

type dbRepo struct {
	db *sql.DB
	sync.RWMutex
}

func NewPostgresRepo(dbName string) (*dbRepo, error) {
	connStr := fmt.Sprintf("user=postgres dbname=%s sslmode=disable", dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(createTableInterval); err != nil {
		return nil, err
	}

	return &dbRepo{
		db: db,
	}, nil
}

func (r *dbRepo) Create(i pomodoro.Interval) (int64, error) {
	// Create entry in the repository
	r.Lock()
	defer r.Unlock()
	// Prepare INSERT statement
	row := r.db.QueryRow(
		`INSERT INTO interval 
(start_time, planned_duration, actual_duration, category, state) 
VALUES($1, $2, $3, $4, $5) 
RETURNING id`,
		i.StartTime, i.PlannedDuration, i.ActualDuration, i.Category, i.State)
	// INSERT results
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *dbRepo) Update(i pomodoro.Interval) error {
	// Update entry in the repository
	r.Lock()
	defer r.Unlock()
	// Prepare UPDATE statement
	updStmt, err := r.db.Prepare(
		"UPDATE interval SET start_time=$1, actual_duration=$2, state=$3 WHERE id=$4")
	if err != nil {
		return err
	}
	defer updStmt.Close()
	// Exec UPDATE statement
	res, err := updStmt.Exec(i.StartTime, i.ActualDuration, i.State, i.ID)
	if err != nil {
		return err
	}
	// UPDATE results
	_, err = res.RowsAffected()
	return err
}

func (r *dbRepo) ByID(id int64) (pomodoro.Interval, error) {
	// Search items in the repository by ID
	r.RLock()
	defer r.RUnlock()
	// Query DB row based on ID
	row := r.db.QueryRow("SELECT * FROM interval WHERE id=$1", id)
	var i pomodoro.Interval
	err := row.Scan(&i.ID, &i.StartTime, &i.PlannedDuration, &i.ActualDuration, &i.Category, &i.State)
	return i, err
}

func (r *dbRepo) Last() (pomodoro.Interval, error) {
	// Search last item in the repository
	r.RLock()
	defer r.RUnlock()
	// Query and parse last row into Interval struct
	var last pomodoro.Interval
	err := r.db.QueryRow("SELECT * FROM interval ORDER BY id desc LIMIT 1").Scan(
		&last.ID, &last.StartTime, &last.PlannedDuration, &last.ActualDuration, &last.Category, &last.State,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return last, pomodoro.ErrNoIntervals
	}
	if err != nil {
		return last, err
	}
	return last, nil
}

func (r *dbRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	// Search last n items of type break in the repository
	r.RLock()
	defer r.RUnlock()
	// Define SELECT query for breaks
	stmt := `SELECT * FROM interval WHERE category LIKE '%Break' ORDER BY id DESC LIMIT $1`
	// Query DB for breaks
	rows, err := r.db.Query(stmt, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Parse data into slice of Interval
	var data []pomodoro.Interval
	for rows.Next() {
		var i pomodoro.Interval
		err = rows.Scan(&i.ID, &i.StartTime, &i.PlannedDuration, &i.ActualDuration, &i.Category, &i.State)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	// Return data
	return data, err
}

func (r *dbRepo) CategorySummary(day time.Time, filter string) (time.Duration, error) {
	// Return a daily summary
	r.RLock()
	defer r.RUnlock()
	// Define SELECT query for daily summary
	stmt := `SELECT sum(actual_duration) FROM interval 
WHERE category LIKE $1 AND 
to_date(start_time::TEXT, 'YYYY-MM-DD HH24:MI:SS')=
to_date($2::TEXT, 'YYYY-MM-DD HH24:MI:SS');`
	var ds sql.NullInt64
	err := r.db.QueryRow(stmt, filter, day).Scan(&ds)
	var d time.Duration
	if ds.Valid {
		d = time.Duration(ds.Int64)
	}
	return d, err
}
