package main

import (
	"log"
	"os"
	"service_billing/db"
	"service_billing/services"

	"github.com/joho/godotenv"
)

func loadEnv() {
	// Carga el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env, se usarán las variables de entorno.")
	}
}

func connectDatabase() {
	// Obtener la cadena de conexión desde las variables de entorno
	connectionString := os.Getenv("PG_CONNECTION")
	if connectionString == "" {
		log.Fatal("Error: La variable PG_CONNECTION no está definida.")
	}

	// Conectar a la base de datos
	db.ConnectDB(connectionString)
	log.Println("Conexión a PostgreSQL completada.")
}

func main() {
	// Cargar variables de entorno
	loadEnv()

	// Conectar a la base de datos
	connectDatabase()
	defer db.DB.Close()

	// Iniciar servicios
	log.Println("Aplicación iniciada.")
	services.GenerateInvoices()
	log.Println("Facturación completada con éxito.")
}
