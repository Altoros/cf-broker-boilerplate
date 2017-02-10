package cmd

type CommandOpts struct {
	// -----> Global options
	Version               func() error `long:"version" short:"v" description:"Show Service Brokere version"`
	Port                  int          `long:"port" short:"p" description:"Port of Service Broker" default:"8080" env:"PORT"`
	ConfigPath            string       `long:"config" short:"c" description:"Path to config with plans"`
	ServiceID             string       `long:"service-broker-id" descrition:"service broker service id" default:"service-broker" env:"SERVICE_BROKER_ID"`
	Name                  string       `long:"service-broker-name" descrition:"service broker name" default:"service-broker" env:"SERVICE_BROKER_NAME"`
	Description           string       `long:"service-broker-description" descrition:"service broker name" default:"Service for running destructive tests on a app" env:"SERVICE_BROKER_DESCRIPTION"`
	ServiceBrokerUsername string       `long:"service-broker-username" description:"broker-username" default:"service-broker"`
	ServiceBrokerPassword string       `long:"service-broker-password" description:"broker-password" default:"c1oudc0w"`
	Help                  bool         `long:"help" short:"h" description:"Show this help message"`
}
