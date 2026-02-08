package model

import "time"

type FallingLetter struct {
	Char      rune
	Column    int
	Row       int
	Speed     time.Duration
	lastTick  time.Time
	Active    bool
	SpawnedAt time.Time
}

func NewFallingLetter(char rune, column, row int, speed time.Duration, now time.Time) FallingLetter {
	return FallingLetter{
		Char:      char,
		Column:    column,
		Row:       row,
		Speed:     speed,
		lastTick:  now,
		Active:    true,
		SpawnedAt: now,
	}
}

func (f *FallingLetter) Step(now time.Time) bool {
	if !f.Active {
		return false
	}
	if now.Sub(f.lastTick) >= f.Speed {
		f.Row++
		f.lastTick = now
		return true
	}
	return false
}
