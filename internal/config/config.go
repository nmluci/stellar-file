package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName    string          `json:"serviceName"`
	ServiceAddress string          `json:"servicePort"`
	ServiceID      string          `json:"serviceID"`
	RPCAddress     string          `json:"rpcAddress"`
	TrustedService map[string]bool `json:"trustedService"`
	Environment    Environment     `json:"environment"`

	WorkerConfig  WorkerConfig     `json:"workerConfig"`
	MariaDBConfig MariaDBConfig    `json:"mariaDBConfig"`
	MongoDBConfig MongoDBConfig    `json:"mongoDBConfig"`
	RedisConfig   RedisConfig      `json:"redisConfig"`
	StellarConfig StellarRPCConfig `json:"stellarRPCConfig"`
}

const logTagConfig = "[Init Config]"

var config *Config

func Init() {
	godotenv.Load("conf/.env")

	conf := Config{
		ServiceName:    os.Getenv("SERVICE_NAME"),
		ServiceAddress: os.Getenv("SERVICE_ADDR"),
		ServiceID:      os.Getenv("SERVICE_ID"),
		RPCAddress:     os.Getenv("GPRC_ADDR"),
		MariaDBConfig: MariaDBConfig{
			Address:  os.Getenv("MARIADB_ADDRESS"),
			Username: os.Getenv("MARIADB_USERNAME"),
			Password: os.Getenv("MARIADB_PASSWORD"),
			DBName:   os.Getenv("MARIADB_DBNAME"),
		},
		MongoDBConfig: MongoDBConfig{
			Address:  os.Getenv("MONGODB_ADDRESS"),
			Username: os.Getenv("MONGODB_USERNAME"),
			Password: os.Getenv("MONGODB_PASSWORD"),
			DBName:   os.Getenv("MONGODB_DBNAME"),
		},
		RedisConfig: RedisConfig{
			Address:  os.Getenv("REDIS_ADDRESS"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		StellarConfig: StellarRPCConfig{
			AuthKey: os.Getenv("STELLAR_AUTH_KEY"),
		},
		WorkerConfig: WorkerConfig{
			Arc: ArchivalWorkerConfig{
				DefaultDir: os.Getenv("ARC_DEFAULT_DIR"),
			},
			Down: DownloaderWorkerConfig{
				DefaultDir: os.Getenv("DOWN_DEFAULT_DIR"),
			},
		},
	}

	if conf.ServiceName == "" {
		log.Fatalf("%s service name should not be empty", logTagConfig)
	}

	if conf.ServiceAddress == "" {
		log.Fatalf("%s service port should not be empty", logTagConfig)
	}

	if conf.MariaDBConfig.Address == "" || conf.MariaDBConfig.DBName == "" {
		log.Fatalf("%s address and db name cannot be empty", logTagConfig)
	}

	envString := os.Getenv("ENVIRONMENT")
	if envString != "dev" && envString != "prod" {
		log.Fatalf("%s environment must be either dev or prod, found: %s", logTagConfig, envString)
	}

	conf.Environment = Environment(envString)

	conf.TrustedService = map[string]bool{conf.ServiceID: true}
	if trusted := os.Getenv("TRUSTED_SERVICES"); trusted == "" {
		conf.TrustedService["STELLAR_HENTAI"] = true
	} else {
		for _, svc := range strings.Split(trusted, ",") {
			if _, ok := conf.TrustedService[svc]; !ok {
				conf.TrustedService[svc] = true
			}
		}
	}

	config = &conf
}

func Get() (conf *Config) {
	conf = config
	return
}
