package configs

import (
	"fmt"
	"os"
	"strconv"
	"subscriber-topic-stars/src/entities"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

type DBConnection struct {
	Gorm *gorm.DB
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
}

func ConnectDatabase() *DBConnection {
	fmt.Println("Connecting to the database...")
	gorm := connecting()

	GormDB = gorm

	connection := &DBConnection{
		Gorm: GormDB,
	}

	return connection
}

func connecting() *gorm.DB {
	config := loadConfig()
	connectionString := config.ConnectionString()

	database, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Printf("failed to connect to the database: %v\n", err)
	}

	sqlDB, err := database.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	fmt.Println("✅ GORM connected to the database successfully!")

	return database
}

func loadConfig() DatabaseConfig {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Println("invalid port number: %w", err)
	}

	tz := os.Getenv("DB_TIMEZONE")
	if tz == "" {
		tz = "UTC"
	}

	return DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
		Timezone: tz,
	}
}

func (config DatabaseConfig) ConnectionString() string {
	if config.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=%s TimeZone=%s",
			config.Host, config.Port, config.User, config.DBName, config.SSLMode, config.Timezone,
		)
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode, config.Timezone,
	)
}

func RunMigrations(db *gorm.DB) {
	fmt.Println("Running migrations...")
	gormModels := make([]interface{}, 0, len(entities.RegisteredEntities))
	for _, entity := range entities.RegisteredEntities {
		gormModels = append(gormModels, entity)
	}
	migrator := db.Migrator()

	allExist := true
	for _, m := range gormModels {
		if !migrator.HasTable(m) {
			allExist = false
			break
		}
	}
	if allExist {
		fmt.Println("ℹ️  All tables already exist (GORM), skipping migration")
		return
	}

	fmt.Println("🔧 Running GORM AutoMigrate…")
	if err := db.AutoMigrate(gormModels...); err != nil {
		fmt.Printf("❌ GORM AutoMigrate failed: %v\n", err)
	} else {
		fmt.Println("✅ GORM migration complete")
	}
}
