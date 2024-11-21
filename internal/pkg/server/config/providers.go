package config

type OvhProviders struct {
	Endpoint string `mapstructure:"endpoint"`
	ApplicationKey string `mapstructure:"application_key"`
	ApplicationSecret string `mapstructure:"application_secret"`
	ConsumerKey string `mapstructure:"consumer_key"`
}
