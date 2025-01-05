package types

import (
	"container/heap"
	"sync"
)

type Order struct {
	ID       int
	TraderID int
	Price    float64
	Quantity int
	IsBuy    bool // true for buy, false for sell
	Stock    string
}

// OrderBook maintains orders in priority order.
type OrderBook struct {
	BuyOrders  PriorityQueue // Max-heap for buyers
	SellOrders PriorityQueue // Min-heap for sellers
	mu         sync.Mutex
}

// Transaction represents a completed trade.
type Transaction struct {
	BuyOrderID  int
	SellOrderID int
	Price       float64
	Quantity    int
	Stock       string
}

// Trader represents a market participant.
type Trader struct {
	ID        int
	Cash      float64
	Portfolio map[string]int // Stock symbol to quantity
	Strategy  func(t *Trader, ob *OrderBook)
}

// PriorityQueue implements heap.Interface and holds Orders.
type PriorityQueue []*Order

func (pq PriorityQueue) Len() int { return len(pq) }

// For buy orders, higher price has higher priority.
// For sell orders, lower price has higher priority.
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].IsBuy {
		return pq[i].Price > pq[j].Price
	}
	return pq[i].Price < pq[j].Price
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	order := x.(*Order)
	*pq = append(*pq, order)
}

func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return item
}

// NewOrderBook creates a new order book.
func NewOrderBook() *OrderBook {
	return &OrderBook{
		BuyOrders:  make(PriorityQueue, 0),
		SellOrders: make(PriorityQueue, 0),
	}
}

// AddOrder places an order into the appropriate book.
func (ob *OrderBook) AddOrder(order *Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	if order.IsBuy {
		heap.Push(&ob.BuyOrders, order)
	} else {
		heap.Push(&ob.SellOrders, order)
	}
}

// MatchOrders tries to match buy and sell orders.
func (ob *OrderBook) MatchOrders(traders map[int]*Trader) []Transaction {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	var transactions []Transaction

	for len(ob.BuyOrders) > 0 && len(ob.SellOrders) > 0 {
		buyOrder := ob.BuyOrders[0]
		sellOrder := ob.SellOrders[0]

		if buyOrder.Price >= sellOrder.Price && buyOrder.Stock == sellOrder.Stock {
			quantity := min(buyOrder.Quantity, sellOrder.Quantity)
			transactions = append(transactions, Transaction{
				BuyOrderID:  buyOrder.ID,
				SellOrderID: sellOrder.ID,
				Price:       sellOrder.Price,
				Quantity:    quantity,
				Stock:       buyOrder.Stock,
			})

			// Update quantities and portfolios
			buyOrder.Quantity -= quantity
			sellOrder.Quantity -= quantity

			traders[buyOrder.TraderID].Portfolio[buyOrder.Stock] += quantity
			// traders[sellOrder.TraderID].Portfolio[sellOrder.Stock] -= quantity
			traders[buyOrder.TraderID].Cash += (float64(quantity)*buyOrder.Price - float64(quantity)*sellOrder.Price)
			traders[sellOrder.TraderID].Cash += float64(quantity) * sellOrder.Price

			if buyOrder.Quantity == 0 {
				heap.Pop(&ob.BuyOrders)
			}
			if sellOrder.Quantity == 0 {
				heap.Pop(&ob.SellOrders)
			}
		} else {
			break
		}
	}
	return transactions
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
