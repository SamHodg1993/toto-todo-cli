package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var sql_get_todos string = "SELECT id, title, completed FROM todos"
var sql_get_todos_LONG string = "SELECT id, title, description, project_id, created_at, updated_at, completed FROM todos"

var fullDate bool = false

var getCmd = &cobra.Command{
	Use:   "list",
	Short: "List todo's for current project",
	Long:  "Get a list of all the todo's for the current project (defined by the current directory).",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		rows, err := GetTodosForFilepath()
		if err != nil {
			fmt.Printf("%v.\n", err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Todo", "Status"})
		table.SetBorder(true)
		table.SetRowLine(true)

		strikethrough := color.New(color.CrossedOut).SprintFunc()

		for rows.Next() {
			var (
				id        int
				title     string
				completed bool
			)

			err := rows.Scan(&id, &title, &completed)
			if err != nil {
				fmt.Printf("Error reading row: %v\n", err)
				return
			}

			// If todo is completed, apply strikethrough to the title
			if completed {
				title = strikethrough(title)
			}

			status := "Pending"
			if completed {
				status = "Done"
			}

			table.Append([]string{
				fmt.Sprintf("%d", id),
				title,
				status,
			})
		}

		table.Render()

		if err := rows.Err(); err != nil {
			fmt.Printf("Error iterating over rows: %v\n", err)
		}
	},
}

var getCmdLong = &cobra.Command{
	Use:   "list-long",
	Short: "List todo's with more data for the current project.",
	Long:  "Get a more detailed list of all the todo's for the current project (defined by the current directory)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		rows, err := GetTodosForFilepath_LONG()
		if err != nil {
			fmt.Printf("%v.\n", err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Todo", "Description", "Project Id", "Created At", "Updated At", "Status"})
		table.SetBorder(true)
		table.SetRowLine(true)

		strikethrough := color.New(color.CrossedOut).SprintFunc()

		for rows.Next() {
			var (
				id          int
				title       string
				completed   bool
				description sql.NullString
				projectId   int
				createdAt   time.Time
				updatedAt   time.Time
			)

			err := rows.Scan(&id, &title, &description, &projectId, &createdAt, &updatedAt, &completed)
			if err != nil {
				fmt.Printf("Error reading row: %v\n", err)
				return
			}

			// If todo is completed, apply strikethrough to the title
			if completed {
				title = strikethrough(title)
			}

			status := "Pending"
			if completed {
				status = "Done"
			}

			if fullDate {
				table.Append([]string{
					fmt.Sprintf("%d", id),
					title,
					description.String,
					strconv.Itoa(projectId),
					createdAt.Format(time.RFC3339),
					updatedAt.Format(time.RFC3339),
					status,
				})
			} else {
				table.Append([]string{
					fmt.Sprintf("%d", id),
					title,
					description.String,
					strconv.Itoa(projectId),
					createdAt.Format("02-01-2006"),
					updatedAt.Format("02-01-2006"),
					status,
				})
			}
		}

		table.Render()

		if err := rows.Err(); err != nil {
			fmt.Printf("Error iterating over rows: %v\n", err)
		}
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List outstanding todo's",
	Long:  "Get a list of all the todo titles that are outstanding",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		rows, err := GetTodosForFilepath()
		if err != nil {
			if err.Error() == "operation cancelled by user" {
				return
			}
			fmt.Printf("%v\n", err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Todo", "Status"})
		table.SetBorder(true)
		table.SetRowLine(true)

		strikethrough := color.New(color.CrossedOut).SprintFunc()

		for rows.Next() {
			var (
				id        int
				title     string
				completed bool
			)

			err := rows.Scan(&id, &title, &completed)
			if err != nil {
				fmt.Printf("Error reading row: %v\n", err)
				return
			}

			// If todo is completed, apply strikethrough to the title
			if completed {
				title = strikethrough(title)
			}

			status := "Pending"
			if completed {
				status = "Done"
			}

			table.Append([]string{
				fmt.Sprintf("%d", id),
				title,
				status,
			})
		}

		table.Render()

		if err := rows.Err(); err != nil {
			fmt.Printf("Error iterating over rows: %v\n", err)
		}
	},
}

var lsCmdLong = &cobra.Command{
	Use:   "lsl",
	Short: "List todo's with more data",
	Long:  "Get a more detailed list of all the todo's for the current project (defined by the current directory)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		rows, err := GetTodosForFilepath_LONG()
		if err != nil {
			fmt.Printf("%v.\n", err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Todo", "Description", "ProjectId", "Created At", "Updated At", "Status"})
		table.SetBorder(true)
		table.SetRowLine(true)

		strikethrough := color.New(color.CrossedOut).SprintFunc()

		for rows.Next() {
			var (
				id          int
				title       string
				completed   bool
				description sql.NullString
				projectId   int
				createdAt   time.Time
				updatedAt   time.Time
			)

			err := rows.Scan(&id, &title, &description, &projectId, &createdAt, &updatedAt, &completed)
			if err != nil {
				fmt.Printf("Error reading row: %v\n", err)
				return
			}

			// If todo is completed, apply strikethrough to the title
			if completed {
				title = strikethrough(title)
			}

			status := "Pending"
			if completed {
				status = "Done"
			}

			if fullDate {
				table.Append([]string{
					fmt.Sprintf("%d", id),
					title,
					description.String,
					strconv.Itoa(projectId),
					createdAt.Format(time.RFC3339),
					updatedAt.Format(time.RFC3339),
					status,
				})
			} else {
				table.Append([]string{
					fmt.Sprintf("%d", id),
					title,
					description.String,
					strconv.Itoa(projectId),
					createdAt.Format("02-01-2006"),
					updatedAt.Format("02-01-2006"),
					status,
				})
			}
		}

		table.Render()

		if err := rows.Err(); err != nil {
			fmt.Printf("Error iterating over rows: %v\n", err)
		}
	},
}

func init() {
	lsCmdLong.Flags().BoolVarP(&fullDate, "Full-Date", "D", false, "Return the dates as full timestamps")
	getCmdLong.Flags().BoolVarP(&fullDate, "Full-Date", "D", false, "Return the dates as full timestamps")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(getCmdLong)
	rootCmd.AddCommand(lsCmdLong)
}
