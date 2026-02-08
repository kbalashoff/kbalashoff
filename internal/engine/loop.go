package engine

import (
	"math/rand"
	"time"

	"falltype/internal/model"
)

type Stats struct {
	TotalSpawned int
	CorrectHits  int
	Misses       int
	WrongKeys    int
}

func (s Stats) Accuracy() float64 {
	total := s.CorrectHits + s.Misses + s.WrongKeys
	if total == 0 {
		return 0
	}
	return float64(s.CorrectHits) / float64(total)
}

type FrameState struct {
	Lesson       Lesson
	LessonIndex  int
	TotalLessons int
	Falling      []model.FallingLetter
	Stats        Stats
	StatusLine   string
}

type Engine struct {
	columns map[rune]int
	topRow  int
	botRow  int
	rng     *rand.Rand
}

func New(columns map[rune]int, topRow, bottomRow int) *Engine {
	return &Engine{
		columns: columns,
		topRow:  topRow,
		botRow:  bottomRow,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (e *Engine) RunLesson(lesson Lesson, lessonIndex, totalLessons int, input <-chan rune, quit <-chan struct{}, frame func(FrameState), tick time.Duration) (Stats, bool, bool) {
	stats := Stats{}
	falling := make([]model.FallingLetter, 0, lesson.MaxSimultaneous)
	lastSpawn := time.Time{}
	status := ""

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for {
		now := time.Now()
		select {
		case <-quit:
			return stats, false, true
		case ch := <-input:
			if ch == 27 || ch == 'q' {
				return stats, false, true
			}
			if !model.IsLetterAllowed(ch) {
				continue
			}
			matched := false
			for i := range falling {
				if falling[i].Active && falling[i].Char == ch {
					falling[i].Active = false
					stats.CorrectHits++
					matched = true
					break
				}
			}
			if !matched {
				stats.WrongKeys++
				status = "✗ wrong key"
			} else {
				status = ""
			}
		case <-ticker.C:
			if len(activeOnly(falling)) < lesson.MaxSimultaneous && (lastSpawn.IsZero() || now.Sub(lastSpawn) >= lesson.SpawnDelay) && stats.TotalSpawned < lesson.RequiredAttempts {
				falling = append(falling, e.spawn(lesson, now))
				stats.TotalSpawned++
				lastSpawn = now
			}

			for i := range falling {
				if !falling[i].Active {
					continue
				}
				if falling[i].Step(now) && falling[i].Row >= e.botRow {
					falling[i].Active = false
					stats.Misses++
					status = "⚠ miss"
				}
			}

			falling = activeOnly(falling)
			frame(FrameState{Lesson: lesson, LessonIndex: lessonIndex, TotalLessons: totalLessons, Falling: falling, Stats: stats, StatusLine: status})
		}

		if stats.TotalSpawned >= lesson.RequiredAttempts && len(falling) == 0 {
			passed := stats.Accuracy() >= lesson.RequiredAccuracy
			return stats, passed, false
		}
	}
}

func (e *Engine) spawn(lesson Lesson, now time.Time) model.FallingLetter {
	char := lesson.ActiveLetters[e.rng.Intn(len(lesson.ActiveLetters))]
	col := e.columns[char]
	return model.NewFallingLetter(char, col, e.topRow, lesson.FallSpeed, now)
}

func activeOnly(in []model.FallingLetter) []model.FallingLetter {
	out := in[:0]
	for _, f := range in {
		if f.Active {
			out = append(out, f)
		}
	}
	return out
}
