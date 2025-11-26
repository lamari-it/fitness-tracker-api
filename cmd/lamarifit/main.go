package main

import (
	"fmt"
	"lamari-fit-api/cmd/lamarifit/commands"
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
		fmt.Println("LamariFit API CLI v1.0.0")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`LamariFit API CLI

Usage:
  lamarifit <command> [arguments]

Available Commands:
  db        Database operations (seed, reset, migrate, etc.)
  help      Show this help message
  version   Show version information

Database Commands:
  lamarifit db seed              Run database seeders
  lamarifit db seed:fresh        Drop all data and re-seed
  lamarifit db reset             Reset database (migrate:fresh + seed)
  lamarifit db migrate           Run pending migrations
  lamarifit db migrate:fresh     Drop all tables and re-run migrations
  lamarifit db migrate:rollback  Rollback last migration batch

Examples:
  lamarifit db seed              # Run all seeders
  lamarifit db reset             # Reset database and seed
  lamarifit db migrate           # Run migrations

Use "lamarifit help <command>" for more information about a command.`)
}
