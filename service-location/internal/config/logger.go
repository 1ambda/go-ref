package config

import "go.uber.org/zap"

func GetLogger() *zap.SugaredLogger {
	if IsTestEnv(Spec) || IsLocalEnv(Spec) {
		log, _ := zap.NewDevelopment()
		return log.Sugar()
	}

	log, _ := zap.NewProduction()
	return log.Sugar()
}
