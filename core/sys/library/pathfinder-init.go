package library

import (
	"log/slog"
	"os"
)

func (l *Library) initPFAncestries(f string, logger *slog.Logger) {
	content, err := loadAncestriesPF(f)
	if err != nil {
		logger.Error("error from ancestries", "error", err.Error())
		os.Exit(2)
	}
	for _, v := range content {
		if v.Name != "" {
			l.PFAncestry.Append(v)
		}
	}
}

func (l *Library) initPFBackgrounds(f string, logger *slog.Logger) {
	content, err := loadBackgroundsPF(f)
	if err != nil {
		logger.Error("error from backgrounds", "error", err.Error())
		os.Exit(2)
	}
	for _, v := range content {
		if v.Name != "" {
			l.PFBackground.Append(v)
		}
	}
}

func (l *Library) initPFClasses(f string, logger *slog.Logger) {
	content, err := loadClassesPF(f)
	if err != nil {
		logger.Error("error from classes", "error", err.Error())
		os.Exit(2)
	}
	for _, v := range content {
		if v.Name != "" {
			l.PFClass.Append(v)
		}
	}
}
