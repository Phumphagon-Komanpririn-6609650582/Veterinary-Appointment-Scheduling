package main

import (
	"fmt"
	"veterinary-api/config"
)

func main() {

	config.InitDB()

	fmt.Println("Server is running...")
}
