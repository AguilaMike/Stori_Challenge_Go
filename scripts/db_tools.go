package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// This script is used to run database migrations and generate sqlc code.
// It is invoked by the Makefile in the root of the project.
// Usage: go run db_tools.go [migrate|sqlc] [args...]
// Example: go run db_tools.go migrate -direction up
// Example: go run db_tools.go sqlc
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run db_tools.go [migrate|sqlc] [args...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migrate":
		runMigrations()
	case "sqlc":
		generateSqlc()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runMigrations() {
	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	direction := migrateCmd.String("direction", "up", "Migration direction (up or down)")
	migrateCmd.Parse(os.Args[2:])

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	m, err := migrate.New("file://scripts/migrations", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	switch *direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}

	fmt.Println("Migration completed successfully")
}

func generateSqlc() {
	cmd := exec.Command("sqlc", "generate")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error generating sqlc code: %v\n%s", err, output)
	}
	fmt.Println("sqlc code generation completed successfully")
}
