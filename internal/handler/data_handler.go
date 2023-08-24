package handler

import (
	"monapi/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// get data request
func (h *Handler) getData(c *gin.Context) {

	logrus.Printf("getData(): BEGIN")

	var page = c.DefaultQuery("page", "1")
	var limit = c.DefaultQuery("limit", "10000")

	// check for query params
	// ID
	id_str, ok := c.GetQuery("id")
	if !ok || id_str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": "'id' param not found or empty"})
		logrus.Error("getData(): 'id' param not found or empty")
		return
	}
	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": "Cannot parse 'id'"})
		logrus.Errorf("getData(): Cannot parse 'id', error = %s", err.Error())
		return
	}

	// from
	from_str, ok := c.GetQuery("from")
	if !ok || from_str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": "'from' param not found or empty"})
		logrus.Error("getData(): 'from' param not found or empty")
		return
	}

	from, err := time.Parse("2006-01-02", from_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": "Cannot parse 'from'"})
		logrus.Errorf("getData(): Cannot parse 'from', error = %s", err.Error())
		return
	}

	// to
	to_str, ok := c.GetQuery("to")
	if !ok || to_str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": "'to' param not found or empty"})
		logrus.Error("getData(): 'to' param not found or empty")
		return
	}
	to, err := time.Parse("2006-01-02", to_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error: ": "Cannot parse 'from'"})
		logrus.Errorf("getData(): Cannot parse 'from', error = %s", err.Error())
		return
	}

	// amount
	default_amount := 100000
	amount := 0
	amount_str, ok := c.GetQuery("amount")
	if !ok || amount_str == "" {
		amount = default_amount
		logrus.Printf("getData(): 'amount' param not found or empty. Use the default amount = %d", amount)
	}
	amount, err = strconv.Atoi(amount_str)
	if err != nil {
		amount = default_amount
		logrus.Errorf("getData(): Cannot parse 'amount', error = %s, use the default amount=%d", err.Error(), amount)
	}

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)

	if intLimit > 10000 {
		intLimit = 10000
	}

	offset := (intPage - 1) * intLimit

	logrus.Printf("getData(): Try to get data for id = %s, from = %s, to = %s", id_str, from_str, to_str)
	items, err := h.services.Data.GetData(id, from, to, intLimit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		logrus.Errorf("getData(): Cannot get data for id = %s, from = %s, to = %s, error = %s", id_str, from_str, to_str, err.Error())
		return
	}

	logrus.Printf("getData(): Try to reduce amnount from %d to = %d", len(items), amount)

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

	c.JSON(http.StatusOK, gin.H{"total": total, "amount": len(result), "step": step, "data": result})

	logrus.Printf("getData(): END")
}
