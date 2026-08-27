package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/kevinkiplangat432/arvis/internal/policy"
	"github.com/kevinkiplangat432/arvis/internal/tokenize"
)

var tokenizeDemoCmd = &cobra.Command{
	Use:   "tokenize-demo [text]",
	Short: "Locally demo tokenize -> detokenize on sample text, no provider call needed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rdb, err := policy.Connect(cfg.RedisAddr)
		if err != nil {
			return err
		}
		defer rdb.Close()

		ctx := context.Background()
		requestID := uuid.NewString()
		original := []byte(args[0])

		fmt.Println("ORIGINAL:")
		fmt.Println(" ", string(original))

		tokenized, err := tokenize.Tokenize(ctx, rdb, requestID, original)
		if err != nil {
			return fmt.Errorf("tokenize failed: %w", err)
		}
		fmt.Println("\nTOKENIZED — this is what actually leaves the network:")
		fmt.Println(" ", string(tokenized))

		reconstructed, err := tokenize.Detokenize(ctx, rdb, requestID, tokenized)
		if err != nil {
			return fmt.Errorf("detokenize failed: %w", err)
		}
		fmt.Println("\nRECONSTRUCTED — this is what the caller sees:")
		fmt.Println(" ", string(reconstructed))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokenizeDemoCmd)
}