package handler

import (
	"monapi/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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

	logrus.Printf("getData(): Try to reduce amnount from %d to = %d", len(items), parser.Amount)

	total := len(items)
	step := 1

	var result []model.Data_new

	if total <= parser.Amount {
		result = items
	} else {
		step = total / parser.Amount
		result = make([]model.Data_new, parser.Amount)
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

	c.JSON(http.StatusOK, gin.H{"total": total, "amount": len(result), "step": step, "data": result})

	logrus.Printf("getData(): END")
}
