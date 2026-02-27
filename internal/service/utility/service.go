package utility

import (
	"database/sql"
	"fmt"
	"os"
	"slices"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/odgy8/toto/internal/models"
	todoHelper "github.com/odgy8/toto/internal/service/todo"
	"github.com/odgy8/toto/internal/utilities"
)

// TodoServiceInterface defines methods needed from todo service
type TodoServiceInterface interface {
	RemoveCompletedTodosForProject(projectId int) error
}

// ProjectServiceInterface defines methods needed from project service
type ProjectServiceInterface interface {
	HandleNoExistingProject() (int, error)
	AddNewProjectWithPrompt() error
}

// Service handles utility command operations
type Service struct {
	db             *sql.DB
	todoService    TodoServiceInterface
	projectService ProjectServiceInterface
}

// New creates a new utility commands service
func New(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

// SetDependencies allows injecting service dependencies
func (s *Service) SetDependencies(todoService TodoServiceInterface, projectService ProjectServiceInterface) {
	s.todoService = todoService
	s.projectService = projectService
}

// CleanAndPrintTodos clears screen, removes completed todos, and displays remaining todos
func (s *Service) CleanAndPrintTodos(reverseList bool) error {
	utilities.ClearScreen()

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	row := s.db.QueryRow("SELECT id FROM projects WHERE filepath = ?", currentDir)

	var projectId int = 0
	err = row.Scan(&projectId)
	if err != nil {
		choice, err := s.projectService.HandleNoExistingProject()
		if err != nil {
			return err
		}
		if choice == 2 {
			if err := s.projectService.AddNewProjectWithPrompt(); err != nil {
				return err
			}
			return s.CleanAndPrintTodos(reverseList)
		}
	}

	if projectId == 0 {
		fmt.Printf("No project exists for this working directory, defaulting to the global project\n")
		projectId = 1
	}

	err = s.todoService.RemoveCompletedTodosForProject(projectId)
	if err != nil {
		return err
	}

	rows, err := s.db.Query("SELECT id, title, completed FROM todos WHERE project_id = ?", projectId)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	defer rows.Close()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Todo", "Status"})
	table.SetBorder(true)
	table.SetRowLine(true)

	strikethrough := color.New(color.CrossedOut).SprintFunc()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err = rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		todos = append(todos, todo)
	}

	if reverseList {
		slices.Reverse(todos)
	}

	for _, todo := range todos {
		table.Append(todoHelper.FormatTodoTableRow(todo, strikethrough))
	}

	table.Render()
	return nil
}
