# Trading Simulation Program

## Overview
This is a Go-based trading simulation program that models a simple stock exchange. Traders can place buy and sell orders for specific stocks, and an order matching engine periodically processes these orders. The program supports interactive commands to switch between traders, place orders, and view the status of the order book and trader portfolios.

## Features
- **Multiple Traders:** Simulate multiple traders with individual portfolios and cash balances.
- **Order Book Management:** Buy and sell orders are stored and matched based on price and quantity.
- **Interactive CLI:** Users can switch between traders, place orders, and view the current state of the market.
- **Periodic Order Matching:** Orders are matched every second in the background.

## Project Structure
```
/StockMarketSimulationGO
|-- 1.go                  # Main program with interactive CLI and order processing
|-- types/                # Package containing core types (Trader, Order, OrderBook)
|   |-- types.go          # Defines all structuctures and methods
|-- README.md             # Project documentation
```

## How to Run
1. Clone the repository.
2. Navigate to the project directory.
3. Run the following command to start the simulation:
   ```
   go run 1.go
   ```
4. Follow the interactive prompts to switch traders, place orders, and view summaries.

## Commands
- **`help`** - Display the list of available commands.
- **`switch <trader_id>`** - Switch to a specific trader.
- **`order <buy|sell> <quantity> <stock> <price>`** - Place a buy or sell order.
- **`show`** - Display the current order book.
- **`summary`** - Show traders' portfolio and cash balances.
- **`exit`** - Stop the simulation.

## Example Usage
```
> switch 1
Switched to Trader 1
> order buy 5 AAPL 150.00
Order added: {ID:101 TraderID:1 Price:150 Quantity:5 IsBuy:true Stock:AAPL}
> summary
Trader 1 - Cash: 250.00, Portfolio: map[AAPL:10 GOOG:5]
Trader 2 - Cash: 1500.00, Portfolio: map[AAPL:7 AMZN:8]
Trader 3 - Cash: 800.00, Portfolio: map[GOOG:6 AMZN:4]
```

## Requirements
- Go 1.18 or higher

## Notes
- Orders are randomly assigned IDs.
- The program periodically processes orders automatically.
- Ensure that the `types` package contains necessary struct definitions for `Trader`, `Order`, and `OrderBook`.

## Future Enhancements
- Adding AI factor to control Traders - implementing various strategies
- Graphical user interface (GUI).

## License
This project is open-source and licensed under the MIT License.

