package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	_userHttpDelivery "simple_crud/user/delivery/http"
	_userHttpDeliveryMiddleware "simple_crud/user/delivery/http/middleware"
	_userRepo "simple_crud/user/repository/mysql"
	_userUsecase "simple_crud/user/usecase"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}

	}()

	e := echo.New()
	middL := _userHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	ur := _userRepo.NewMysqlUserRepository(dbConn)

	timeoutContext := time.Duration(viper.GetDuration("context.timeout")) * time.Second
	uu := _userUsecase.NewUserUsecase(ur, timeoutContext)
	_userHttpDelivery.NewArticleHandler(e, uu)

	log.Fatal(e.Start(viper.GetString("server.address")))

}
