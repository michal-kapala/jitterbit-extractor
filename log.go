package main

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) logError(err error) {
	runtime.LogError(a.ctx, err.Error())
}

func (a *App) logWarning(msg string) {
	runtime.LogPrint(a.ctx, msg)
}
