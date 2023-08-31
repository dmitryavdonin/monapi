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
	Algo   int
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
	default_algo := 2
	default_D := 50

	amount := 0
	amount_str, ok := c.GetQuery("amount")
	if !ok || amount_str == "" {
		return nil, errors.New("'amount' param not found or empty")
	}
	amount, err = strconv.Atoi(amount_str)
	if err != nil {
		return nil, fmt.Errorf("cannot parse 'amount', error = %s", err.Error())
	}

	if amount == 0 {
		return nil, errors.New("'amount' cannot be 0")
	}

	if amount > default_amount {
		logrus.Printf("getData(): 'amount' = %d cannot be more than default_amount = %d", amount, default_amount)
		amount = default_amount
	}

	algo := 0
	algo_str, ok := c.GetQuery("algo")
	if !ok || algo_str == "" {
		algo = default_algo
		logrus.Printf("getData(): 'algo' param not found or empty. Use the default algo = %d", algo)
	}
	algo, err = strconv.Atoi(algo_str)
	if err != nil {
		algo = default_algo
		logrus.Errorf("getData(): Cannot parse 'algo', error = %s, use the default algo=%d", err.Error(), algo)
	}
	if algo > 2 {
		algo = default_algo
	}

	D := 0
	D_str, ok := c.GetQuery("D")
	if !ok || D_str == "" {
		D = default_D
		logrus.Printf("getData(): 'algo' param not found or empty. Use the default algo = %d", algo)
	}
	D, err = strconv.Atoi(D_str)
	if err != nil {
		D = default_D
		logrus.Errorf("getData(): Cannot parse 'D', error = %s, use the default D = %d", err.Error(), D)
	}
	if D > 100 || D <= 0 {
		D = default_D
	}

	return &ParamParser{
		ID:     id,
		From:   from,
		To:     to,
		Amount: amount,
		Algo:   algo,
		D:      D,
	}, nil
}
