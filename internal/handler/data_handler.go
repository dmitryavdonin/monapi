package handler

import (
	"monapi/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) algo1(items []model.Data_new, amount int) []model.DTO {

	logrus.Printf("algo1(): BEGIN")

	total := len(items)

	var result []model.DTO

	step := total / amount
	result = make([]model.DTO, amount)
	source_index := 0

	for i := 0; i < amount; i++ {
		result[i] = model.DTO{
			Temperature: items[source_index].Temperature,
			Humidity:    items[source_index].Temperature,
			Time:        items[source_index].DtWr.Format("2006-01-02 15:04:05")}

		source_index += step

		if (source_index + step) >= total {
			break
		}
	}

	logrus.Printf("algo1(): END")

	return result
}

func (h *Handler) findNearest(items []model.Data_new, point time.Time, D_range int) *model.Data_new {
	var result *model.Data_new = nil

	left_point := point.Add(-time.Second * time.Duration(D_range))
	right_pont := point.Add(time.Second * time.Duration(D_range))

	for _, item := range items {

		if item.DtWr.Unix() >= left_point.Unix() && item.DtWr.Unix() < right_pont.Unix() {
			result = &item
			break
		}
	}

	return result
}

// algo2 implementation
func (h *Handler) algo2(items []model.Data_new, amount int, from time.Time, to time.Time, D int) []model.DTO {

	logrus.Printf("algo2(): BEGIN")

	var result []model.DTO

	span := to.Sub(from)

	step_in_sec := int(span.Seconds()) / amount

	D_range := step_in_sec * D / 100

	result = make([]model.DTO, amount)

	for i := 0; i < amount; i++ {
		point := from.Add(time.Second * time.Duration(step_in_sec*i))
		item := h.findNearest(items, point, D_range)
		if item != nil {
			result[i] = model.DTO{
				Temperature: item.Temperature,
				Humidity:    item.Humidity,
				Time:        item.DtWr.Format("2006-01-02 15:04:05"),
			}
		} else {
			result[i] = model.DTO{
				Temperature: 1000,
				Humidity:    1000,
				Time:        point.Format("2006-01-02 15:04:05"),
			}
		}
	}

	logrus.Printf("algo2(): END")

	return result
}

// get data request
func (h *Handler) getData(c *gin.Context) {

	logrus.Printf("getData(): BEGIN")

	var page = c.DefaultQuery("page", "1")
	var limit = c.DefaultQuery("limit", "10000")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)

	if intLimit > 10000 {
		intLimit = 10000
	}

	offset := (intPage - 1) * intLimit

	parser, err := InitParamParser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
		logrus.Errorf("getData(): Cannot parse params, error = %s", err.Error())
		return
	}

	logrus.Printf("getData(): Try to get data for id = %d, from = %s, to = %s",
		parser.ID, parser.From.Format("2006-01-02 15:04:05"), parser.To.Format("2006-01-02 15:04:05"))

	items, err := h.services.Data.GetData(parser.ID, parser.From, parser.To, intLimit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		logrus.Errorf("getData(): Cannot get data for id = %d, from = %s, to = %s, error = %s",
			parser.ID, parser.From.Format("2006-01-02 15:04:05"), parser.To.Format("2006-01-02 15:04:05"), err.Error())
		return
	}

	total := len(items)

	if total < parser.Amount {
		logrus.Printf("getData(): total amount = %d is less than requested amount = %d", total, parser.Amount)
		//parser.Amount = total
	}

	var result []model.DTO
	if parser.Algo == 2 {
		result = h.algo2(items, parser.Amount, parser.From, parser.To, parser.D)
	} else {
		result = h.algo1(items, parser.Amount)
	}

	c.JSON(http.StatusOK, gin.H{"amount": len(result), "data": result})

	logrus.Printf("getData(): END")
}

func (h *Handler) getLastValue(c *gin.Context) {
	id_str := c.Param("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		logrus.Errorf("getData(): Cannot parse 'id', error = %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": err.Error()})
		return
	}

	item, err := h.services.GetLastValue(id)
	if err != nil {
		logrus.Errorf("getData(): Cannot find data for 'id' int the last 1 day and last 7 days, error = %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"Error: ": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.DTO{Temperature: item.Temperature, Humidity: item.Humidity, Time: item.DtWr.Format("2006-01-02 15:04:05")})
}
