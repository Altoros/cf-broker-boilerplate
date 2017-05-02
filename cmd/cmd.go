package cmd

import (
	"fmt"
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
)

type Cmd struct {
	CommandOpts CommandOpts
	Opts        interface{}

	deps boshcmd.BasicDeps
}

func NewCmd(CommandOpts CommandOpts, opts interface{}, deps boshcmd.BasicDeps) Cmd {
	return Cmd{CommandOpts, opts, deps}
}

type cmdConveniencePanic struct {
	Err error
}

func (c Cmd) Execute() (cmdErr error) {
	// Catch convenience panics from panicIfErr
	defer func() {
		if r := recover(); r != nil {
			if cp, ok := r.(cmdConveniencePanic); ok {
				cmdErr = cp.Err
			} else {
				panic(r)
			}
		}
	}()

	c.configureUI()
	c.configureFS()

	deps := c.deps

	switch opts := c.Opts.(type) {
	case *StartServerOpts:
		return NewStartServerCmd(deps).Run(*opts)
	case *MessageOpts:
		deps.UI.PrintBlock(opts.Message)
		return nil

	default:
		return fmt.Errorf("Unhandled command: %#v", c.Opts)
	}
}

func (c Cmd) configureUI() {
	c.deps.UI.EnableTTY(c.CommandOpts.TTYOpt)

	if !c.CommandOpts.NoColorOpt {
		c.deps.UI.EnableColor()
	}

	if c.CommandOpts.JSONOpt {
		c.deps.UI.EnableJSON()
	}

	if c.CommandOpts.NonInteractiveOpt {
		c.deps.UI.EnableNonInteractive()
	}
}

func (c Cmd) configureFS() {
	tmpDirPath, err := c.deps.FS.ExpandPath("~/.boilerplate-service-broker/tmp")
	c.panicIfErr(err)

	err = c.deps.FS.ChangeTempRoot(tmpDirPath)
	c.panicIfErr(err)
}

func (c Cmd) panicIfErr(err error) {
	if err != nil {
		panic(cmdConveniencePanic{err})
	}
}
