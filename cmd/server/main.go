package main

import (
	"fmt"

	"github.com/brenoproti/go-api/configs"
)

func main() {
	config := configs.LoadConfig(".")
	fmt.Printf("%v", config)
}
