package library

import "log/slog"

const (
	D10HM            string = "D10HomeMade" // D10 based on World of Darkness
	PFD20            string = "Pathfinder"  // Pathfinder d20
	PFD20Ancestries  string = "PFD20Ancestries"
	PFD20Backgrounds string = "PFD20Backgrounds"
	PFD20Classes     string = "PFD20Classes"
)

func (l *Library) ImportRPGLibrary(name string, logger *slog.Logger, config map[string]string) {
	switch name {
	case D10HM:
		l.initDescriptions(config[D10HM], logger)
	case PFD20:
		l.initDescriptions(config[PFD20], logger)
		l.initPFAncestries(config[PFD20Ancestries], logger)
		l.initPFBackgrounds(config[PFD20Backgrounds], logger)
		l.initPFClasses(config[PFD20Classes], logger)
	}
}
