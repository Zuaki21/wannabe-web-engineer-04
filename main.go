package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	//"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	e := echo.New()

	e.POST("/add", addCityHandler)

	e.POST("/delete/:name", deleteCityHandler)

	e.GET("/cities/:cityName", getCityInfoHandler)

	e.Start(":4000")
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")
	fmt.Println(cityName)

	city := City{}
	db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName)
	if city.Name == "" {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, city)
}

func addCityHandler(c echo.Context) error {
	var cityData City
	if err := c.Bind(&cityData); err != nil {
		return c.JSON(http.StatusBadRequest, cityData)
	}

	cityState := `INSERT INTO city (ID, Name, CountryCode, District, Population) VALUES (?, ?, ?, ? ,?)`
	db.MustExec(cityState, cityData.ID, cityData.Name, cityData.CountryCode, cityData.District, cityData.Population)

	return c.String(http.StatusOK, "New city("+cityData.Name+") is added.")
}

func deleteCityHandler(c echo.Context) error {
	cityName := c.Param("name")

	cityStateDelete := `DELETE FROM city WHERE Name = ?`
	db.MustExec(cityStateDelete, cityName)

	return c.String(http.StatusOK, cityName+" is added.")
}
