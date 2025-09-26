package lg

import (
	"context"

	"github.com/Lexographics/logar"
)

var (
	UserLogs   = "user-logs"
	SystemLogs = "system-logs"
	AdminLogs  = "admin-logs"
)

var instance logar.App

func Set(l logar.App) {
	instance = l
}

func Get() logar.App {
	return instance
}

func Logger() logar.Logger {
	return instance.GetLogger()
}

func WithContext(ctx context.Context) logar.Logger {
	return instance.GetLogger().WithContext(ctx)
}

func ActionManager() logar.ActionManager {
	return instance.GetActionManager()
}

func Analytics() logar.Analytics {
	return instance.GetAnalytics()
}
