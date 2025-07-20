package main

import (
	"fit-flow-api/cmd/fitflow/commands"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "db":
		commands.HandleDatabaseCommand(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	case "version", "-v", "--version":
		fmt.Println("FitFlow API CLI v1.0.0")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`FitFlow API CLI

Usage:
  fitflow <command> [arguments]

Available Commands:
  db        Database operations (seed, reset, migrate, etc.)
  help      Show this help message
  version   Show version information

Database Commands:
  fitflow db seed              Run database seeders
  fitflow db seed:fresh        Drop all data and re-seed
  fitflow db reset             Reset database (migrate:fresh + seed)
  fitflow db migrate           Run pending migrations
  fitflow db migrate:fresh     Drop all tables and re-run migrations
  fitflow db migrate:rollback  Rollback last migration batch

Examples:
  fitflow db seed              # Run all seeders
  fitflow db reset             # Reset database and seed
  fitflow db migrate           # Run migrations

Use "fitflow help <command>" for more information about a command.`)
}