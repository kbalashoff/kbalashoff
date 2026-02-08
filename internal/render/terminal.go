package render

import (
	"fmt"
	"os"
	"strings"

	"falltype/internal/engine"
)

type Renderer struct {
	Width     int
	Height    int
	TopRow    int
	BottomRow int
}

func NewRenderer(width, height int) Renderer {
	if width < 40 {
		width = 40
	}
	if height < 20 {
		height = 20
	}
	return Renderer{
		Width:     width,
		Height:    height,
		TopRow:    3,
		BottomRow: height - 5,
	}
}

func (r Renderer) Clear() {
	fmt.Fprint(os.Stdout, "\033[2J\033[H")
}

func (r Renderer) Draw(state engine.FrameState) {
	grid := make([][]rune, r.Height)
	for y := range grid {
		grid[y] = []rune(strings.Repeat(" ", r.Width))
	}

	title := fmt.Sprintf(" falltype | %s | lesson %d/%d ", state.Lesson.Name, state.LessonIndex+1, state.TotalLessons)
	for i, ch := range title {
		if i+1 < r.Width-1 {
			grid[0][i+1] = ch
		}
	}

	meta := fmt.Sprintf("Spawned:%d  Hits:%d  Misses:%d  Wrong:%d  Accuracy:%.1f%%  (ESC/Q to quit)",
		state.Stats.TotalSpawned,
		state.Stats.CorrectHits,
		state.Stats.Misses,
		state.Stats.WrongKeys,
		state.Stats.Accuracy()*100,
	)
	for i, ch := range meta {
		if i+1 < r.Width-1 {
			grid[1][i+1] = ch
		}
	}

	for x := 0; x < r.Width; x++ {
		grid[r.TopRow-1][x] = '─'
		grid[r.BottomRow][x] = '─'
	}

	for _, f := range state.Falling {
		if !f.Active || f.Row < 0 || f.Row >= r.Height {
			continue
		}
		if f.Column >= 0 && f.Column < r.Width {
			grid[f.Row][f.Column] = f.Char
		}
	}

	for _, line := range grid {
		fmt.Fprintln(os.Stdout, string(line))
	}

	if state.StatusLine != "" {
		fmt.Fprintln(os.Stdout, state.StatusLine)
	}
}

func (r Renderer) DrawLessonResult(stats engine.Stats, lesson engine.Lesson, passed bool) {
	r.Clear()
	fmt.Printf("\n%s\n", lesson.Name)
	fmt.Println(strings.Repeat("=", len(lesson.Name)))
	fmt.Printf("Spawned: %d\n", stats.TotalSpawned)
	fmt.Printf("Hits:    %d\n", stats.CorrectHits)
	fmt.Printf("Misses:  %d\n", stats.Misses)
	fmt.Printf("Wrong:   %d\n", stats.WrongKeys)
	fmt.Printf("Accuracy %.1f%% (required %.1f%%)\n", stats.Accuracy()*100, lesson.RequiredAccuracy*100)
	if passed {
		fmt.Println("Result: PASS")
	} else {
		fmt.Println("Result: RETRY")
	}
	fmt.Println("Press any key to continue...")
}

func (r Renderer) DrawFinalSummary(all []engine.Stats) {
	r.Clear()
	fmt.Println("falltype complete")
	fmt.Println("=================")
	for i, s := range all {
		fmt.Printf("Lesson %d -> Accuracy %.1f%% | Hits %d | Misses %d | Wrong %d\n",
			i+1,
			s.Accuracy()*100,
			s.CorrectHits,
			s.Misses,
			s.WrongKeys,
		)
	}
}

func (r Renderer) Bottom() int { return r.BottomRow }
func (r Renderer) Top() int    { return r.TopRow }
