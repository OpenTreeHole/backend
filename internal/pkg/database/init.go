package database

type Config struct {
	GormModels         []any
	EnableGorm         bool
	EnableCache        bool
	EnableSearchEngine bool
}

func Init(config Config) {
	if config.EnableGorm {
		initGorm(config.GormModels...)
	}

	if config.EnableCache {
		initCache()
	}

	if config.EnableSearchEngine {
		initSearchEngine()
	}
}
