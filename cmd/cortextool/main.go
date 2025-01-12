package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/cortexproject/cortex-tools/pkg/commands"
	"github.com/cortexproject/cortex-tools/pkg/version"
)

var (
	ruleCommand              commands.RuleCommand
	alertCommand             commands.AlertCommand
	alertmanagerCommand      commands.AlertmanagerCommand
	logConfig                commands.LoggerConfig
	pushGateway              commands.PushGatewayConfig
	loadgenCommand           commands.LoadgenCommand
	remoteReadCommand        commands.RemoteReadCommand
	aclCommand               commands.AccessControlCommand
	analyseCommand           commands.AnalyseCommand
	bucketValidateCommand    commands.BucketValidationCommand
	overridesExporterCommand = commands.NewOverridesExporterCommand()
)

func main() {
	app := kingpin.New("cortextool", "A command-line tool to manage cortex.")
	logConfig.Register(app)
	alertCommand.Register(app)
	alertmanagerCommand.Register(app)
	ruleCommand.Register(app)
	pushGateway.Register(app)
	loadgenCommand.Register(app)
	remoteReadCommand.Register(app)
	overridesExporterCommand.Register(app)
	aclCommand.Register(app)
	analyseCommand.Register(app)
	bucketValidateCommand.Register(app)

	app.Command("version", "Get the version of the cortextool CLI").Action(func(k *kingpin.ParseContext) error {
		fmt.Print(version.Template)
		version.CheckLatest()

		return nil
	})

	kingpin.MustParse(app.Parse(os.Args[1:]))

	pushGateway.Stop()
}
