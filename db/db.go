package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Driver para PostgreSQL
)

var DB *sql.DB

func ConnectDB(connectionString string) {
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("La base de datos no está disponible: %v", err)
	}

	fmt.Println("Conexión exitosa a PostgreSQL!")
}
