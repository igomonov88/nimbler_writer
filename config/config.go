package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Service      Service      `yaml:"web"`
	Zipkin       Zipkin       `yaml:"zipkin"`
	Database     Database     `yaml:"database"`
	S3           S3           `yaml:"aws_s3"`
	KeyGenerator KeyGenerator `yaml:"key_generator"`
}

type Service struct {
	APIHost         string        `yaml:"api_host"`
	DebugHost       string        `yaml:"debug_host"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type Zipkin struct {
	LocalEndpoint string  `yaml:"local_endpoint"`
	ReporterURI   string  `yaml:"reporter_uri"`
	ServiceName   string  `yaml:"service_name"`
	Probability   float64 `yaml:"probability"`
}

type Database struct {
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	Name       string `yaml:"name"`
	DisableTLS bool   `yaml:"disable_tls"`
}

type S3 struct {
	AccessKeyID string `yaml:"access_key_id"`
	SecretKey   string `yaml:"secret_key"`
	BucketName  string `yaml:"bucket_name"`
}

type KeyGenerator struct {
	APIHost         string        `yaml:"api_host"`
	DebugHost       string        `yaml:"debug_host"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

func Parse(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "reading file path")
	}
	defer file.Close()
	var cfg Config
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, errors.Wrap(err, "decoding config file")
	}
	return &cfg, nil
}
