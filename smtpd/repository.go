package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlRepository struct {
	db *sql.DB
}

func NewMysqlRepository(db *sql.DB) (*MysqlRepository, error) {
	repo := new(MysqlRepository)
	repo.db = db
	_, err := repo.db.Exec(createTableEmail)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (repo *MysqlRepository) SaveEmail(from, tos string, data []byte) error {
	if data == nil {
		data = []byte{}
	}
	_, err := repo.db.Exec(`insert into email("from","tos","data") values(?,?,?)`, from, tos, sql.RawBytes(data))
	return err
}
