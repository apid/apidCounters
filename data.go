package apidCounters

import (
	"database/sql"
	"github.com/apid/apid-core"
)

const (
	databaseID = "apidCounters"
)

var (
	data apid.DataService
	db   apid.DB
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
	row := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='apidCounters';")
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Errorf("unable to access DB: %v", err)
		return err
	}
	if count == 0 {
		_, err := db.Exec("CREATE TABLE apidCounters (id text, counter integer);")
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

	res, err := db.Exec("update apidCounters set counter = counter + 1 where id = ?;", id)
	var nRows int64
	if err == nil {
		// note: no potential race condition because events are queued
		nRows, err = res.RowsAffected()
		if nRows == 0 {
			res, err = db.Exec("insert into apidCounters values(?, 1);", id)
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
	log.Debugf("get DB counter: %s", id)

	var count int
	err := db.QueryRow("select counter from apidCounters where id = ?;", id).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debugf("counter %s: %v", id, 0)
			return 0, nil
		}
		log.Errorf("unable to query DB: %v", err)
	}

	log.Debugf("counter %s: %v", id, count)

	return count, err
}

// return a map of all the stored counters
func getDBCounters() (result map[string]int, err error) {
	log.Debug("get DB counters")

	var (
		id    string
		count int
		rows  *sql.Rows
	)

	rows, err = db.Query("select id, counter from apidCounters order by id;")
	defer rows.Close()
	if err != nil {
		log.Errorf("unable to query DB: %v", err)
		return
	}
	result = make(map[string]int)
	for rows.Next() {
		rows.Scan(&id, &count)
		result[id] = count
	}

	log.Debugf("counters: %v", result)

	return
}
