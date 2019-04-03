package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
)

const (
	insertQueryShort = "INSERT INTO \"table\""
)

func TestGetStorage(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		log.Panicf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := postgresStorage{db: db}

	if storage.Get() != db {
		t.Fail()
	}
}

func TestTransactionSuccessful(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Panicf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := postgresStorage{db: db}

	sqlMock.ExpectBegin()

	var expectedLastInsertID int64 = 1
	var expectedLastRowsAffected int64 = 1
	expectedResult := sqlmock.NewResult(expectedLastInsertID, expectedLastRowsAffected)
	sqlMock.ExpectExec(insertQueryShort).WithArgs("arg1", "arg2", "arg3").WillReturnResult(expectedResult)

	sqlMock.ExpectCommit()

	txErr := storage.Transaction(context.TODO(), func(context context.Context, tx *sql.Tx) error {
		actualResult, err := tx.Exec(insertQueryShort, "arg1", "arg2", "arg3")
		actualLastInsertID, _ := actualResult.LastInsertId()
		actualRowsAffected, _ := actualResult.RowsAffected()

		if actualLastInsertID != expectedLastInsertID || actualRowsAffected != expectedLastRowsAffected {
			t.Fail()
		}
		return err
	})

	if txErr != nil {
		t.Fail()
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionFailToBeginTransaction(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Panicf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()

	storage := postgresStorage{db: db}
	txErr := storage.Transaction(context.TODO(), func(context context.Context, tx *sql.Tx) error {
		return nil
	})

	if txErr == nil {
		t.Fail()
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionSuccessfulRollback(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Panicf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	storage := postgresStorage{db: db}
	expectedErr := fmt.Errorf("unexpected error")
	txErr := storage.Transaction(context.TODO(), func(context context.Context, tx *sql.Tx) error {
		return expectedErr
	})

	if txErr != expectedErr {
		t.Fail()
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionFailRollback(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Panicf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()

	storage := postgresStorage{db: db}
	expectedErr := fmt.Errorf("unexpected error")
	txErr := storage.Transaction(context.TODO(), func(context context.Context, tx *sql.Tx) error {
		return expectedErr
	})

	if txErr == nil || txErr == expectedErr {
		t.Fail()
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionFailToCommitTransaction(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Panicf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := postgresStorage{db: db}
	txErr := storage.Transaction(context.TODO(), func(context context.Context, tx *sql.Tx) error {
		return nil
	})

	if txErr == nil {
		t.Fail()
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCloseStorage(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		log.Panicf("an error '%s' was not expected when opening a stub database connection", err)
	}

	storage := postgresStorage{db: db}
	storage.Close()

	err = db.Ping()
	if err == nil {
		t.Fail()
	}
}
