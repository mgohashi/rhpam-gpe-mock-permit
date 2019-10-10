package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

//Status s
type Status bool

//Permit p
type Permit struct {
	ID     int    `json:"id" form:"id" query:"id"`
	Pid    int    `json:"pid" form:"pid" query:"pid"`
	Status Status `json:"status" form:"status" query:"status"`
}

var electricalDb []*Permit
var structuralDb []*Permit
var elecStatus Status
var structStatus Status

func main() {
	e := echo.New()
	elecStatus = true
	structStatus = true

	e.GET("/", getTypes)
	e.GET("/:type/:id", getPermit)
	e.POST("/:type", createPermit)
	e.POST("/:type/status/:status", cancelPermits)

	e.Logger.Fatal(e.Start(":8082"))
}

func getTypes(c echo.Context) error {
	types := []string{"electrical", "structural"}
	return c.JSON(http.StatusOK, types)
}

func getPermit(c echo.Context) error {
	t := c.Param("type")
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var permit *Permit

	switch t {
	case "electrical":
		permit, err = getPermitFromDB(electricalDb, id)
	case "structural":
		permit, err = getPermitFromDB(structuralDb, id)
	default:
		return echo.NewHTTPError(http.StatusNotFound, "Type not found!")
	}

	if permit != nil {
		fmt.Printf("%v permit %v found\n", t, toJSON(permit))
		return c.JSON(http.StatusOK, permit)
	}

	return echo.NewHTTPError(http.StatusNotFound, err.Error())

}

func createPermit(c echo.Context) error {
	t := c.Param("type")
	p := new(Permit)

	if err := c.Bind(p); err != nil {
		return err
	}

	switch t {
	case "electrical":
		id := len(electricalDb) + 1
		p.ID = id
		p.Status = elecStatus
		electricalDb = append(electricalDb, p)
	case "structural":
		id := len(structuralDb) + 1
		p.ID = id
		p.Status = structStatus
		structuralDb = append(structuralDb, p)
	default:
		return echo.NewHTTPError(http.StatusNotFound, "Type not found!")
	}

	fmt.Printf("%v permit %v created\n", t, toJSON(p))

	return c.JSON(http.StatusOK, p)
}

func cancelPermits(c echo.Context) error {
	t := c.Param("type")
	status, err := strconv.ParseBool(c.Param("status"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch t {
	case "electrical":
		elecStatus = Status(status)
	case "structural":
		structStatus = Status(status)
	default:
		return echo.NewHTTPError(http.StatusNotFound, "Type not found!")
	}

	fmt.Printf("All %v permit's statuses will be %v\n", t, status)

	return c.NoContent(http.StatusOK)
}

func getPermitFromDB(db []*Permit, id int) (*Permit, error) {
	if len(db) > id-1 {
		for i, v := range db {
			if i == id-1 {
				return v, nil
			}
		}
	}

	return nil, fmt.Errorf("%d not fount", id)
}

func toJSON(a interface{}) string {
	b, err := json.Marshal(a)

	if err != nil {
		fmt.Println(err.Error())
		return "[Marshal Error]"
	}

	return string(b)
}
