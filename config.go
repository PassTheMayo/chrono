package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Tasks   []TaskConfig `yaml:"tasks"`
	LogFile string       `yaml:"log_file"`
}

type TaskConfig struct {
	Label     string             `yaml:"label"`
	SilentRun bool               `yaml:"silent_run"`
	Interval  TaskIntervalConfig `yaml:"interval"`
	Execute   TaskRunnerConfig   `yaml:"execute"`
	Logs      TaskLogsConfig     `yaml:"logs"`
}

type TaskIntervalConfig struct {
	Type        string `yaml:"type"`
	Value       string `yaml:"value"`
	StartAtZero bool   `yaml:"start_at_zero"`
}

type TaskRunnerConfig struct {
	Type    string            `yaml:"type"`
	Process string            `yaml:"process"`
	Args    []string          `yaml:"args"`
	Env     map[string]string `yaml:"env"`
}

type TaskLogsConfig struct {
	Enable      bool   `yaml:"enable"`
	OutFile     string `yaml:"out_file"`
	ErrorFile   string `yaml:"error_file"`
	CombineLogs bool   `yaml:"combine_logs"`
}

func (c *Config) ReadFile(file string) error {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}
