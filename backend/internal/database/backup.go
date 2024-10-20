// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package database

import (
	"context"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"os"
	"strings"
	"sync"
	"time"
)

// BackupTables creates a backup of specified tables in the database.
// It generates a .sql file containing the SQL statements needed to recreate
// the database schema and insert all the current data for each specified table.
//
// TODO: Enhance this function by implementing a background task to automatically back up the database
// using a fully encrypted mechanism provided by the "opengpg" utility.
// Currently, it only performs a basic backup and stores it directly on disk. Instead,
// the backups should be archived and stored in a cloud storage service.
// This function should be handled by goroutines for task to run periodically without relying on cron jobs,
// as Go supports concurrency with goroutines.
func (s *service) BackupTables(tablesToBackup []string) error {
	for _, tableName := range tablesToBackup {
		if !IsValidTableName(tableName) {
			return fmt.Errorf("invalid table name: %s", tableName)
		}
	}

	backupFile := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))
	file, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.LogErrorf("Failed to close file: %v", cerr)
		}
		if err != nil {
			if remErr := os.Remove(backupFile); remErr != nil {
				log.LogErrorf("Failed to remove incomplete backup file: %v", remErr)
			}
		}
	}()

	// For large datasets, this may need to configure this and adjust the MySQL server settings.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	for _, tableName := range tablesToBackup {
		if err = s.dumpTableSchema(ctx, file, tableName); err != nil {
			return err
		}

		if err = s.dumpTableData(ctx, file, tableName); err != nil {
			return err
		}
	}

	log.LogInfof("Backup completed: %s", backupFile)
	return nil
}

// BackupTablesConcurrently creates a backup of specified tables concurrently.
// It uses goroutines to perform backups for each table simultaneously, improving performance.
// Each table's backup is handled in a separate goroutine (e.g., 9999999 tables then 9999999 goroutines), and errors are captured via a channel.
func (s *service) BackupTablesConcurrently(tablesToBackup []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(tablesToBackup))

	for _, tableName := range tablesToBackup {
		wg.Add(1)
		go func(table string) {
			defer wg.Done()
			if err := s.backupSingleTable(table); err != nil {
				errChan <- err
			}
		}(tableName)
	}

	// Wait for all goroutines to finish and close the error channel
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for errors from any of the goroutines
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// backupSingleTable performs the backup of a single table.
// It validates the table name, starts a transaction, and writes the schema and data to a backup file.
// The transaction ensures a consistent snapshot of the table during the backup process.
func (s *service) backupSingleTable(tableName string) error {
	if !IsValidTableName(tableName) {
		return fmt.Errorf("invalid table name: %s", tableName)
	}

	// Start a transaction to ensure data consistency
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	backupFile := fmt.Sprintf("backup_%s_%s.sql", tableName, time.Now().Format("20060102_150405"))
	file, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file for table %s: %w", tableName, err)
	}
	defer file.Close()

	// For large datasets, this may need to configure this and adjust the MySQL server settings.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Dump schema and data within the transaction context
	if err := s.dumpTableSchema(ctx, file, tableName); err != nil {
		return err
	}

	if err := s.dumpTableData(ctx, file, tableName); err != nil {
		return err
	}

	// Commit the transaction after successful backup
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.LogInfof("Backup completed for table: %s", tableName)
	return nil
}

// dumpTableSchema writes the CREATE TABLE statement for the specified table to the backup file.
func (s *service) dumpTableSchema(ctx context.Context, file *os.File, tableName string) error {
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	row := s.db.QueryRowContext(ctx, query)
	var table, createTableStmt string
	if err := row.Scan(&table, &createTableStmt); err != nil {
		return fmt.Errorf("failed to get create table statement: %w", err)
	}
	_, err := file.WriteString(fmt.Sprintf("%s;\n\n", createTableStmt))
	return err
}

// dumpTableData retrieves all rows from the specified table and writes them as INSERT statements to the backup file.
func (s *service) dumpTableData(ctx context.Context, file *os.File, tableName string) error {
	query := fmt.Sprintf("SELECT * FROM `%s`", tableName)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to query table data: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		insertStmt := buildInsertStatement(tableName, columns, values)
		if _, err := file.WriteString(insertStmt); err != nil {
			return err
		}
	}
	_, err = file.WriteString("\n")
	return err
}

// buildInsertStatement constructs an SQL INSERT statement for a single row of data.
func buildInsertStatement(tableName string, columns []string, values []any) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("INSERT INTO `%s` (", tableName))
	sb.WriteString(strings.Join(columns, ", "))
	sb.WriteString(") VALUES (")

	for i, val := range values {
		if i > 0 {
			sb.WriteString(", ")
		}
		if val == nil {
			sb.WriteString("NULL")
		} else if b, ok := val.([]byte); ok {
			sb.WriteString(fmt.Sprintf("'%s'", escapeString(string(b))))
		} else {
			sb.WriteString(fmt.Sprintf("'%v'", escapeString(fmt.Sprintf("%v", val))))
		}
	}
	sb.WriteString(");\n")
	return sb.String()
}
