package prepare

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"riverqueue.com/riverui"
)

func MustNewRiver(ctx context.Context, dbPool *pgxpool.Pool) *riverui.Handler {
	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: mustNewWorkers(),
	})
	if err != nil {
		log.Fatalf("riverClient create err: %s", err)
	}

	if err = riverClient.Start(ctx); err != nil {
		log.Fatalf("riverClient err: %s", err)
	}

	endpoints := riverui.NewEndpoints(riverClient, nil)

	handler, err := riverui.NewHandler(&riverui.HandlerOpts{
		Endpoints: endpoints,
		Logger:    slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Prefix:    "/riverui",
	})
	if err != nil {
		log.Fatalf("riverui.NewHandler %s", err)
	}

	// запускаем фоновые задачи UI
	err = handler.Start(ctx)
	if err != nil {
		log.Fatalf("riverui.Start %s", err)
	}

	return handler

}

func mustNewWorkers() *river.Workers {
	workers := river.NewWorkers()
	if err := river.AddWorkerSafely(workers, &SortWorker{}); err != nil {
		panic("handle this error")

	}
	return workers
}

type SortArgs struct {
	// Strings is a slice of strings to sort.
	Strings []string `json:"strings"`
}

func (SortArgs) Kind() string { return "sort" }

type SortWorker struct {
	// An embedded WorkerDefaults sets up default methods to fulfill the rest of
	// the Worker interface:
	river.WorkerDefaults[SortArgs]
}

func (w *SortWorker) Work(ctx context.Context, job *river.Job[SortArgs]) error {
	sort.Strings(job.Args.Strings)
	fmt.Printf("Sorted strings: %+v\n", job.Args.Strings)
	return nil
}
