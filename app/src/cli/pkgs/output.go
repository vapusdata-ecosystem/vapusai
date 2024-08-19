package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/list"
	table "github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rs/zerolog"
)

var defaultTableStyle = table.StyleRounded

func NewTableWritter() table.Writer {
	tw := table.NewWriter()
	tw.SetStyle(defaultTableStyle)
	tw.SetOutputMirror(os.Stdout)
	return tw
}

func NewListWritter(items []interface{}, style list.Style) list.Writer {
	lw := list.NewWriter()
	lw.SetStyle(style)
	lw.AppendItems(items)
	return lw
}

func GetSpinner(charSet int) *spinner.Spinner {
	if charSet == 0 {
		charSet = 43
	}
	s := spinner.New(spinner.CharSets[charSet], 100*time.Millisecond)
	s.Color("white")
	s.Suffix = " :Please wait..."
	return s
}

func LogTitles(mess string, logger zerolog.Logger) {
	xx := text.FormatTitle.Apply(mess)
	xx = text.Underline.Sprintf(xx)
	logger.Info().Msg("\n")
	logger.Info().Msg(xx)
}
func LogTitlesf(mess string, logger zerolog.Logger, args ...interface{}) {
	xx := text.FormatTitle.Apply(mess)
	xx = text.Underline.Sprintf(xx)
	logger.Info().Msg("\n")
	for _, arg := range args {
		xx = fmt.Sprintf(xx+" : %v", arg)
	}
	logger.Info().Msg(xx)
}

func LogDescription(mess string, logger zerolog.Logger, args ...interface{}) {
	xx := text.FormatDefault.Apply(mess)
	logger.Info().Msg("\n")
	logger.Info().Msg(fmt.Sprintf(xx, args...))
	logger.Info().Msg("\n")
}

func UnixTransformer(tt int64) string {
	t := time.Unix(tt, 0)
	// tf := text.NewUnixTimeTransformer("dd-MM-yyyy", time.Local)
	return t.Format("2006-01-02 15:04:05")
}
