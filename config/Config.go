package config

type Config struct {
	Server struct {
		Port        int    `mapstructure:"port"`
		Environment string `mapstructure:"environment"`
	} `mapstructure:"server"`
	Db struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Address  string `mapstructure:"address"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"db"`
}
