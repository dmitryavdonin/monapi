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

var Version string = "2.0.7"

// Поиск ближайшей фактической точки к расчетной точке в рамках заданной дельта-окрестности слева и справа
func (h *Handler) findNearest(items []model.Data_new, point time.Time, D_range int) *model.Data_new {
	var result *model.Data_new = nil

	left_point := point.Add(-time.Second * time.Duration(D_range))
	right_point := point.Add(time.Second * time.Duration(D_range))

	m_left := make(map[int64]model.Data_new)
	m_right := make(map[int64]model.Data_new)

	for _, item := range items {
		// фактические точки слева
		if item.DtWr.Unix() >= left_point.Unix() && item.DtWr.Unix() < point.Unix() {
			span := point.Sub(item.DtWr)
			m_left[span.Milliseconds()] = item
		}
		// фактические точки справа
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

	// найти ближайшую слева
	if len(m_left) > 0 {
		for k := range m_left {
			if min_left == math.MinInt64 {
				min_left = k
				continue
			}
			if k < min_left {
				min_left = k
			}
		}
	}

	// найти ближайшую справа
	if len(m_right) > 0 {
		for k := range m_right {
			if min_right == math.MinInt64 {
				min_right = k
				continue
			}
			if k < min_right {
				min_right = k
			}
		}
	}

	// выбрать ближайшую между левой и правой и положить в массив возвращаемого результата
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

// Поиск ближайшей фактической точки к расчетной точке в рамках заданной дельта-окрестности справа
func (h *Handler) findNearestRightOnly(items []model.Data_new, point time.Time, D_range int) *model.Data_new {
	var result *model.Data_new = nil
	right_point := point.Add(time.Second * time.Duration(D_range))

	m_right := make(map[int64]model.Data_new)

	for _, item := range items {
		// точки справа
		if item.DtWr.Unix() >= point.Unix() && item.DtWr.Unix() < right_point.Unix() {
			span := item.DtWr.Sub(point)
			m_right[span.Milliseconds()] = item
		}
		if item.DtWr.Unix() >= right_point.Unix() {
			break
		}
	}

	var min_right int64 = math.MinInt64

	// выбрать ближайшую справа
	if len(m_right) > 0 {
		for k := range m_right {
			if min_right == math.MinInt64 {
				min_right = k
				continue
			}
			if k < min_right {
				min_right = k
			}
		}
	}

	// добавить в массив возвращаемых значений
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

// Поиск ближайшей фактической точки к расчетной точке в рамках заданной дельта-окрестности слева
func (h *Handler) findNearestLeftOnly(items []model.Data_new, point time.Time, D_range int) *model.Data_new {
	var result *model.Data_new = nil

	left_point := point.Add(-time.Second * time.Duration(D_range))

	m_left := make(map[int64]model.Data_new)

	for _, item := range items {
		// точки слева
		if item.DtWr.Unix() >= left_point.Unix() && item.DtWr.Unix() < point.Unix() {
			span := point.Sub(item.DtWr)
			m_left[span.Milliseconds()] = item
		}
		if item.DtWr.Unix() >= point.Unix() {
			break
		}
	}

	var min_left int64 = math.MinInt64

	// выбрать ближайшую слева
	if len(m_left) > 0 {
		for k := range m_left {
			if min_left == math.MinInt64 {
				min_left = k
				continue
			}
			if k < min_left {
				min_left = k
			}
		}
	}

	// положить в массив возвращаемых значений
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
			if min_left == math.MinInt64 {
				min_left = k
				continue
			}
			if k < min_left {
				min_left = k
			}
		}
	}

	// find the min in right
	if len(m_right) > 0 {
		for k := range m_right {
			if min_right == math.MinInt64 {
				min_right = k
				continue
			}
			if k < min_right {
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
			if min_right == math.MinInt64 {
				min_right = k
				continue
			}
			if k < min_right {
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
			if min_left == math.MinInt64 {
				min_left = k
				continue
			}
			if k < min_left {
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

// Реализация алгоритма прореживания данных
func (h *Handler) prepareData(items []model.Data_new, amount int, from time.Time, to time.Time, D int) []model.DTO {

	logrus.Printf("prepareData(): BEGIN")

	var result []model.DTO

	span := to.Sub(from)

	step_in_sec := int(span.Seconds()) / amount

	D_range := step_in_sec * D / 100

	// Общее правило выравнивания допуска:
	// Если расчетный размер допуска (D) меньше 60 секунд - берем ровно 60 секунд.
	if D_range < 60 {
		D_range = 60
		logrus.Printf("prepareData(): D_range is less than 60 sec, so let have D_range = 60 sec")
	}

	// Уточнение для исключения "нахлеста":
	// НО, если получается так, что допуск размером 60 секунд менее 50% интервала, то нужно взять строго допуск = 50% интервала.
	if D_range < (step_in_sec / 2) {
		D_range = (step_in_sec / 2)
		logrus.Printf("prepareData(): D_range = %d is less than step_in_sec / 2 = %d sec, so let have D_range = step_in_sec / 2",
			D_range, (step_in_sec / 2))
	}

	// основной цикл поиска ближайших фактических точек к расчетным
	for i := 0; i < amount+1; i++ {
		point := from.Add(time.Second * time.Duration(step_in_sec*i))
		var item *model.Data_new
		var dto model.DTO

		if i == 0 {
			// для левой границы диапазона ищем ближайшего соседа справа
			item = h.findNearestRightOnly(items, point, D_range)
		} else if i == amount {
			// для правой границы диапазона ищем ближайшего соседа слева
			item = h.findNearestLeftOnly(items, point, D_range)
		} else {
			// для остальных расчетных точек анализируем соседей как слева, так и справа
			item = h.findNearest(items, point, D_range)
		}

		if item != nil {
			dto = model.DTO{
				Temperature: item.Temperature,
				Humidity:    item.Humidity,
				Time:        item.DtWr.Format("2006-01-02 15:04:05"),
			}
		} else {
			// если не нашли фактическую точку в пределах дельта-окрестности - возвращаем фиктивные значения 1000
			// и временную метку расчетной точки
			dto = model.DTO{
				Temperature: 1000,
				Humidity:    1000,
				Time:        point.Format("2006-01-02 15:04:05"),
			}
		}
		result = append(result, dto)
	}

	logrus.Printf("prepareData(): END")

	return result
}

// Реализация алгоритма прореживания данных в режиме отладки
func (h *Handler) prepareDataDebug(items []model.Data_new, amount int, from time.Time, to time.Time, D int) ([]model.DTODebug, int, int) {

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
		logrus.Printf("algo2Debug(): D_range = %d is less than step_in_sec / 2 = %d sec, so let have D_range = step_in_sec / 2",
			D_range, (step_in_sec / 2))
		D_range = (step_in_sec / 2)
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

// Обработка запроса GET host:port/api/data
// Параметры запроса:
// id - числовой идентификатор датчика, данные которого интересуют
// from - начало временного диапазона в формате YYYY-MM-DD HH:mm:ss
// to - конец временного диапазона в формате YYYY-MM-DD HH:mm:ss
// amount - кол-во временных интервалов, на которое нужно разбить временной диапазон
// D - размер дельта-окрестности (в процентах от интервала), в пределах которой ищется точка
// Возвращает список в формате JSON:
// Пример:
//
//	{
//	   "amount": 5,
//	   "data": [
//	       {
//	           "T": 1000,
//	           "H": 1000,
//	           "t": "2021-10-31 23:50:00"
//	       },
//	       {
//	           "T": 1000,
//	           "H": 1000,
//	           "t": "2021-10-31 23:55:00"
//	       }]
//	}
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

	// достаем данные из БД
	items, err := h.services.Data.GetData(parser.ID, parser.From, parser.To, intLimit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		logrus.Errorf("getData(): Cannot get data for id = %d, from = %s, to = %s, error = %s",
			parser.ID, parser.From.Format("2006-01-02 15:04:05"), parser.To.Format("2006-01-02 15:04:05"), err.Error())
		return
	}

	total := len(items)

	var result []model.DTO

	if parser.Amount == 0 {
		// если запрошено нулевое количество данных - возвращаем все данные из БД в заданном диапазоне
		for i := 0; i < total; i++ {
			dto := model.DTO{
				Temperature: items[i].Temperature,
				Humidity:    items[i].Humidity,
				Time:        items[i].DtWr.Format("2006-01-02 15:04:05"),
			}
			result = append(result, dto)
		}
	} else {
		if total < parser.Amount {
			logrus.Printf("getData(): total amount = %d is less than requested amount = %d", total, parser.Amount)
		}
		// запуск алгоритма прореживания данных
		result = h.prepareData(items, parser.Amount, parser.From, parser.To, parser.D)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})

	logrus.Printf("getData(): END")
}

// Обработка запроса GET host:port/api/data/debug
// Параметры запроса:
// id - числовой идентификатор датчика, данные которого интересуют
// from - начало временного диапазона в формате YYYY-MM-DD HH:mm:ss
// to - конец временного диапазона в формате YYYY-MM-DD HH:mm:ss
// amount - кол-во временных интервалов, на которое нужно разбить временной диапазон
// D - размер дельта-окрестности (в процентах от интервала), в пределах которой ищется точка
// Возвращает список в формате JSON c дополнительными данными для отладки
// Пример:
//
//	{
//	   "D_in_sec": 150,
//	   "amount": 5,
//	   "data": [
//	       {
//	           "T": 20.61,
//	           "H": 51.13,
//	           "t": "2021-11-01 00:00:20",
//	           "point": "2021-11-01 00:00:00",
//	           "left": "2021-10-31 23:57:30",
//	           "right": "2021-11-01 00:02:30"
//	       }
//	   ],
//	   "step_in_sec": 300
//	}
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

	// получаем данные из БД
	items, err := h.services.Data.GetData(parser.ID, parser.From, parser.To, intLimit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		logrus.Errorf("getDataDebug(): Cannot get data for id = %d, from = %s, to = %s, error = %s",
			parser.ID, parser.From.Format("2006-01-02 15:04:05"), parser.To.Format("2006-01-02 15:04:05"), err.Error())
		return
	}

	total := len(items)

	var result []model.DTODebug
	var step_in_sec = 0
	var D_in_sec = 0

	if parser.Amount == 0 {
		// если запрошено нулевое количество данных - возвращаем все данные из БД в заданном диапазоне
		for i := 0; i < total; i++ {
			dto := model.DTODebug{
				Temperature: items[i].Temperature,
				Humidity:    items[i].Humidity,
				Time:        items[i].DtWr.Format("2006-01-02 15:04:05"),
			}
			result = append(result, dto)
		}
	} else {

		if total < parser.Amount {
			logrus.Printf("getDataDebug(): total amount = %d is less than requested amount = %d", total, parser.Amount)
		}

		// выполнение алгоритма прореживания данных в режиме отладки
		result, step_in_sec, D_in_sec = h.prepareDataDebug(items, parser.Amount, parser.From, parser.To, parser.D)
	}

	c.JSON(http.StatusOK, gin.H{"step_in_sec": step_in_sec, "D_in_sec": D_in_sec, "data": result})

	logrus.Printf("getDataDebug(): END")
}

// Обработка запроса GET host:port/api/last/:id
// Параметры запроса:
// id - идентификатор датчика
// Возвращает последние доступные в БД значения для заданного датчика в формате JSON
// Пример:
//
//	{
//	   "T": 21.35,
//	   "H": 44.17,
//	   "t": "2023-09-24 23:59:48"
//	}
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
