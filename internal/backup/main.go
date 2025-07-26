package main

import (
	"context"
	"os"
	"os/signal"
	"quotobot/pkg/config"
	"quotobot/pkg/logger"

	"github.com/go-co-op/gocron/v2"
)

type Backup struct {
	Logger    *logger.Logger
	Config    *config.Config
	Scheduler gocron.Scheduler
}

func NewBackup() *Backup {
	l := logger.NewLogger()
	c := config.LoadConfig(l)
	s, err := gocron.NewScheduler()
	if err != nil {
		l.Error.Fatalf("Failed to create scheduler: %v\n", err)
		return nil
	}

	return &Backup{
		Logger:    l,
		Config:    c,
		Scheduler: s,
	}
}

func (bak *Backup) Start(ctx context.Context) {
	bak.Logger.Info.Println("Backup service started")

	j, err := bak.Scheduler.NewJob(
		gocron.CronJob(bak.Config.Backup.Cron, false),
		gocron.NewTask(
			bak.performBackup,
			ctx,
		),
		gocron.JobOption(gocron.WithStartImmediately()),
	)
	if err != nil {
		bak.Logger.Error.Fatalf("Failed to create backup job: %v\n", err)
		return
	}
	bak.Logger.Info.Printf("Scheduled backup job: %s\n", j.ID())

	bak.Scheduler.Start()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	backup := NewBackup()
	backup.Start(ctx)

	<-ctx.Done()
	backup.Scheduler.Shutdown()
	backup.Logger.Info.Println("Backup service stopped")
	os.Exit(0)
}
