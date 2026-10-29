//go:build integration

// Package integration testa os repositórios contra um Postgres 15 real
// (testcontainers), com as migrations de database/migrations aplicadas.
package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:15",
		postgres.WithDatabase("elephrases_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Printf("failed to start postgres container: %v", err)
		return 1
	}
	defer func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}()

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Printf("failed to get connection string: %v", err)
		return 1
	}

	testDB, err = gorm.Open(gormpostgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Printf("failed to connect: %v", err)
		return 1
	}

	if err := applyMigrations(testDB); err != nil {
		log.Printf("failed to apply migrations: %v", err)
		return 1
	}

	return m.Run()
}

// applyMigrations executa os arquivos *.up.sql reais, em ordem de versão.
func applyMigrations(db *gorm.DB) error {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "database", "migrations")

	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no migrations found in %s", dir)
	}
	sort.Strings(files)

	for _, f := range files {
		sql, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if err := db.Exec(string(sql)).Error; err != nil {
			return fmt.Errorf("%s: %w", filepath.Base(f), err)
		}
	}
	return nil
}

// resetDB limpa todas as tabelas entre os testes.
func resetDB(t *testing.T) {
	t.Helper()
	if err := testDB.Exec("TRUNCATE reviews, flashcards, decks, users CASCADE").Error; err != nil {
		t.Fatalf("reset db: %v", err)
	}
}

func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}
