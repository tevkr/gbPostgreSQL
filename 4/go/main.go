package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type (
	Port int
)

type AvailablePortSearchHint struct {
	Port Port
}

func findAvailableUdpPort(ctx context.Context, dbpool *pgxpool.Pool) ([]AvailablePortSearchHint, error) {
	const sql = `SELECT port FROM available_udp_ports WHERE availability = TRUE LIMIT 1;`

	rows, err := dbpool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	// В слайс hints будут собраны все строки, полученные из базы
	var hints []AvailablePortSearchHint

	// rows.Next() итерируется по всем строкам, полученным из базы.
	for rows.Next() {
		var hint AvailablePortSearchHint

		// Scan записывает значения столбцов в свойства структуры hint
		err = rows.Scan(&hint.Port)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		hints = append(hints, hint)
	}

	// Проверка, что во время выборки данных не происходило ошибок
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}

	return hints, nil
}

type (
	Id string
	Name string
	LobbyId string
)

type User struct {
	Id Id
	Name Name
	LobbyId LobbyId
}

func findLobbyUsers(ctx context.Context, dbpool *pgxpool.Pool, lobbyId string) ([]User, error) {
	const sql = `SELECT id, name FROM users WHERE lobby_id = $1;`

	rows, err := dbpool.Query(ctx, sql, lobbyId)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	// В слайс hints будут собраны все строки, полученные из базы
	var hints []User

	// rows.Next() итерируется по всем строкам, полученным из базы.
	for rows.Next() {
		var hint User

		// Scan записывает значения столбцов в свойства структуры hint
		err = rows.Scan(&hint.Id, &hint.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		hints = append(hints, hint)
	}

	// Проверка, что во время выборки данных не происходило ошибок
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}

	return hints, nil
}

func updateUsersLobbyId(ctx context.Context, dbpool *pgxpool.Pool, lobbyId string, userId string) (LobbyId, error) {
	const sql = `UPDATE users SET lobby_id = $1 WHERE id = $2 RETURNING lobby_id;`

	var hint User
	err := dbpool.QueryRow(ctx, sql, lobbyId, userId).Scan(&hint.LobbyId)
	if err != nil {
		return hint.LobbyId, fmt.Errorf("failed to update user: %w", err)
	}

	return hint.LobbyId, nil
}


type AttackResults struct {
	Duration         time.Duration
	Threads          int
	QueriesPerformed uint64
}

func findAvailableUdpPortAttack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults {
	var queries uint64

	attacker := func(stopAt time.Time) {
		for {
			_, err := findAvailableUdpPort(ctx, dbpool)
			if err != nil { log.Fatal(err) }

			atomic.AddUint64(&queries, 1)

			if time.Now().After(stopAt) {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(threads)

	startAt := time.Now()
	stopAt := startAt.Add(duration)

	for i := 0; i < threads; i++ {
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
}

func findLobbyUsersAttack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults {
	var queries uint64

	attacker := func(stopAt time.Time) {
		for {
			_, err := findLobbyUsers(ctx, dbpool,"f7677b26-c3f9-4028-9fa8-d4732e7e56f6")
			if err != nil { log.Fatal(err) }

			atomic.AddUint64(&queries, 1)

			if time.Now().After(stopAt) {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(threads)

	startAt := time.Now()
	stopAt := startAt.Add(duration)

	for i := 0; i < threads; i++ {
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
}

func updateUsersLobbyIdAttack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults {
	var queries uint64

	attacker := func(stopAt time.Time) {
		for {
			_, err := updateUsersLobbyId(ctx, dbpool,"f7677b26-c3f9-4028-9fa8-d4732e7e56f6", "01d3e2b5-9834-4358-bdbc-5b136f44f166")
			if err != nil { log.Fatal(err) }

			atomic.AddUint64(&queries, 1)

			if time.Now().After(stopAt) {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(threads)

	startAt := time.Now()
	stopAt := startAt.Add(duration)

	for i := 0; i < threads; i++ {
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func main() {
	ctx := context.Background()

	url := "postgres://myuser:secret@localhost:5432/mydb"

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil { log.Fatal(err) }

	aList := []func(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults{findAvailableUdpPortAttack, findLobbyUsersAttack, updateUsersLobbyIdAttack}

	for i := 0; i < 2*3; i++ {
		if (i % 2 == 0) {
			cfg.MaxConns = 1
			cfg.MinConns = 1
		} else {
			cfg.MaxConns = 8
			cfg.MinConns = 8
		}
		dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
		if err != nil { log.Fatal(err) }

		duration := time.Duration(10 * time.Second)
		threads := 1000
		fmt.Println("======== " + getFunctionName(aList[i/2]) + "; cfg conns = " + fmt.Sprint(cfg.MaxConns) + " ========")
		res := aList[i/2](ctx, duration, threads, dbpool)

		fmt.Println("duration:", res.Duration)
		fmt.Println("threads:", res.Threads)
		fmt.Println("queries:", res.QueriesPerformed)
		qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
		fmt.Println("QPS:", qps)

		dbpool.Close()
	}
}