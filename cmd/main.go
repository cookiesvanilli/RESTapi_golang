package main

import (
	todo "github.com/cookiesvanilli/go_app"
	"github.com/cookiesvanilli/go_app/pkg/handler"
	"github.com/cookiesvanilli/go_app/pkg/repository"
	"github.com/cookiesvanilli/go_app/pkg/service"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	// экземпляр сервера
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server) // .New() Это делается, когда в структуре есть различные поля и мы хотим передать их значения в конструкторе.
	//Когда у нас пустая структура, создавать конструктор не обязательно, можно воспользоваться стандартной конструкцией new()
	if err := srv.Run(viper.GetString("8000"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Error while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

//docker pull postgres
//docker run --name=todo-db -e POSTGRES_PASSWORD='12345' -p 5436:5432 -d --rm postgres
//docker ps
