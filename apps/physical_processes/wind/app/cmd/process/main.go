package main

import (
	"fmt"
	"wind-process/internal/random"
)

func main() {
	fmt.Println(random.CreateRandInt(4, 423))
	fmt.Println(random.CreateRandFloat(3.43, 3.4343))
}
