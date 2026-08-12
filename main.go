package main

import "fmt"

func main() {
	test := [32][64]int{}

	for row, column := range test {
		fmt.Println(row, column)
	}
}
