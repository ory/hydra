package connectdb

import (
	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strings"
)

func ConnectToMySQL(path string) (*sqlx.DB, error) {
	parts := strings.Split(path, "/")
	database := parts[len(parts)-1]
	if database == "mysql" {
		return sqlx.Connect("mysql", path)
	}

	db, err := sqlx.Connect("mysql", strings.Join(parts[:len(parts)-1], "/") + "/mysql")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database)); err != nil {
		return nil, errors.Wrap(err, "Unable to create database")
	} else if _, err = db.Exec(fmt.Sprintf("USE %s", database)); err != nil {
		return nil, errors.Wrap(err, "Unable to use database")
	}

	return db, err
}
