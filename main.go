package main

import (
	"log"
	"os"
	"service_billing/db"
	"service_billing/services"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	connectionString := os.Getenv("PG_CONNECTION")
	if connectionString == "" {
		log.Fatal("Error: La variable PG_CONNECTION no está definida")
	}

	db.ConnectDB(connectionString)
	defer db.DB.Close()

	log.Println("Conexión a PostgreSQL completada. Aplicación iniciada.")
	services.GenerateInvoices()
	log.Println("Facturación completada.")
}
