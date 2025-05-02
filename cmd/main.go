package main

import "fmt"

func main() {
	fmt.Println("Hello world!")
	str := "абв"
	for _, s := range str {
		fmt.Println(string(s))
	}

}
