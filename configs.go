package rlog

import (
	"sync"
)

const (
	defaultLogLevel         = "debug" //  "debug | info | warn | error"
	defaultSlackrusLogLevel = "warn"
	defaultLogChannelSize   = 2000
)

var logsChannel chan [][]byte
var logsWaitGroup = new(sync.WaitGroup)
var logsDone chan bool

type SlackrusCfg struct {
	Enabled  bool
	Endpoint string // Endpoint for your Slack hook
	LogLevel string //  "debug | info | warn | error"
}

type LogConfig struct {
	EnvPrefix   string
	Formatter   string // "text" | "json"
	LogLevel    string //  "debug | info | warn | error"
	LogChanSize int
	SlackrusCfg SlackrusCfg
}
