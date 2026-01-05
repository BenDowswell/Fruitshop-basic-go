package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Product struct {
	Name  string
	Price float64
}

func printInventory(products []Product) {
	fmt.Println("Printing Inventory:")
	for _, p := range products {
		fmt.Printf("%s costs £%.2f\n", p.Name, p.Price)
	}
}

func buildinventory(path string) ([]Product, error) {
	file, err := os.Open(path)
	if err != nil {
		// generic + underlying cause
		return nil, fmt.Errorf("Please check the data source is correct: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	products := []Product{}

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Please check the data source is correct: line %d: expected format name,£price", lineNum)
		}

		name := strings.TrimSpace(parts[0])
		if name == "" {
			return nil, fmt.Errorf("Please check the data source is correct: line %d: product name is empty", lineNum)
		}

		pricePart := strings.TrimSpace(parts[1])
		if !strings.HasPrefix(pricePart, "£") {
			return nil, fmt.Errorf("Please check the data source is correct: line %d: price must start with £", lineNum)
		}

		priceStr := strings.TrimPrefix(pricePart, "£")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, fmt.Errorf("Please check the data source is correct: line %d: invalid price %q: %w", lineNum, priceStr, err)
		}

		products = append(products, Product{
			Name:  name,
			Price: price,
		})
	}

	// scanner-level error (I/O issues etc.)
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Please check the data source is correct: %w", err)
	}

	return products, nil
}

func runCart(products []Product) {
	reader := bufio.NewReader(os.Stdin)

	cart := map[int]int{}
	var total float64

	for {
		printMenu(products)
		fmt.Print("Pick an item number, or 'q' to quit:")

		input := readLine(reader)

		if strings.EqualFold(input, "q") {
			break
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(products) {
			fmt.Println("Please enter a valid item number.")
			continue
		}
		// convert 1-based menu choice to 0-based slice index
		idx := choice - 1

		fmt.Printf("how many %s would you like to buy? ", products[idx].Name)
		qtyStr := readLine(reader)
		qty, err := strconv.Atoi(qtyStr)

		if err != nil || qty <= 0 {
			fmt.Println("Please enter a quantity greater than 0.")
			continue
		}

		cart[idx] += qty
		addedCost := float64(qty) * products[idx].Price
		total += addedCost

		fmt.Printf("Added %d x %s (£%.2f each). Cart total is now £%.2f\n",
			qty, products[idx].Name, products[idx].Price, total)

		fmt.Print("Add another item? (y/n): ")
		again := readLine(reader)
		if strings.EqualFold(again, "n") || strings.EqualFold(again, "no") {
			break
		}
	}

	// Checkout / summary
	if len(cart) == 0 {
		fmt.Println("\nYour cart is empty. Goodbye!")
		return
	}

	fmt.Println("\nCart summary:")
	for idx, p := range products {
		qty := cart[idx]
		if qty == 0 {
			continue
		}
		line := float64(qty) * p.Price
		fmt.Printf("- %d x %s = £%.2f\n", qty, p.Name, line)
	}
	fmt.Printf("Total: £%.2f\n", total)

	fmt.Print("Would you like to pay? (yes/no): ")

	pay := readLine(reader)

	if strings.EqualFold(pay, "yes") || strings.EqualFold(pay, "y") {
		fmt.Printf("Payment accepted. Final total is £%.2f\n", total)
	} else {
		fmt.Println("No problem — goodbye.")
	}

}

func printMenu(products []Product) {
	fmt.Println("\nInventory:")
	for i, p := range products {
		fmt.Printf("%d) %s - £%.2f\n", i+1, p.Name, p.Price)
	}
}

func readLine(r *bufio.Reader) string {
	s, _ := r.ReadString('\n')
	return strings.TrimSpace(s)
}

func main() {
	// start program
	fmt.Println("Fruitshop started")

	products, err := buildinventory("values.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	runCart(products)

}
