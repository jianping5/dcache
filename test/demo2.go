package main

import (
	"fmt"
)
	

func TestReturn() (t int) {
	fmt.Println(12)
	returnX()
	fmt.Println(13)
	return 15
}

func returnX() (t int) {
	return 5;
}

func main() {
	k := TestReturn()
	fmt.Println(k)
}