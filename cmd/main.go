package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kevinkiplangat432/arvis/internal/api"
)

const banner = `
 █████╗  ██████╗  ██╗   ██╗ ██╗ ███████╗
██╔══██╗ ██╔══██╗ ██║   ██║ ██║ ██╔════╝
███████║ ██████╔╝ ██║   ██║ ██║ ███████╗
██╔══██║ ██╔══██╗ ╚██╗ ██╔╝ ██║ ╚════██║
██║  ██║ ██║  ██║  ╚████╔╝  ██║ ███████║
╚═╝  ╚═╝ ╚═╝  ╚═╝   ╚═══╝   ╚═╝ ╚══════╝

Automated Runtime Visibility & Intelligence System
Version 0.4.0
`

func main() {
	srvRouter := api.NewRouter()

	fmt.Print(banner)
	if err := http.ListenAndServe(":8080", srvRouter); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}