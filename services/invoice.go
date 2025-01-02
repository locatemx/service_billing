package services

import (
	"log"
	"service_billing/db"
)

// InvoiceItem representa un ítem de factura
type InvoiceItem struct {
	ServiceID   *string
	ProductID   *string
	ServiceName string
	ProductName string
	Price       float64
}

// GenerateInvoices genera facturas para cuentas activas
func GenerateInvoices() {
	// Consulta para obtener cuentas, dispositivos, servicios y productos relacionados
	query := `
        SELECT acc.id AS account_id,
               COALESCE(bs.id, NULL) AS service_id,
               COALESCE(bp.id, NULL) AS product_id,
               COALESCE(bsb.invoice_name, 'Sin Nombre') AS service_name,
               COALESCE(bp.alias, 'Sin Nombre') AS product_name,
               COALESCE(bs.price, bsb.default_price, 0) AS service_price,
               COALESCE(bp.price, 0) AS product_price
        FROM account AS acc
        LEFT JOIN account_device AS ad ON ad.account = acc.id
        LEFT JOIN billing_service AS bs ON bs.id = ad.service
        LEFT JOIN billing_service_base AS bsb ON bsb.id = bs.base
        LEFT JOIN billing_product AS bp ON bp.account = acc.id
        WHERE acc.active = true AND acc.deleted = false;
    `

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Fatalf("Error en la consulta: %v", err)
	}
	defer rows.Close()

	accountTotals := make(map[string]float64)
	invoiceItems := make(map[string][]InvoiceItem)

	// Leer los resultados y consolidar totales por cuenta
	for rows.Next() {
		var accountID, serviceName, productName string
		var serviceID, productID *string
		var servicePrice, productPrice float64

		err := rows.Scan(&accountID, &serviceID, &productID, &serviceName, &productName, &servicePrice, &productPrice)
		if err != nil {
			log.Fatalf("Error al leer los datos: %v", err)
		}

		// Sumar los precios de servicios y productos al total de la cuenta
		accountTotals[accountID] += servicePrice + productPrice

		// Agregar ítems de servicio
		if serviceID != nil {
			invoiceItems[accountID] = append(invoiceItems[accountID], InvoiceItem{
				ServiceID:   serviceID,
				ServiceName: serviceName,
				Price:       servicePrice,
			})
		}

		// Agregar ítems de producto
		if productID != nil {
			invoiceItems[accountID] = append(invoiceItems[accountID], InvoiceItem{
				ProductID:   productID,
				ProductName: productName,
				Price:       productPrice,
			})
		}
	}

	// Crear facturas y sus ítems
	for accountID, totalAmount := range accountTotals {
		// Validar si hay ítems para esta cuenta
		if len(invoiceItems[accountID]) == 0 {
			log.Printf("La cuenta %s no tiene ítems asociados. No se generará factura.", accountID)
			continue
		}

		// Crear factura
		invoiceID := createInvoice(accountID, totalAmount)

		// Insertar los ítems para la factura
		for _, item := range invoiceItems[accountID] {
			insertInvoiceItem(invoiceID, item)
		}
	}
}

// createInvoice inserta un registro en la tabla invoice
func createInvoice(accountID string, totalAmount float64) string {
	query := `
        INSERT INTO invoice (id, invoice_date, total_amount, status, created_at, updated_at, due_date)
        VALUES (gen_random_uuid(), CURRENT_DATE, $1, 'pending', NOW(), NOW(), CURRENT_DATE + INTERVAL '7 days')
        RETURNING id;
    `

	var invoiceID string
	err := db.DB.QueryRow(query, totalAmount).Scan(&invoiceID)
	if err != nil {
		log.Fatalf("Error al crear la factura para la cuenta %s: %v", accountID, err)
	}

	return invoiceID
}

// insertInvoiceItem inserta un ítem en la tabla invoice_item
func insertInvoiceItem(invoiceID string, item InvoiceItem) {
	// Validar si el ítem tiene un serviceID o productID
	if item.ServiceID == nil && item.ProductID == nil {
		log.Printf("El ítem no tiene un serviceID ni productID válido. Saltando inserción.")
		return
	}

	_, err := db.DB.Exec(`
        INSERT INTO invoice_item (id, invoice_id, service, product, quantity, price, created_at)
        VALUES (gen_random_uuid(), $1, $2, $3, 1, $4, NOW());
    `, invoiceID, item.ServiceID, item.ProductID, item.Price)
	if err != nil {
		log.Fatalf("Error al insertar ítem de factura %s: %v", invoiceID, err)
	}
}
