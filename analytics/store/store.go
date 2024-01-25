package store

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"log"
)

type AnalyticsStore struct {
	session *pgx.Conn
}

type Product struct {
	Quantity  int32
	Amount    float64
	ProductId string
}

type Customer struct {
	Customers uuid.UUID
	Quantity  int32
	Amount    float64
}

func NewAnalyticsStore(session *pgx.Conn) *AnalyticsStore {
	return &AnalyticsStore{session: session}
}

func (store *AnalyticsStore) GetTotalSales(ctx context.Context) (float64, error) {
	log.Println("Get Total Sales")
	var amount float64
	err := store.session.QueryRow(ctx,
		"SELECT SUM(AMOUNT) FROM transactions").Scan(
		&amount,
	)
	if err != nil {
		return 0.0, err
	}
	return amount, nil
}

func (store *AnalyticsStore) GetSalesByProduct(ctx context.Context, id uuid.UUID) (*Product, error) {
	log.Println("Get Sales by Product")
	var product = &Product{}
	err := store.session.QueryRow(ctx,
		"SELECT SUM(quantity),SUM(amount), product_id FROM transactions WHERE product_id = $1 group by product_id", id).Scan(
		&product.Quantity,
		&product.Amount,
		&product.ProductId,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (store *AnalyticsStore) GetTopCustomers(ctx context.Context) ([]Customer, error) {
	log.Println("Get Top Customers")
	//var customer = []&Customer{}
	var customers []Customer
	rows, err := store.session.Query(ctx,
		"SELECT customer_id, sum(quantity) as quantity ,sum(amount) as sales  FROM transactions GROUP BY customer_id ORDER BY sales desc limit 5")
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.Customers, &customer.Quantity, &customer.Amount)
		if err != nil {
			log.Fatal(err)
		}
		customers = append(customers, customer)
	}
	if err != nil {
		return nil, err
	}
	return customers, nil
}
