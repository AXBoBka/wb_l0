package store

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func OpenConnection() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось подключиться к базе данных: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func GetOrder(id int, conn *pgx.Conn) (string, error) {
	order := ""
	err := conn.QueryRow(context.Background(), fmt.Sprintf("select order_data from orders where order_id=%d", id)).Scan(&order)
	return order, err
}

func AddOrder(order []byte, conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO orders (order_data) VALUES ('%s');", order))

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Заказ добавлен в базу данных!")
	return nil
}

func GetAllOrders(conn *pgx.Conn) map[string]string {
	var id, order string
	orders := make(map[string]string)

	rows, _ := conn.Query(context.Background(), "SELECT * FROM orders")

	_, err := pgx.ForEachRow(rows, []any{&id, &order}, func() error {
		orders[id] = order
		return nil
	})

	if err != nil {
		log.Println("...")
		return nil
	}

	return orders
}
