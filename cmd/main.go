package main

import (
	"context"
	todo "github.com/cookiesvanilli/go_app"
	"github.com/cookiesvanilli/go_app/pkg/handler"
	"github.com/cookiesvanilli/go_app/pkg/repository"
	"github.com/cookiesvanilli/go_app/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //Импорт библиотеки lib/pq дает нам возможность работать с драйвером Postgres, без него мы просто не сможем подключиться к БД. Именно его мы указали в методе sqlx.Open("postgres", ...).
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env variables: %s", err.Error())
	}

	//инициализация базы данных
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"), // создали .env и добавили в .gitignore, чтобы не светить пароль
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize db: %s", err.Error())
	}

	// экземпляр сервера
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server) // .New() Это делается, когда в структуре есть различные поля и мы хотим передать их значения в конструкторе.
	//Когда у нас пустая структура, создавать конструктор не обязательно, можно воспользоваться стандартной конструкцией new()

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

//docker pull postgres
//docker run --name=todo-db -e POSTGRES_PASSWORD=12345 -p 5436:5432 -d postgres
//docker ps
//docker exec -it
//docker ps
//docker exec -it 14354t5y6rye  /bin/bash
// psql -U postgres
// select * from users;
