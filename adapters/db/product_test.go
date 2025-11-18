package db_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/JosimarReis/go-hexagonal/adapters/db"
	"github.com/JosimarReis/go-hexagonal/application"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var Db *sql.DB

func setUp() {
	var err error

	Db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	createTable(Db)
	createProduct(Db)
}

func createTable(db *sql.DB) {
	table := `CREATE TABLE products (
    			id STRING, 
    			name STRING, 
    			price FLOAT, 
    			status STRING);`

	stmt, err := db.Prepare(table)
	if err != nil {
		log.Fatal(err.Error())
	}
	stmt.Exec()
}

func createProduct(db *sql.DB) {
	insert := `INSERT INTO products (id, name, price, status) VALUES ("abc", "Product Test", 0, "disabled")`
	stmt, err := db.Prepare(insert)
	if err != nil {
		log.Fatal(err.Error())
	}
	stmt.Exec()
}

func TestProductDb_Get(t *testing.T) {
	setUp()
	defer Db.Close()
	productDb := db.NewProductDb(Db)
	product, err := productDb.Get("abc")
	require.Nil(t, err)
	require.Equal(t, "Product Test", product.GetName())
	require.Equal(t, 0.0, product.GetPrice())
	require.Equal(t, "disabled", product.GetStatus())

}

func TestProductDb_Save(t *testing.T) {
	setUp()
	defer Db.Close()
	productDb := db.NewProductDb(Db)

	product := application.NewProduct()
	product.Name = "Product Test"
	product.Price = 25

	result, err := productDb.Save(product)
	require.Nil(t, err)
	require.Equal(t, product.Name, result.GetName())
	require.Equal(t, product.Price, result.GetPrice())
	require.Equal(t, product.Status, result.GetStatus())

	product.Status = "enabled"
	result, err = productDb.Save(product)
	require.Nil(t, err)
	require.Equal(t, product.Name, result.GetName())
	require.Equal(t, product.Price, result.GetPrice())
	require.Equal(t, product.Status, result.GetStatus())
}
