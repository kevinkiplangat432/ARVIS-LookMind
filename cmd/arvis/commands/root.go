package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/spf13/cobra"
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

// cfg holds the global application configurations
var cfg *config.Config

var rootCmd = &cobra.Command{
	Use: "arvis",
	Short: "ARVIS - AI infrastructure monitoring proxy",
	Long: banner,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(banner)
		fmt.Println("Select mode:")
		fmt.Println("  [1] Run server")
		fmt.Println("  [2] Run tests")
		fmt.Println("  [3] Run migrations")
		fmt.Print("\nEnter choice: ")
		
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			return serverCmd.RunE(cmd, args)
		case "2":
			return testCmd.RunE(cmd, args)
		case "3":
			return migrateCmd.RunE(cmd, args)
		default:
			return fmt.Errorf("invalid choice %q — enter 1, 2, or 3", choice)
		}
	},
}




func Execute() {
	_ = godotenv.Load()

	loaded, err := config.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}
	cfg = loaded

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}


func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(migrateCmd)
}