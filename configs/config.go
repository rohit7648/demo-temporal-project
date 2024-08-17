package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	// Config -.
	Config struct {
		Env string
		*App
		*Server
		*Log
		*Temporal
		*DemoTemporalProject
	}

	// Server -.
	Server struct {
		*Grpc
		*Http
	}

	// App
	App struct {
		Name    string
		Version string
	}

	// Grpc -.
	Grpc struct {
		Network string
		Address string
		Timeout *time.Duration
	}

	// Http -.
	Http struct {
		Network string
		Address string
		Timeout *time.Duration
	}

	// Log -.
	Log struct {
		Level string
	}

	Temporal struct {
		BaseUrl       string
		WorkerAddress string
		RetryPolicy   *RetryPolicy
		Workflow      *Workflow
		Activity      *Activity
		Worker        *Worker
	}
	DemoTemporalProject struct {
		BaseUrl       string
		GrpcUrl       string
		UserBatchSize uint64
	}

	RetryPolicy struct {
		InitialInterval    time.Duration
		BackoffCoefficient float64
		MaximumInterval    time.Duration
		MaximumAttempts    int32
	}

	Workflow struct {
		TaskTimeout time.Duration
	}

	Activity struct {
		StartToCloseTimeout time.Duration
	}

	Worker struct {
		MaxConcurrentActivityExecutionSize      int
		MaxConcurrentLocalActivityExecutionSize int
		WorkerActivitiesPerSecond               float64
		TaskQueueActivitiesPerSecond            float64
		MaxConcurrentWorkflowTaskExecutionSize  int
		MaxConcurrentActivityTaskPollers        int
		MaxConcurrentWorkflowTaskPollers        int
	}
)

var Cfg *Config

// Bind all env variables to Viper Keys
func bindEnvVariables() {
	viper.BindEnv("env", "ENV")
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.BindEnv("server.grpc.address", "GRPC_SERVER_ADDRESS")
	viper.BindEnv("server.http.address", "HTTP_SERVER_ADDRESS")
	viper.BindEnv("server.grpc.timeout", "GRPC_SERVER_TIMEOUT")
	viper.BindEnv("server.http.timeout", "HTTP_SERVER_TIMEOUT")
	viper.BindEnv("demoTemporalProject.baseUrl", "DEMO_TEMPORAL_PROJECT_BASE_URL")
	viper.BindEnv("demotemporalproject.grpcUrl", "DEMO_TEMPORAL_PROJECT_GRPC_URL")
	viper.BindEnv("temporal.baseUrl", "TEMPORAL_BASE_URL")
	viper.BindEnv("temporal.schedulerBatchSize", "SCHEDULER_BATCH_SIZE")
	viper.BindEnv("temporal.workerAddress", "TEMPORAL_WORKER_ADDRESS")
}

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	// Load ENV from .env file
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	defaultConfigName := "config"
	// Setting up default configs
	viper.SetConfigName(defaultConfigName)
	viper.AddConfigPath("configs")
	viper.SetConfigType("yaml")
	// Read default configs
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading viper config: %w", err)
	}

	//Resolving env and setting env specific config
	if env := strings.ToLower(os.Getenv("ENV")); strings.Compare(env, "") != 0 {
		envConfigName := defaultConfigName + "." + env
		viper.SetConfigName(envConfigName)
	}

	//merging env configs
	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("error merging env specific config: %w", err)
	}

	bindEnvVariables()

	Cfg = &Config{}

	if err := viper.Unmarshal(Cfg); err != nil {
		return nil, fmt.Errorf("error in converting config to struct - %w", err)
	}
	return Cfg, nil
}
