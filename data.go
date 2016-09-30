package apidCounters

import (
	"database/sql"
	"github.com/30x/apid"
)

const (
	databaseID = "apidExamplePlugin"
)

var (
	data apid.DataService
	db   *sql.DB
)

// initializes the data source
func initDB(services apid.Services) error {

	data = services.Data()

	// retrieve a named database (`data.DB()` would retrieve a common DB)
	var err error
	db, err = data.DBForID(databaseID)
	if err != nil {
		log.Errorf("unable to access DB: %v", err)
		return err
	}

	// if table doesn't already exist, create it
	row := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='apidExamplePlugin';")
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Errorf("unable to access DB: %v", err)
		return err
	}
	if count == 0 {
		_, err := db.Exec("CREATE TABLE apidExamplePlugin (id text, counter integer);")
		if err != nil {
			log.Errorf("unable to query DB: %v", err)
			return err
		}
	}

	return nil
}

// increment the counter for a given id
func incrementDBCounter(id string) error {
	log.Debugf("increment DB counter for %s", id)

	res, err := db.Exec("update apidExamplePlugin set counter = counter + 1 where id = ?;", id)
	var nRows int64
	if err == nil {
		// note: no potential race condition because events are queued
		nRows, err = res.RowsAffected()
		if nRows == 0 {
			res, err = db.Exec("insert into apidExamplePlugin values(?, 1);", id)
			if err == nil {
				nRows, err = res.RowsAffected()
			}
		}
	}
	if err != nil || nRows == 0 {
		log.Errorf("unable to increment counter in DB: %v", err)
		return err
	}

	return nil
}

// return the counter for a given id
func getDBCounter(id string) (int, error) {
	log.Debug("get DB counter")

	var count int
	err := db.QueryRow("select counter from apidExamplePlugin where id = ?;", id).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		log.Errorf("unable to query DB: %v", err)
	}

	return count, err
}

// return a map of all the stored counters
func getDBCounters() (result map[string]int, err error) {
	log.Debug("get DB counter")

	var (
		id    string
		count int
		rows  *sql.Rows
	)

	rows, err = db.Query("select id, counter from apidExamplePlugin order by id;")
	if err != nil {
		log.Errorf("unable to query DB: %v", err)
		return
	}
	result = make(map[string]int)
	for rows.Next() {
		rows.Scan(&id, &count)
		result[id] = count
	}

	return
}
