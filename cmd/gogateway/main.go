package main

import (
	"github.com/PDK1744/gogateway/internal/server"
	// "github.com/PDK1744/gogateway/internal/config"
)

func main() {
	// _, err := config.LoadConfig("/home/kobeb/KobeCodes/gogateway/sampleconfig.yaml")
	// if err != nil {
	// 	fmt.Printf("ERROR: %v", err)
	// }
	server.StartServer()
}
