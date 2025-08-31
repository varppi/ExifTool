package progress

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

type ProgressBar struct {
	Max      int
	Disabled bool
	progress int
	stop     bool
}

func (p *ProgressBar) Progress() {
	p.progress++
}

func (p *ProgressBar) StartLoading() {
	if p.Disabled {
		return
	}
	div := float64(p.Max / 25)
	if p.Max < 25 {
		div = float64(p.Max)
	}
	go func() {
		for !p.stop {
			pr := float64(p.progress)
			ma := float64(p.Max)
			blocksAmount := int(math.Floor(pr/div)) + 1
			emptyAmount := int(math.Floor((ma-pr)/div)) + 1
			fmt.Fprintf(os.Stderr, " %s%s %d/%d\r", strings.Repeat("▰", blocksAmount), strings.Repeat(" ", emptyAmount), p.progress, p.Max)
			time.Sleep(10 * time.Millisecond)
		}
		fmt.Println(strings.Repeat(" ", 75))
		p.stop = false
	}()
}

func (p *ProgressBar) StopLoading() {
	p.stop = true
	time.Sleep(15 * time.Millisecond)
}
