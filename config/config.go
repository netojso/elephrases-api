package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Valores usados quando a variável correspondente não está no .env.
const (
	defaultMaxOpenConns       = 25
	defaultMaxIdleConns       = 10
	defaultConnMaxLifetimeMin = 30
	defaultConnMaxIdleTimeMin = 5
	pingTimeout               = 5 * time.Second
	slowQueryThreshold        = 200 * time.Millisecond
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	DBMaxOpenConns         int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns         int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBConnMaxLifetimeMin   int    `mapstructure:"DB_CONN_MAX_LIFETIME_MIN"`
	DBConnMaxIdleTimeMin   int    `mapstructure:"DB_CONN_MAX_IDLE_TIME_MIN"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AwsS3URL               string `mapstructure:"AWS_S3_URL"`
	AwsBucketName          string `mapstructure:"AWS_BUCKET_NAME"`
	AwsAccessKeyID         string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey     string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsRegion              string `mapstructure:"AWS_REGION"`
}

func LoadConfig() (*Env, error) {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
		return nil, err
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
		return nil, err
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env, nil
}

func NewPostgresDatabase(env *Env) *gorm.DB {
	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass
	dbName := env.DBName

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// O ping com timeout é feito abaixo, depois de configurar o pool.
		DisableAutomaticPing: true,
		Logger: logger.New(log.New(os.Stdout, "", log.LstdFlags), logger.Config{
			SlowThreshold:             slowQueryThreshold,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database instance:", err)
	}

	sqlDB.SetMaxOpenConns(orDefault(env.DBMaxOpenConns, defaultMaxOpenConns))
	sqlDB.SetMaxIdleConns(orDefault(env.DBMaxIdleConns, defaultMaxIdleConns))
	sqlDB.SetConnMaxLifetime(time.Duration(orDefault(env.DBConnMaxLifetimeMin, defaultConnMaxLifetimeMin)) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(orDefault(env.DBConnMaxIdleTimeMin, defaultConnMaxIdleTimeMin)) * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("database did not respond within %s at %s:%s: %v", pingTimeout, dbHost, dbPort, err)
	}

	return db
}

func orDefault(value, fallback int) int {
	if value <= 0 {
		return fallback
	}
	return value
}

func ClosePostgresDBConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database instance:", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Fatal("failed to close database connection:", err)
	}

	log.Println("Connection to PostgreSQL closed.")
}
