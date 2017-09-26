package main

import (
	"database/sql"
	"errors"
)

type dialog struct {
	ID    int     `json:"id"`
	string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *dialog) getProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *dialog) updateProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *dialog) deleteProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *dialog) createProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getdialogs(db *sql.DB, start, count int) ([]dialog, error) {
	return nil, errors.New("Not implemented")
}
