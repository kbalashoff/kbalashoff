package engine

import "time"

type Lesson struct {
	Name             string
	ActiveLetters    []rune
	MaxSimultaneous  int
	FallSpeed        time.Duration
	RequiredAccuracy float64
	RequiredAttempts int
	SpawnDelay       time.Duration
}

func DefaultLessons() []Lesson {
	return []Lesson{
		{
			Name:             "Lesson 1: g h",
			ActiveLetters:    []rune{'g', 'h'},
			MaxSimultaneous:  1,
			FallSpeed:        750 * time.Millisecond,
			RequiredAccuracy: 0.85,
			RequiredAttempts: 20,
			SpawnDelay:       1100 * time.Millisecond,
		},
		{
			Name:             "Lesson 2: f g h j",
			ActiveLetters:    []rune{'f', 'g', 'h', 'j'},
			MaxSimultaneous:  1,
			FallSpeed:        620 * time.Millisecond,
			RequiredAccuracy: 0.85,
			RequiredAttempts: 24,
			SpawnDelay:       1000 * time.Millisecond,
		},
		{
			Name:             "Lesson 3: home row core",
			ActiveLetters:    []rune{'d', 'f', 'g', 'h', 'j', 'k'},
			MaxSimultaneous:  2,
			FallSpeed:        520 * time.Millisecond,
			RequiredAccuracy: 0.86,
			RequiredAttempts: 28,
			SpawnDelay:       850 * time.Millisecond,
		},
		{
			Name:             "Lesson 4: home row full",
			ActiveLetters:    []rune{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'ö', 'ä'},
			MaxSimultaneous:  3,
			FallSpeed:        430 * time.Millisecond,
			RequiredAccuracy: 0.87,
			RequiredAttempts: 36,
			SpawnDelay:       700 * time.Millisecond,
		},
		{
			Name:             "Lesson 5: full alphabet",
			ActiveLetters:    []rune{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', 'å', 'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'ö', 'ä', 'z', 'x', 'c', 'v', 'b', 'n', 'm'},
			MaxSimultaneous:  5,
			FallSpeed:        300 * time.Millisecond,
			RequiredAccuracy: 0.88,
			RequiredAttempts: 50,
			SpawnDelay:       450 * time.Millisecond,
		},
	}
}
