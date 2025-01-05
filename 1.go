package main

import (
	"bufio"
	"day_17/types"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	AvailableStocks := []string{"AAPL", "GOOG", "AMZN"}
	StockOrderBooks := make(map[string]*types.OrderBook)

	for _, s := range AvailableStocks {
		StockOrderBooks[s] = types.NewOrderBook()
	}

	var wg sync.WaitGroup

	traders := map[int]*types.Trader{
		1: {ID: 1, Cash: 1000, Portfolio: map[string]int{"AAPL": 10, "GOOG": 5}},
		2: {ID: 2, Cash: 1500, Portfolio: map[string]int{"AAPL": 7, "AMZN": 8}},
		3: {ID: 3, Cash: 800, Portfolio: map[string]int{"GOOG": 6, "AMZN": 4}},
	}

	stop := make(chan struct{})

	// Goroutine to process orders periodically
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				for stock, ob := range StockOrderBooks {
					transactions := ob.MatchOrders(traders)
					for _, t := range transactions {
						fmt.Printf("Transaction: BuyOrder %d, SellOrder %d, Price %.2f, Quantity %d, Stock %s\n",
							t.BuyOrderID, t.SellOrderID, t.Price, t.Quantity, stock)
					}
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	currentTraderID := 0
	reader := bufio.NewReader(os.Stdin)

	// Interactive command loop
	fmt.Println("Enter commands (type 'help' for a list of commands):")
	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		args := strings.Split(cmd, " ")

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  switch <trader_id> - Switch to a specific trader")
			fmt.Println("  order <buy|sell> <quantity> <stock> <price> - Place an order")
			fmt.Println("  show - Display the current order book")
			fmt.Println("  traders - List all available traders")
			fmt.Println("  summary - Show traders' portfolio and cash")
			fmt.Println("  exit - Stop the simulation")

		case "switch":
			if len(args) < 2 {
				fmt.Println("Usage: switch <trader_id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil || traders[id] == nil {
				fmt.Println("Invalid trader ID")
				continue
			}
			currentTraderID = id
			fmt.Printf("Switched to Trader %d\n", currentTraderID)

		case "order":
			if currentTraderID == 0 {
				fmt.Println("Switch to a trader first using 'switch <trader_id>'")
				continue
			}
			if len(args) < 5 {
				fmt.Println("Usage: order <buy|sell> <quantity> <stock> <price>")
				continue
			}
			isBuy := args[1] == "buy"
			quantity, err1 := strconv.Atoi(args[2])
			price, err2 := strconv.ParseFloat(args[4], 64)
			if err1 != nil || err2 != nil {
				fmt.Println("Invalid quantity or price")
				continue
			}
			stock := args[3]
			if !isBuy {
				if traders[currentTraderID].Portfolio[stock] < quantity {
					fmt.Println("Invalid quantity")
					continue
				}
				traders[currentTraderID].Portfolio[stock] -= quantity

			}
			if isBuy {
				if traders[currentTraderID].Cash < float64(quantity)*price {
					fmt.Println("Invalid price")
					continue
				}
				traders[currentTraderID].Cash -= float64(quantity) * price
			}
			order := &types.Order{
				ID:       rand.Intn(1000),
				TraderID: currentTraderID,
				Price:    price,
				Quantity: quantity,
				IsBuy:    isBuy,
				Stock:    stock,
			}
			StockOrderBooks[stock].AddOrder(order)
			fmt.Printf("Order added: %+v\n", order)

		case "show":
			fmt.Println("Order Book:")
			for key, val := range StockOrderBooks {
				fmt.Printf("%s\n", key)
				for i := 0; i < max(len(val.BuyOrders), len(val.SellOrders)); i++ {
					if i < len(val.BuyOrders) {
						fmt.Printf("Buy %d: %+v    ", i+1, val.BuyOrders[i])
					}
					if i >= len(val.BuyOrders) {
						fmt.Printf("                ")
					}
					if i < len(val.SellOrders) {
						fmt.Printf("Sell %d: %+v\n", i+1, val.SellOrders[i])
					}
					if i >= len(val.SellOrders) {
						fmt.Printf("\n")
					}
				}
			}

		case "traders":
			fmt.Println("Available Traders:")
			for id, t := range traders {
				fmt.Printf("Trader %d - Cash: %.2f, Portfolio: %+v\n", id, t.Cash, t.Portfolio)
			}

		case "summary":
			fmt.Println("Traders Summary:")
			for _, t := range traders {
				fmt.Printf("Trader %d - Cash: %.2f, Portfolio: %+v\n", t.ID, t.Cash, t.Portfolio)
			}

		case "exit":
			close(stop)
			wg.Wait()
			return

		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
