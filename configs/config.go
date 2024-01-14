package configs

import "github.com/spf13/viper"

type Config struct {
	Env        string
	Port       int
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
	DbSsl      string
	RedisHost  string
	RedisPort  int
}

func Load() Config {
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()
	conf := Config{}
	conf.Env = viper.GetString("ENV")
	conf.Port = viper.GetInt("PORT")
	conf.DbHost = viper.GetString("POSTGRES_HOST")
	conf.DbPort = viper.GetInt("POSTGRES_PORT")
	conf.DbUser = viper.GetString("POSTGRES_USER")
	conf.DbPassword = viper.GetString("POSTGRES_PASSWORD")
	conf.DbName = viper.GetString("POSTGRES_DB")
	conf.DbSsl = viper.GetString("POSTGRES_SSL")
	conf.RedisHost = viper.GetString("REDIS_HOST")
	conf.RedisPort = viper.GetInt("REDIS_PORT")
	return conf
}
