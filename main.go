package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
	bilog "github.com/cloudfoundry/bosh-cli/logger"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshuifmt "github.com/cloudfoundry/bosh-cli/ui/fmt"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/Altoros/cf-broker-boilerplate/cmd"
)

const mainLogTag = "main"

func main() {
	logger := newLogger()
	defer handlePanic()

	ui := boshui.NewConfUI(logger)
	defer ui.Flush()

	cmdFactory := cmd.NewFactory(boshcmd.NewBasicDeps(ui, logger))

	cmd, err := cmdFactory.New(os.Args[1:])
	if err != nil {
		fail(err, ui, logger)
	}

	logger.Debug(mainLogTag, "Starting CF Service Broker")
	err = cmd.Execute()

	if err != nil {
		fail(err, ui, logger)
	} else {
		success(ui, logger)
	}
}

func newLogger() boshlog.Logger {
	level := boshlog.LevelNone

	logLevelString := os.Getenv("SERVICE_BROKER_LOG_LEVEL")

	if logLevelString != "" {
		var err error
		level, err = boshlog.Levelify(logLevelString)
		if err != nil {
			err = bosherr.WrapError(err, "Invalid SERVICE_BROKER_LOG_LEVEL value")
			logger := boshlog.NewLogger(boshlog.LevelError)
			ui := boshui.NewConsoleUI(logger)
			fail(err, ui, logger)
		}
	}

	return newSignalableLogger(boshlog.NewLogger(level))
}

func newSignalableLogger(logger boshlog.Logger) boshlog.Logger {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	signalableLogger, _ := bilog.NewSignalableLogger(logger, c)
	return signalableLogger
}

func fail(err error, ui boshui.UI, logger boshlog.Logger) {
	if err != nil {
		logger.Error("CLI", err.Error())
		ui.ErrorLinef(boshuifmt.MultilineError(err))
	}
	ui.ErrorLinef("Exit code 1")
	ui.Flush() // todo make sure UI is flushed
	os.Exit(1)
}

func success(ui boshui.UI, logger boshlog.Logger) {
	logger.Debug("CLI", "Succeeded")
	ui.PrintLinef("Succeeded")
}

func handlePanic() {
	panic := recover()

	if panic != nil {
		var msg string

		switch obj := panic.(type) {
		case string:
			msg = obj
		case fmt.Stringer:
			msg = obj.String()
		case error:
			msg = obj.Error()
		default:
			msg = fmt.Sprintf("%#v", obj)
		}

		// Always output to regardless of main logger's level
		logger := boshlog.NewLogger(boshlog.LevelError)
		logger.ErrorWithDetails("CLI", "Panic: %s", msg, debug.Stack())

		ui := boshui.NewConsoleUI(logger)
		fail(nil, ui, logger)
	}
}
