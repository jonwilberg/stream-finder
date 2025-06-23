package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func NewProgressBar(description string, barLength int) *progressbar.ProgressBar {
	return progressbar.NewOptions(barLength,
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stdout, "\n")
		}),
	)
}
