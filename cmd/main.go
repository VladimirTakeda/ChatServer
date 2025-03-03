package main

import (
	"ChatServer/cmd/internal"
	websocket2 "ChatServer/pkg/connection/websocket"
	"ChatServer/pkg/handlers/http"
	"ChatServer/pkg/handlers/ws"
	minio2 "ChatServer/pkg/repository/minio"
	"ChatServer/pkg/repository/postgres"
	redis2 "ChatServer/pkg/repository/redis"
	"ChatServer/pkg/service"
	"context"
	"database/sql"
	"fmt"
	"github.com/antelman107/net-wait-go/wait"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config %s", err.Error())
		os.Exit(1)
	}

	if err := gotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables %s", err.Error())
		os.Exit(1)
	}

	if !wait.New(
		wait.WithProto("tcp"),
		wait.WithWait(200*time.Millisecond),
		wait.WithBreak(1000*time.Millisecond),
		wait.WithDeadline(15*time.Second),
		wait.WithDebug(true),
	).Do([]string{fmt.Sprintf("%s:%s", viper.GetString("postgres.host"), viper.GetString("postgres.port"))}) {
		logrus.Fatalf("db is not available")
		return
	}

	pool, err := postgres.NewPostgresDb(postgres.Config{
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetString("postgres.port"),
		UserName: viper.GetString("postgres.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   viper.GetString("postgres.dbname"),
		SSLMode:  viper.GetString("postgres.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to connect to db, %s", err.Error())
		os.Exit(1)
	}

	pgURL := fmt.Sprintf("postgres://postgres:%s@%s:%s/postgres?sslmode=disable", os.Getenv("DB_PASSWORD"), viper.GetString("postgres.host"), viper.GetString("postgres.port"))

	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		logrus.Fatalf("failed to prepare db, %s", err.Error())
		os.Exit(1)
	}

	err = postgres.RunMigrations("file://migrations", db)

	redisCache, err := redis2.NewRedisDb(redis2.Config{
		Host:     viper.GetString("cache.host"),
		Port:     viper.GetString("cache.port"),
		UserName: viper.GetString("cache.username"),
		Password: viper.GetString("cache.password"),
		DbNumber: viper.GetInt("cache.dbnumber"),
	})

	if err != nil {
		logrus.Fatalf("failed to open redis cache, %s", err.Error())
		os.Exit(1)
	}

	redisPubSub, err := redis2.NewRedisDb(redis2.Config{
		Host:     viper.GetString("pubsub.host"),
		Port:     viper.GetString("pubsub.port"),
		UserName: viper.GetString("pubsub.username"),
		Password: viper.GetString("pubsub.password"),
		DbNumber: viper.GetInt("pubsub.dbnumber"),
	})

	if err != nil {
		logrus.Fatalf("failed to open redis pubsub, %s", err.Error())
		os.Exit(1)
	}

	minio, err := minio2.NewMinioClient(minio2.Config{
		KeyID:     viper.GetString("minio.keyID"),
		AccessKey: viper.GetString("minio.accessKey"),
		EndPoint:  viper.GetString("minio.endpoint"),
	})

	if err != nil {
		logrus.Fatalf("failed to create minio client, %s", err.Error())
		os.Exit(1)
	}

	s3Storage := minio2.NewS3Storage(minio)

	repos := postgres.NewRepository(pool)
	cache := redis2.NewCache(redisCache)
	pubsub := redis2.NewPubSub(redisPubSub)

	services := service.NewService(repos, cache, pubsub, s3Storage)

	httpHandlers := http.NewHandler(services)
	serverHub := websocket2.NewServerHub(services)
	wsHandlers := ws.NewHandler(serverHub, services)
	go serverHub.Run()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	httpHandlers.SetupRoutes(router)
	wsHandlers.SetupRoutes(router)

	srv := new(internal.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), router); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
			os.Exit(1)
		}
	}()

	logrus.Print("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	pool.Close()
}

func initConfig() error {
	viper.AddConfigPath("cmd/internal")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
