package main

import (
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	AmqpURL      *url.URL // AMQP URL (password comes from the token)
	AmqpExchange string   // Exchange to shovel messages
	AmqpToken    string   // File location of the token
	ListenPort   int
	ListenIp     string
	DestUdp      []string
	Debug        bool
	Verify       bool
}

func (c *Config) ReadConfig() {
	viper.SetConfigName("config")                            // name of config file (without extension)
	viper.SetConfigType("yaml")                              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/xrootd-monitoring-shoveler/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.xrootd-monitoring-shoveler") // call multiple times to add many search paths
	viper.AddConfigPath(".")                                 // optionally look for config in the working directory
	viper.AddConfigPath("config/")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Warningln("Unable to read in config file, will check environment for configuration:", err)
	}
	viper.SetEnvPrefix("SHOVELER")

	// Autmatically look to the ENV for all "Gets"
	viper.AutomaticEnv()
	// Look for environment variables with underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("amqp.exchange", "shoveled-xrd")
	viper.SetDefault("amqp.token_location", "/etc/xrootd-monitoring-shoveler/token")

	// Get the AMQP URL
	c.AmqpURL, err = url.Parse(viper.GetString("amqp.url"))
	if err != nil {
		panic(fmt.Errorf("Fatal error parsing AMQP URL: %s \n", err))
	}
	log.Debugln("AMQP URL:", c.AmqpURL.String())

	// Get the AMQP Exchange
	c.AmqpExchange = viper.GetString("amqp.exchange")
	log.Debugln("AMQP Exchange:", c.AmqpExchange)

	// Get the Token location
	c.AmqpToken = viper.GetString("amqp.token_location")
	log.Debugln("AMQP Token location:", c.AmqpToken)

	// Get the UDP listening parameters
	viper.SetDefault("listen.port", 9993)
	c.ListenPort = viper.GetInt("listen.port")
	c.ListenIp = viper.GetString("listen.ip")

	c.DestUdp = viper.GetStringSlice("outputs.destinations")

	c.Debug = viper.GetBool("debug")

	viper.SetDefault("verify", true)
	c.Verify = viper.GetBool("verify")

	// Metrics defaults
	viper.SetDefault("metrics.enable", true)
	viper.SetDefault("metrics.port", 8000)

}
