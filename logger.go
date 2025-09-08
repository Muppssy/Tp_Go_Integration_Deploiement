package main

import "gitlab.com/ggpack/logchain-go"

func initLogging() logchain.Logger {
	params := logchain.Params{
		"template":  "{{.timestamp}} " + version + " {{.levelLetter}} {{.fileLine}} {{.msg}}",
		"verbosity": 3,
	}
	chainer := logchain.NewLogChainer(params)
	return chainer.InitLogging()
}

var Logger = initLogging()
