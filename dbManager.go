package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var levelMapping = map[string]int{
	"EASY":   1,
	"MEDIUM": 2,
	"HARD":   3,
}

type dbManager struct {
	*sql.DB
}

func newDBManager(configData config) (dbManager, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", configData.ID, configData.PW, configData.DB)
	db, err := sql.Open(configData.StoreType, dsn)
	if err != nil {
		return dbManager{}, err
	}

	dbm := dbManager{db}

	// check for connection
	err = dbm.Ping()
	if err != nil {
		return dbManager{}, err
	}

	return dbm, nil
}

func (dbm *dbManager) insertProblem(input []string) error {
	// 0: id, 1: problem name, 2: level, 3: url
	id, err := strconv.Atoi(input[0])
	if err != nil {
		return fmt.Errorf("insertProblem err : %s", err.Error())
	}
	insert, err := dbm.Query("INSERT INTO `leetcode_problem`(id, problem_name, level_id, url) VALUES(?, ?, ?, ?)", id, input[1], levelMapping[input[2]], input[3])

	if err != nil {
		return fmt.Errorf("insertProblem err : %s", err.Error())
	}
	defer insert.Close()

	return nil
}

func (dbm *dbManager) checkExist(problemID int) bool {
	rows, err := dbm.Query("SELECT COUNT(*) FROM `leetcode_problem` WHERE id=?", problemID)
	if err != nil {
		fmt.Printf("checkExist: err %s", err.Error())
		return false
	}

	defer rows.Close()
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			fmt.Printf("checkExist: result err : %s", err.Error())
			return false
		}
	}

	if id != 0 {
		return true
	}
	return false
}
