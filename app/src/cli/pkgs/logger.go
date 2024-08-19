package pkg

import (
	zerolog "github.com/rs/zerolog"
	dmlogger "github.com/vapusdata-ecosystem/vapusai-studio/core/pkgs/logger"
)

var DmLogger zerolog.Logger

func InitLogger(debugMode bool) {
	DmLogger = dmlogger.NewZeroLogger(debugMode, true, false, "", "caller", "level", "fields", "time")
}
