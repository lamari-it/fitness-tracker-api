package commands

import (
	"fit-flow-api/config"
	"fit-flow-api/database"
	"fmt"
	"os"
)

func HandleDatabaseCommand(args []string) {
	if len(args) == 0 {
		printDatabaseUsage()
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	subCommand := args[0]

	switch subCommand {
	case "seed":
		err = database.RunSeeders(db)
		if err != nil {
			fmt.Printf("Error running seeders: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database seeded successfully!")

	case "seed:fresh":
		fmt.Println("Dropping all data and re-seeding...")
		err = database.DropAllData(db)
		if err != nil {
			fmt.Printf("Error dropping data: %v\n", err)
			os.Exit(1)
		}
		err = database.RunSeeders(db)
		if err != nil {
			fmt.Printf("Error running seeders: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database seeded successfully!")

	case "re-seed":
		if cfg.Environment == "production" {
			fmt.Println("Error: re-seed command is not allowed in production environment!")
			fmt.Println("This command can only be run when APP_ENV is set to 'development' or 'staging'")
			os.Exit(1)
		}
		
		fmt.Println("🚨 WARNING: This will DELETE ALL DATA and re-seed the database!")
		fmt.Println("Environment: " + cfg.Environment)
		fmt.Print("Are you sure you want to continue? (yes/no): ")
		
		var response string
		fmt.Scanln(&response)
		
		if response != "yes" {
			fmt.Println("Operation cancelled")
			os.Exit(0)
		}
		
		fmt.Println("Dropping all data...")
		err = database.DropAllData(db)
		if err != nil {
			fmt.Printf("Error dropping data: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Re-seeding database...")
		err = database.RunSeeders(db)
		if err != nil {
			fmt.Printf("Error running seeders: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Database re-seeded successfully!")

	case "reset":
		fmt.Println("Resetting database...")
		if cfg.UseMigrations {
			err = database.MigrateFresh(db, cfg)
			if err != nil {
				fmt.Printf("Error running migrations: %v\n", err)
				os.Exit(1)
			}
		} else {
			database.DB = db
			database.AutoMigrate()
		}
		err = database.RunSeeders(db)
		if err != nil {
			fmt.Printf("Error running seeders: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database reset successfully!")

	case "migrate":
		if cfg.UseMigrations {
			err = database.Migrate(db, cfg)
			if err != nil {
				fmt.Printf("Error running migrations: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Migrations completed successfully!")
		} else {
			database.DB = db
			database.AutoMigrate()
			fmt.Println("Auto-migration completed successfully!")
		}

	case "migrate:fresh":
		if cfg.UseMigrations {
			err = database.MigrateFresh(db, cfg)
			if err != nil {
				fmt.Printf("Error running migrations: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Fresh migrations completed successfully!")
		} else {
			fmt.Println("Dropping all tables and re-running auto-migrate...")
			err = database.DropAllTables(db)
			if err != nil {
				fmt.Printf("Error dropping tables: %v\n", err)
				os.Exit(1)
			}
			database.DB = db
			database.AutoMigrate()
			fmt.Println("Fresh auto-migration completed successfully!")
		}

	case "migrate:rollback":
		if cfg.UseMigrations {
			err = database.RollbackMigration(db, cfg)
			if err != nil {
				fmt.Printf("Error rolling back migration: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Migration rolled back successfully!")
		} else {
			fmt.Println("Rollback is only available when using migrations (USE_MIGRATIONS=true)")
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown database command: %s\n", subCommand)
		printDatabaseUsage()
		os.Exit(1)
	}
}

func printDatabaseUsage() {
	fmt.Println(`Database Commands:
  fitflow db seed              Run database seeders
  fitflow db seed:fresh        Drop all data and re-seed
  fitflow db reset             Reset database (migrate:fresh + seed)
  fitflow db migrate           Run pending migrations
  fitflow db migrate:fresh     Drop all tables and re-run migrations
  fitflow db migrate:rollback  Rollback last migration batch`)
}