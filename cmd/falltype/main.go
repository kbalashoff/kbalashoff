package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"falltype/internal/engine"
	"falltype/internal/input"
	"falltype/internal/model"
	"falltype/internal/render"
)

type Config struct {
	InitialSpeedMs         int `json:"initial_speed_ms"`
	SpeedDecreasePerLesson int `json:"speed_decrease_per_lesson"`
	MaxSimultaneousLimit   int `json:"max_simultaneous_limit"`
}

func main() {
	cfgPath := flag.String("config", "", "optional path to JSON config")
	flag.Parse()

	lessons := engine.DefaultLessons()
	if *cfgPath != "" {
		cfg, err := loadConfig(*cfgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
			os.Exit(1)
		}
		overrideLessons(&lessons, cfg)
	}

	r := render.NewRenderer(96, 28)
	columns := model.BuildColumnMap(4)
	e := engine.New(columns, r.Top(), r.Bottom())

	kbd := input.NewKeyboard()
	in, errCh, err := kbd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to enable raw mode: %v\n", err)
		os.Exit(1)
	}
	defer kbd.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	quit := make(chan struct{})
	go func() {
		select {
		case <-sig:
			close(quit)
		case err := <-errCh:
			fmt.Fprintf(os.Stderr, "input error: %v\n", err)
			close(quit)
		}
	}()

	results := make([]engine.Stats, 0, len(lessons))
	r.Clear()
	for i := 0; i < len(lessons); i++ {
		lesson := lessons[i]
		stats, passed, stopped := e.RunLesson(lesson, i, len(lessons), in, quit, func(fs engine.FrameState) {
			r.Clear()
			r.Draw(fs)
		}, 25*time.Millisecond)
		if stopped {
			break
		}
		results = append(results, stats)
		r.DrawLessonResult(stats, lesson, passed)
		<-in
		if !passed {
			i-- // retry same lesson
		}
	}

	r.DrawFinalSummary(results)
	fmt.Println("\nPress any key to exit...")
	<-in
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func overrideLessons(lessons *[]engine.Lesson, cfg Config) {
	for i := range *lessons {
		if cfg.InitialSpeedMs > 0 {
			speed := cfg.InitialSpeedMs - i*cfg.SpeedDecreasePerLesson
			if speed < 120 {
				speed = 120
			}
			(*lessons)[i].FallSpeed = time.Duration(speed) * time.Millisecond
		}
		if cfg.MaxSimultaneousLimit > 0 && (*lessons)[i].MaxSimultaneous > cfg.MaxSimultaneousLimit {
			(*lessons)[i].MaxSimultaneous = cfg.MaxSimultaneousLimit
		}
	}
}
