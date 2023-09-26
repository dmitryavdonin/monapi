package handler

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ParamParser struct {
	ID     int
	Amount int
	From   time.Time
	To     time.Time
	D      int // допуск в процентах
}

func InitParamParser(c *gin.Context) (*ParamParser, error) {
	// ID
	id_str, ok := c.GetQuery("id")
	if !ok || id_str == "" {
		return nil, errors.New("'id' param not found or empty")
	}
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return nil, fmt.Errorf("cannot parse 'id', error = %s", err.Error())
	}

	// from
	from_str, ok := c.GetQuery("from")
	if !ok || from_str == "" {
		return nil, errors.New("'from' param not found or empty")
	}

	from, err := time.Parse("2006-01-02 15:04:05", from_str)
	if err != nil {
		return nil, fmt.Errorf("cannot parse 'from', error = %s", err.Error())
	}

	// to
	to_str, ok := c.GetQuery("to")
	if !ok || to_str == "" {
		return nil, errors.New("'to' param not found or empty")
	}
	to, err := time.Parse("2006-01-02 15:04:05", to_str)
	if err != nil {
		return nil, fmt.Errorf("cannot parse 'to', error = %s", err.Error())
	}

	// amount
	default_amount := 10000
	default_D := 50

	amount := 0
	amount_str, ok := c.GetQuery("amount")
	if !ok || amount_str == "" {
		amount = 0
	}
	amount, err = strconv.Atoi(amount_str)
	if err != nil {
		amount = 0
	}

	if amount > default_amount {
		logrus.Printf("getData(): 'amount' = %d cannot be more than default_amount = %d", amount, default_amount)
		amount = default_amount
	}

	D := 0
	D_str, ok := c.GetQuery("D")
	if !ok || D_str == "" {
		D = default_D
		logrus.Printf("getData(): 'D' param not found or empty. Use the default D = %d", D)
	}
	D, err = strconv.Atoi(D_str)
	if err != nil {
		D = default_D
		logrus.Errorf("getData(): Cannot parse 'D', error = %s, use the default D = %d", err.Error(), D)
	}
	if D > 50 || D <= 0 {
		D = default_D
	}

	return &ParamParser{
		ID:     id,
		From:   from,
		To:     to,
		Amount: amount,
		D:      D,
	}, nil
}
