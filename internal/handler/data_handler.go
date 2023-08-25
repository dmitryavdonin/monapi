package handler

import (
	"monapi/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) algo1(items []model.Data_new, amount int) []model.Data_new {

	logrus.Printf("algo1(): BEGIN")

	total := len(items)
	step := 1

	var result []model.Data_new

	if total <= amount {
		result = items
	} else {
		step = total / amount
		result = make([]model.Data_new, amount)
		i := 0
		source_index := 0
		for {
			if i >= len(result) {
				break
			}
			result[i] = items[source_index]
			i++
			source_index += step
			if source_index+step >= total {
				break
			}
		}
	}

	logrus.Printf("algo1(): END")

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
		parser.ID, parser.From.Format("2006-01-02"), parser.To.Format("2006-01-02"))

	items, err := h.services.Data.GetData(parser.ID, parser.From, parser.To, intLimit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		logrus.Errorf("getData(): Cannot get data for id = %d, from = %s, to = %s, error = %s",
			parser.ID, parser.From.Format("2006-01-02"), parser.To.Format("2006-01-02"), err.Error())
		return
	}

	total := len(items)

	if total < parser.Amount {
		logrus.Printf("getData(): Nothing to reduce becaue total amount = %d is less than requested amount = %d. So just return the total amount", total, parser.Amount)
		c.JSON(http.StatusOK, gin.H{"total": total, "amount": total, "step": 1, "data": items})
		return
	}

	step := total / parser.Amount

	logrus.Printf("getData(): Try to reduce amnount from %d to = %d, step = %d, using algo = %d", len(items), parser.Amount, step, parser.Algo)

	var result []model.Data_new
	if parser.Algo == 1 {
		result = h.algo1(items, parser.Amount)
	}

	c.JSON(http.StatusOK, gin.H{"total": total, "amount": len(result), "step": step, "data": result})

	logrus.Printf("getData(): END")
}
