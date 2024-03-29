package cache

import (
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/AXBoBka/wb_l0/internal/store"
	"github.com/jackc/pgx/v5"
)

type Cache struct {
	sync.RWMutex
	orders map[string]string
}

func New(conn *pgx.Conn) *Cache {
	orders := store.GetAllOrders(conn)
	if orders == nil {
		orders = make(map[string]string)
	}
	cache := Cache{
		orders: orders,
	}
	log.Println("Создан интсанс для хранения кэша!")
	return &cache
}

func (c *Cache) AddOrder(order string) {
	id := c.findCurrID()
	log.Println("Заказ добавлен в кэш!")
	c.Lock()
	defer c.Unlock()
	c.orders[id] = order
}

func (c *Cache) GetOrder(id string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	order, found := c.orders[id]
	if !found {
		return "", errors.New("Заказ не найден!")
	}
	return order, nil
}

func (c *Cache) findCurrID() string {
	prevID := len(c.orders)
	return strconv.Itoa(prevID + 1)
}
