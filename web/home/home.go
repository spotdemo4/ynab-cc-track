package home

import (
	"go-htmx-test/db"
	"go-htmx-test/models"
	"go-htmx-test/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func (h HomeHandler) Any(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return h.Get(c)
	case http.MethodPut:
		return h.Put(c)
	case http.MethodPatch:
		return h.Patch(c)
	case http.MethodDelete:
		return h.Delete(c)
	default:
		return c.NoContent(http.StatusMethodNotAllowed)
	}
}

func (h HomeHandler) Get(c echo.Context) error {
	var items []models.Item

	if r := db.DB.Find(&items); r.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, r.Error.Error())
	}

	return utils.Render(c, http.StatusOK, home(items))
}

func (h HomeHandler) Put(c echo.Context) error {
	item := models.Item{
		Name: c.FormValue("name"),
	}

	if item.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	if r := db.DB.Create(&item); r.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, r.Error.Error())
	}

	println("Added item: ", item.Name)
	return utils.CombineRender(c, http.StatusOK, itemsDiv([]models.Item{item}), formElements("showPutForm"))
}

func (h HomeHandler) Patch(c echo.Context) error {
	id := c.FormValue("id")
	name := c.FormValue("name")

	if id == "" || name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id and name are required")
	}

	i, err := strconv.Atoi(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id must be an integer")
	}

	item := models.Item{
		ID:   i,
		Name: c.FormValue("name"),
	}

	if r := db.DB.Save(&item); r.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, r.Error.Error())
	}

	println("Updated item: ", item.Name)
	return utils.CombineRender(c, http.StatusOK, itemDiv(item), formElements("showPatchForm"))
}

func (h HomeHandler) Delete(c echo.Context) error {
	id := c.FormValue("id")

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	if r := db.DB.Delete(&models.Item{}, id); r.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, r.Error.Error())
	}

	println("Deleted item: ", id)
	return c.NoContent(http.StatusOK)
}
