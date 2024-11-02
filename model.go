package main

import (
	"database/sql"
	"errors"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "select id,name,quantity,price from products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (p *product) getProduct(db *sql.DB) error {
	query := "SELECT id, name, quantity, price FROM products WHERE id = ?"
	row := db.QueryRow(query, p.ID) // Use parameterized query
	err := row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *product) createProduct(db *sql.DB) error {
	// Use a parameterized query to safely insert values into the table
	query := "INSERT INTO products (name, quantity, price) VALUES (?, ?, ?)"
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price)
	if err != nil {
		return err
	}

	// Retrieve the ID of the newly inserted record
	lastId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(lastId) // You can convert to int if you’re sure it won’t exceed int size
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	// Correct SQL syntax, use parameterized query
	query := "UPDATE products SET name = ?, quantity = ?, price = ? WHERE id = ?"
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price, p.ID)
	rowsEffected, _ := result.RowsAffected()
	if rowsEffected == 0 {
		return errors.New("no such rows exist")
	}
	return err

}

func (p *product) deleteProduct(db *sql.DB) error {
	query := "delete from products where id=?"
	result, err := db.Exec(query, p.ID)
	if err != nil {
		return err
	}

	rowsEffected, err := result.RowsAffected()
	if rowsEffected == 0 {
		return errors.New("no rows affected")
	}

	return err
}
