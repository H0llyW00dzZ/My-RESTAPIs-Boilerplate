// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"fmt"
	"h0llyw00dz-template/backend/internal/database"
)

// InitializeTables creates all the necessary tables in the database.
//
// Example usage:
//
//	// createUserTable creates the User table in the database if it doesn't exist.
//	func createUserTable(db database.Service) error {
//		// Define the SQL query to create the User table.
//		query := `
//		CREATE TABLE IF NOT EXISTS User (
//			id INT AUTO_INCREMENT PRIMARY KEY,
//			username VARCHAR(255) NOT NULL,
//			email VARCHAR(255) NOT NULL,
//			password VARCHAR(255) NOT NULL,
//			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
//		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
//		`
//
//		// Execute the SQL query without expecting any rows to be returned.
//		err := db.ExecWithoutRow(context.Background(), query)
//		if err != nil {
//			return err
//		}
//
//		// Log a success message indicating that the User table was initialized.
//		log.LogInfo("Successfully initialized the User table.")
//		return nil
//	}
//
//	// InitializeTables creates all the necessary tables in the database.
//	func InitializeTables(db database.Service) error {
//		// Call the createTables function with the database service and table creation functions.
//		return createTables(db,
//			createTable("User", createUserTable),
//			// Add more table creation functions as needed
//		)
//	}
func InitializeTables(db database.Service) error {
	// Note: This approach provides a more flexible and scalable way to initialize database tables compared to using an ORM system.
	// It allows for easy initialization or migration of tables, and can handle a large number of database schemas (e.g, 1 billion database schemas 🔥) without limitations.
	return createTables(db)

}

// createTable is a higher-order function that creates a table in the database.
// It takes the table name and a function that defines the table creation logic.
// It returns a function that accepts a database.Service and executes the table creation logic.
// If an error occurs during table creation, it wraps the error with a descriptive message.
func createTable(name string, fn func(database.Service) error) func(database.Service) error {
	return func(db database.Service) error {
		if err := fn(db); err != nil {
			return fmt.Errorf("failed to create %s table: %v", name, err)
		}
		return nil
	}
}

// createTables creates all the necessary tables in the database.
// It accepts a database.Service and a variadic list of functions that define the table creation logic.
// It iterates over the table creation functions and executes each one.
// If an error occurs during table creation, it returns the error.
func createTables(db database.Service, tables ...func(database.Service) error) error {
	for _, table := range tables {
		if err := table(db); err != nil {
			return err
		}
	}
	return nil
}
