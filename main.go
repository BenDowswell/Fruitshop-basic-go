package main

import (
	"bufio"
	"fmt"
    "strings"
    "strconv"
	"os"
    
)

type Product struct {
	Name  string
	Price float64
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


func main() {
	// start program 
	fmt.Println("Fruitshop started")
    fmt.Println("Printing Inventory:")
	products, err := buildinventory("values.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(products)
}