// Package main is the postgresql lesson `l2_isolation_levels` homework scaffold for Vibe Learn.
//
// Задача: accounts: безопасный Transfer (Serializable+retry или FOR UPDATE) под 100 goroutine.
// Реализуй функции ниже — сигнатуры и тестовая поверхность фиксированы;
// CI (.github/workflows/ci.yml) гоняет `go vet` и `go test ./...`.
// Подробности и критерии приёмки — в README.md.
//
// Драйвер: github.com/jackc/pgx/v5 (+ pgxpool). DATABASE_URL — DSN из env.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Latencies — собранные перцентили для бенчмарка запроса.
type Latencies struct{ P50, P95, P99 time.Duration }

// StandbyInfo — строка из pg_stat_replication для выбора кандидата на promote.
type StandbyInfo struct {
	ClientAddr string
	ReplayLSN  string
	State      string
}

// ----- config -----

// envOr returns the env var for `key` if set, else `fallback`.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// DatabaseURL — DSN PostgreSQL. Дефолт совпадает с docker-compose.yml.
func DatabaseURL() string {
	return envOr("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
}

// Connect открывает пул pgx из DATABASE_URL.
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, DatabaseURL())
}

// ----- TODO #1: Transfer -----
//
// перевод в транзакции: либо Serializable с retry на 40001, либо SELECT FOR UPDATE в детерминированном порядке id
func Transfer(ctx context.Context, pool *pgxpool.Pool, fromID, toID int64, amount int64) error {
	// TODO: implement
	panic("Transfer: not implemented")
}

// ----- TODO #2: TotalBalance -----
//
// SELECT sum(balance) FROM accounts — инвариант не должен меняться
func TotalBalance(ctx context.Context, pool *pgxpool.Pool) (int64, error) {
	// TODO: implement
	panic("TotalBalance: not implemented")
}

// ----- TODO #3: ClaimJob -----
//
// SELECT ... FOR UPDATE SKIP LOCKED LIMIT 1 — забрать задание из pending_jobs
func ClaimJob(ctx context.Context, pool *pgxpool.Pool) (jobID int64, err error) {
	// TODO: implement
	panic("ClaimJob: not implemented")
}

// _refs keeps imports live while the TODO bodies are unimplemented stubs.
// Удали эту функцию, когда реализуешь TODO выше.
var _refs = []any{
	Latencies{},
	StandbyInfo{},
	time.Second,
}

// ----- main entry -----

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Vibe Learn — postgresql lesson %s scaffold up", "l2_isolation_levels")
	log.Printf("DATABASE_URL: %s", DatabaseURL())
	log.Printf("Реализуй TODO-функции, затем `go test ./...`. README.md содержит задачу.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown so `go run .` is interactive — Ctrl-C exits cleanly.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Printf("shutdown signal received")
		cancel()
	}()
	<-ctx.Done()
}
