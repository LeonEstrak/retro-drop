package constants

const (
	LOG_LEVEL                  = "debug"
	MYRIENT_ERISTA_ME_BASE_URL = "https://myrient.erista.me"
	DB_PATH                    = "./database/retro-drop.db"
)

var (
	SYSTEMS_TO_ERISTA_MAPPING = map[string]string{
		"gba":  "/files/No-Intro/Nintendo%20-%20Game%20Boy%20Advance/",
		"snes": "/files/No-Intro/Nintendo%20-%20Super%20Nintendo%20Entertainment%20System/",
	}
)
