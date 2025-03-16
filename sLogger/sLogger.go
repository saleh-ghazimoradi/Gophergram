package sLogger

import (
	"log/slog"
	"os"
)

var SLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
