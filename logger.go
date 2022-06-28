package main

import (
	"os"
)

type TaskLogger struct {
	IsRunning bool
	Config    TaskLogsConfig
	Stdout    *os.File
	Stderr    *os.File
}

func NewLogger(config TaskLogsConfig) *TaskLogger {
	return &TaskLogger{
		IsRunning: false,
		Config:    config,
		Stdout:    nil,
		Stderr:    nil,
	}
}

func (l *TaskLogger) Start() error {
	outFile, err := os.OpenFile(l.Config.OutFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)

	if err != nil {
		return err
	}

	l.Stdout = outFile
	l.Stderr = outFile

	if !l.Config.CombineLogs {
		errFile, err := os.OpenFile(l.Config.ErrorFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)

		if err != nil {
			return err
		}

		l.Stderr = errFile
	}

	l.IsRunning = true

	return nil
}

func (l *TaskLogger) Close() error {
	l.IsRunning = false

	if err := l.Stdout.Close(); err != nil {
		return err
	}

	if !l.Config.CombineLogs {
		return l.Stdout.Close()
	}

	return nil
}
