package handler

import (
	"monapi/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"math"
)

var Version string = "2.0.4"

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
	right_point := point.Add(time.Second * time.Duration(D_range))

	m_left := make(map[int64]model.Data_new)
	m_right := make(map[int64]model.Data_new)

	for _, item := range items {

		// left points
		if item.DtWr.Unix() >= left_point.Unix() && item.DtWr.Unix() < point.Unix() {
			span := point.Sub(item.DtWr)
			m_left[span.Milliseconds()] = item
		}

		// right points
		if item.DtWr.Unix() >= point.Unix() && item.DtWr.Unix() < right_point.Unix() {
			span := item.DtWr.Sub(point)
			m_right[span.Milliseconds()] = item
		}

		if item.DtWr.Unix() >= right_point.Unix() {
			break
		}
	}

	var min_left int64 = math.MinInt64
	var min_right int64 = math.MinInt64

	// find the min in left
	if len(m_left) > 0 {
		for k := range m_left {
			if min_left < k {
				min_left = k
			}
		}
	}

	// find the min in right
	if len(m_right) > 0 {
		for k := range m_right {
			if min_right < k {
				min_right = k
			}
		}
	}

	if min_right != math.MinInt64 && min_left != math.MinInt64 {
		if min_left < min_right {
			result = &model.Data_new{
				ID:          m_left[min_left].ID,
				Temperature: m_left[min_left].Temperature,
				Humidity:    m_left[min_left].Humidity,
				DtWr:        m_left[min_left].DtWr,
			}
		} else {
			result = &model.Data_new{
				ID:          m_right[min_right].ID,
				Temperature: m_right[min_right].Temperature,
				Humidity:    m_right[min_right].Humidity,
				DtWr:        m_right[min_right].DtWr,
			}
		}
	} else if min_right != math.MinInt64 {
		result = &model.Data_new{
			ID:          m_right[min_right].ID,
			Temperature: m_right[min_right].Temperature,
			Humidity:    m_right[min_right].Humidity,
			DtWr:        m_right[min_right].DtWr,
		}
	} else if min_left != math.MinInt64 {
		result = &model.Data_new{
			ID:          m_left[min_left].ID,
			Temperature: m_left[min_left].Temperature,
			Humidity:    m_left[min_left].Humidity,
			DtWr:        m_left[min_left].DtWr,
		}
	}

	return result
}

func (h *Handler) findNearestRightOnly(items []model.Data_new, point time.Time, D_range int) *model.Data_new {
	var result *model.Data_new = nil
	right_point := point.Add(time.Second * time.Duration(D_range))

	m_right := make(map[int64]model.Data_new)

	for _, item := range items {
		// right points
		if item.DtWr.Unix() >= point.Unix() && item.DtWr.Unix() < right_point.Unix() {
			span := item.DtWr.Sub(point)
			m_right[span.Milliseconds()] = item
		}

		if item.DtWr.Unix() >= right_point.Unix() {
			break
		}
	}

	var min_right int64 = math.MinInt64

	// find the min in right
	if len(m_right) > 0 {
		for k := range m_right {
			if min_right < k {
				min_right = k
			}
		}
	}

	if min_right != math.MinInt64 {
		result = &model.Data_new{
			ID:          m_right[min_right].ID,
			Temperature: m_right[min_right].Temperature,
			Humidity:    m_right[min_right].Humidity,
			DtWr:        m_right[min_right].DtWr,
		}
	}

	return result
}

func (h *Handler) findNearestLeftOnly(items []model.Data_new, point time.Time, D_range int) *model.Data_new {
	var result *model.Data_new = nil

	left_point := point.Add(-time.Second * time.Duration(D_range))

	m_left := make(map[int64]model.Data_new)

	for _, item := range items {

		// left points
		if item.DtWr.Unix() >= left_point.Unix() && item.DtWr.Unix() < point.Unix() {
			span := point.Sub(item.DtWr)
			m_left[span.Milliseconds()] = item
		}

		if item.DtWr.Unix() >= point.Unix() {
			break
		}
	}

	var min_left int64 = math.MinInt64

	// find the min in left
	if len(m_left) > 0 {
		for k := range m_left {
			if min_left < k {
				min_left = k
			}
		}
	}

	if min_left != math.MinInt64 {
		result = &model.Data_new{
			ID:          m_left[min_left].ID,
			Temperature: m_left[min_left].Temperature,
			Humidity:    m_left[min_left].Humidity,
			DtWr:        m_left[min_left].DtWr,
		}
	}

	return result
}

func (h *Handler) findNearestDebug(items []model.Data_new, point time.Time, D_range int) (r *model.Data_new, left time.Time, right time.Time) {
	var result *model.Data_new = nil

	left_point := point.Add(-time.Second * time.Duration(D_range))
	right_point := point.Add(time.Second * time.Duration(D_range))

	m_left := make(map[int64]model.Data_new)
	m_right := make(map[int64]model.Data_new)

	for _, item := range items {

		// left ponts
		if item.DtWr.Unix() >= left_point.Unix() && item.DtWr.Unix() < point.Unix() {
			span := point.Sub(item.DtWr)
			m_left[span.Milliseconds()] = item
		}

		// right points
		if item.DtWr.Unix() >= point.Unix() && item.DtWr.Unix() < right_point.Unix() {
			span := item.DtWr.Sub(point)
			m_right[span.Milliseconds()] = item
		}

		if item.DtWr.Unix() >= right_point.Unix() {
			break
		}
	}

	var min_left int64 = math.MinInt64
	var min_right int64 = math.MinInt64

	// find the min in left
	if len(m_left) > 0 {
		for k := range m_left {
			if min_left < k {
				min_left = k
			}
		}
	}

	// find the min in right
	if len(m_right) > 0 {
		for k := range m_right {
			if min_right < k {
				min_right = k
			}
		}
	}

	if min_right != math.MinInt64 && min_left != math.MinInt64 {
		if min_left < min_right {
			result = &model.Data_new{
				ID:          m_left[min_left].ID,
				Temperature: m_left[min_left].Temperature,
				Humidity:    m_left[min_left].Humidity,
				DtWr:        m_left[min_left].DtWr,
			}
		} else {
			result = &model.Data_new{
				ID:          m_right[min_right].ID,
				Temperature: m_right[min_right].Temperature,
				Humidity:    m_right[min_right].Humidity,
				DtWr:        m_right[min_right].DtWr,
			}
		}
	} else if min_right != math.MinInt64 {
		result = &model.Data_new{
			ID:          m_right[min_right].ID,
			Temperature: m_right[min_right].Temperature,
			Humidity:    m_right[min_right].Humidity,
			DtWr:        m_right[min_right].DtWr,
		}
	} else if min_left != math.MinInt64 {
		result = &model.Data_new{
			ID:          m_left[min_left].ID,
			Temperature: m_left[min_left].Temperature,
			Humidity:    m_left[min_left].Humidity,
			DtWr:        m_left[min_left].DtWr,
		}
	}

	return result, left_point, right_point
}

func (h *Handler) findNearestRightOnlytDebug(items []model.Data_new, point time.Time, D_range int) (r *model.Data_new, left time.Time, right time.Time) {
	var result *model.Data_new = nil

	right_point := point.Add(time.Second * time.Duration(D_range))
	m_right := make(map[int64]model.Data_new)

	for _, item := range items {
		// right points
		if item.DtWr.Unix() >= point.Unix() && item.DtWr.Unix() < right_point.Unix() {
			span := item.DtWr.Sub(point)
			m_right[span.Milliseconds()] = item
		}

		if item.DtWr.Unix() >= right_point.Unix() {
			break
		}
	}

	var min_right int64 = math.MinInt64

	// find the min in right
	if len(m_right) > 0 {
		for k := range m_right {
			if min_right < k {
				min_right = k
			}
		}
	}

	if min_right != math.MinInt64 {
		result = &model.Data_new{
			ID:          m_right[min_right].ID,
			Temperature: m_right[min_right].Temperature,
			Humidity:    m_right[min_right].Humidity,
			DtWr:        m_right[min_right].DtWr,
		}
	}

	return result, point, right_point
}

func (h *Handler) findNearestLeftOnlyDebug(items []model.Data_new, point time.Time, D_range int) (r *model.Data_new, left time.Time, right time.Time) {
	var result *model.Data_new = nil

	left_point := point.Add(-time.Second * time.Duration(D_range))

	m_left := make(map[int64]model.Data_new)

	for _, item := range items {

		// left ponts
		if item.DtWr.Unix() >= left_point.Unix() && item.DtWr.Unix() < point.Unix() {
			span := point.Sub(item.DtWr)
			m_left[span.Milliseconds()] = item
		}

		if item.DtWr.Unix() > point.Unix() {
			break
		}
	}

	var min_left int64 = math.MinInt64

	// find the min in left
	if len(m_left) > 0 {
		for k := range m_left {
			if min_left < k {
				min_left = k
			}
		}
	}

	if min_left != math.MinInt64 {
		result = &model.Data_new{
			ID:          m_left[min_left].ID,
			Temperature: m_left[min_left].Temperature,
			Humidity:    m_left[min_left].Humidity,
			DtWr:        m_left[min_left].DtWr,
		}
	}

	return result, left_point, point
}

// algo2 implementation
func (h *Handler) algo2(items []model.Data_new, amount int, from time.Time, to time.Time, D int) []model.DTO {

	logrus.Printf("algo2(): BEGIN")

	var result []model.DTO

	span := to.Sub(from)

	step_in_sec := int(span.Seconds()) / amount

	D_range := step_in_sec * D / 100

	/* Общее правило выравнивания допуска:
	Если расчетный размер допуска (D) меньше 60 секунд - берем ровно 60 секунд.*/
	if D_range < 60 {
		D_range = 60
		logrus.Printf("algo2(): D_range is less than 60 sec, so let have D_range = 60 sec")
	}

	/*Уточнение для исключения "нахлеста":
	НО, если получается так, что допуск размером 60 секунд менее 50% интервала, то нужно взять строго допуск = 50% интервала.*/
	if D_range < (step_in_sec / 2) {
		D_range = (step_in_sec / 2)
		logrus.Printf("algo2(): D_range = %d is less than step_in_sec / 2 = %d sec, so let have D_range = step_in_sec / 2",
			D_range, (step_in_sec / 2))
	}

	for i := 0; i < amount+1; i++ {
		point := from.Add(time.Second * time.Duration(step_in_sec*i))
		var item *model.Data_new
		var dto model.DTO

		if i == 0 {
			// first point
			item = h.findNearestRightOnly(items, point, D_range)
		} else if i == amount {
			item = h.findNearestLeftOnly(items, point, D_range)
		} else {
			// last point
			item = h.findNearest(items, point, D_range)
		}

		if item != nil {
			dto = model.DTO{
				Temperature: item.Temperature,
				Humidity:    item.Humidity,
				Time:        item.DtWr.Format("2006-01-02 15:04:05"),
			}
		} else {
			dto = model.DTO{
				Temperature: 1000,
				Humidity:    1000,
				Time:        point.Format("2006-01-02 15:04:05"),
			}
		}
		result = append(result, dto)
	}

	logrus.Printf("algo2(): END")

	return result
}

func (h *Handler) algo2Debug(items []model.Data_new, amount int, from time.Time, to time.Time, D int) ([]model.DTODebug, int, int) {

	logrus.Printf("algo2Debug(): BEGIN")

	var result []model.DTODebug

	span := to.Sub(from)

	step_in_sec := int(span.Seconds()) / amount

	D_range := step_in_sec * D / 100

	/* Общее правило выравнивания допуска:
	Если расчетный размер допуска (D) меньше 60 секунд - берем ровно 60 секунд.*/
	if D_range < 60 {
		D_range = 60
		logrus.Printf("algo2Debug(): D_range is less than 60 sec, so let have D_range = 60 sec")
	}

	/*Уточнение для исключения "нахлеста":
	НО, если получается так, что допуск размером 60 секунд менее 50% интервала, то нужно взять строго допуск = 50% интервала.*/
	if D_range < (step_in_sec / 2) {
		D_range = (step_in_sec / 2)
		logrus.Printf("algo2Debug(): D_range = %d is less than step_in_sec / 2 = %d sec, so let have D_range = step_in_sec / 2",
			D_range, (step_in_sec / 2))
	}

	for i := 0; i < amount+1; i++ {
		point := from.Add(time.Second * time.Duration(step_in_sec*i))
		var item *model.Data_new
		var dto model.DTODebug
		var left, right time.Time

		if i == 0 {
			item, left, right = h.findNearestRightOnlytDebug(items, point, D_range)
		} else if i == amount {
			item, left, right = h.findNearestLeftOnlyDebug(items, point, D_range)
		} else {
			item, left, right = h.findNearestDebug(items, point, D_range)
		}

		if item != nil {
			dto = model.DTODebug{
				Temperature: item.Temperature,
				Humidity:    item.Humidity,
				Time:        item.DtWr.Format("2006-01-02 15:04:05"),
				Point:       point.Format("2006-01-02 15:04:05"),
				Left:        left.Format("2006-01-02 15:04:05"),
				Right:       right.Format("2006-01-02 15:04:05"),
			}
		} else {
			dto = model.DTODebug{
				Temperature: 1000,
				Humidity:    1000,
				Time:        point.Format("2006-01-02 15:04:05"),
				Point:       point.Format("2006-01-02 15:04:05"),
				Left:        left.Format("2006-01-02 15:04:05"),
				Right:       right.Format("2006-01-02 15:04:05"),
			}
		}
		result = append(result, dto)
	}

	logrus.Printf("algo2Debug(): END")

	return result, step_in_sec, D_range
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

func (h *Handler) getDataDebug(c *gin.Context) {

	logrus.Printf("getDataDebug(): BEGIN")

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
		logrus.Errorf("getDataDebug(): Cannot parse params, error = %s", err.Error())
		return
	}

	logrus.Printf("getDataDebug(): Try to get data for id = %d, from = %s, to = %s",
		parser.ID, parser.From.Format("2006-01-02 15:04:05"), parser.To.Format("2006-01-02 15:04:05"))

	items, err := h.services.Data.GetData(parser.ID, parser.From, parser.To, intLimit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		logrus.Errorf("getDataDebug(): Cannot get data for id = %d, from = %s, to = %s, error = %s",
			parser.ID, parser.From.Format("2006-01-02 15:04:05"), parser.To.Format("2006-01-02 15:04:05"), err.Error())
		return
	}

	total := len(items)

	if total < parser.Amount {
		logrus.Printf("getDataDebug(): total amount = %d is less than requested amount = %d", total, parser.Amount)
		//parser.Amount = total
	}

	result, step_in_sec, D_in_sec := h.algo2Debug(items, parser.Amount, parser.From, parser.To, parser.D)

	c.JSON(http.StatusOK, gin.H{"amount": len(result), "step_in_sec": step_in_sec, "D_in_sec": D_in_sec, "data": result})

	logrus.Printf("getDataDebug(): END")
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

func (h *Handler) getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version:": Version})
}
