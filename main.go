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

func main() {
	fmt.Println("Fruitshop started")

	file, err := os.Open("values.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	products := []Product{}

    for scanner.Scan() {
        line := scanner.Text()

        parts := strings.Split(line, ",")
        if len(parts) != 2 {
            continue
        }

        name := strings.TrimSpace(parts[0])

        priceStr := strings.TrimSpace(parts[1])
        priceStr = strings.TrimPrefix(priceStr, "£")

        price, err := strconv.ParseFloat(priceStr, 64)
        if err != nil {
            continue
        }

        product := Product{
            Name:  name,
            Price: price,
        }

        products = append(products, product)
    }
    fmt.Println("Printing Inventory:")
    for _, p := range products {
	    fmt.Printf("%s costs £%.2f\n", p.Name, p.Price)
    }
}

