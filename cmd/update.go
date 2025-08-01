package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"jet/internal/config"
	"jet/internal/jira"
)

var (
	updateDescription string
	updateDescFile    string
	updateEpic        string
	updateParent      string
)

var editCmd = &cobra.Command{
	Use:   "edit TICKET-KEY",
	Short: "Edit a JIRA ticket",
	Long: `Edit fields of a JIRA ticket.
	
Currently supports editing the description field and epic/parent linking.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticketKey := args[0]

		// Check if any update flags are provided
		if updateDescription == "" && updateDescFile == "" && updateEpic == "" && updateParent == "" {
			return fmt.Errorf("no update fields specified. Use --description, --description-file, --epic, or --parent")
		}

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration error: %w", err)
		}

		// Create JIRA client
		client := jira.NewClient(cfg.URL, cfg.Email, cfg.Username, cfg.Token)

		// Prepare fields to update
		fields := make(map[string]interface{})

		// Handle description update
		if updateDescFile != "" {
			var content []byte
			if updateDescFile == "-" {
				// Read from stdin
				content, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
			} else {
				// Read from file
				content, err = os.ReadFile(updateDescFile)
				if err != nil {
					return fmt.Errorf("failed to read description file: %w", err)
				}
			}
			fields["description"] = strings.TrimSpace(string(content))
		} else if updateDescription != "" {
			fields["description"] = updateDescription
		}

		// Handle epic/parent update
		if updateEpic != "" {
			fields["parent"] = map[string]string{"key": updateEpic}
		}

		if updateParent != "" {
			fields["parent"] = map[string]string{"key": updateParent}
		}

		// Update the ticket
		if err := client.UpdateIssue(ticketKey, fields); err != nil {
			return err
		}

		fmt.Printf("Ticket %s updated successfully\n", ticketKey)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	
	editCmd.Flags().StringVar(&updateDescription, "description", "", "New description for the ticket")
	editCmd.Flags().StringVar(&updateDescFile, "description-file", "", "Read new description from file (use '-' for stdin)")
	editCmd.Flags().StringVar(&updateEpic, "epic", "", "Epic key to link this ticket to")
	editCmd.Flags().StringVar(&updateParent, "parent", "", "Parent ticket key to link this ticket to")
}