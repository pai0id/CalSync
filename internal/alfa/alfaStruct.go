package alfa

type lesson interface {
	getRoomID() int
}

type authRequest struct {
	Email  string `json:"email"`
	APIKey string `json:"api_key"`
}

type authResponse struct {
	Token string `json:"token"`
}

type indexRequest interface {
	getPage() int
	setPage(page int)
}

type lessonRequest struct {
	Status int `json:"status"`
	Page   int `json:"page"`
}

type regularLessonRequest struct {
	Page int `json:"page"`
}

type charRequest struct {
	Page int `json:"page"`
}

type indexResponse interface {
	getTotal() int
	getCount() int
	getItems() []interface{}
}

type lessonResponse struct {
	Total int          `json:"total"`
	Count int          `json:"count"`
	Page  int          `json:"page"`
	Items []lessonItem `json:"items"`
}

type regularLessonResponse struct {
	Total int                 `json:"total"`
	Count int                 `json:"count"`
	Page  int                 `json:"page"`
	Items []regularLessonItem `json:"items"`
}

type charResponse struct {
	Total int        `json:"total"`
	Count int        `json:"count"`
	Page  int        `json:"page"`
	Items []charItem `json:"items"`
}

type lessonItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	RoomId   int    `json:"room_id"`
	TimeFrom string `json:"time_from"`
	TimeTo   string `json:"time_to"`
}

const customerClass = "Customer"
const groupClass = "Group"

type regularLessonItem struct {
	ID           int    `json:"id"`
	RelatedClass string `json:"related_class"`
	RelatedId    int    `json:"related_id"`
	BegDate      string `json:"b_date"`
	EndDate      string `json:"e_date"`
	RoomId       int    `json:"room_id"`
	Day          int    `json:"day"`
	TimeFrom     string `json:"time_from_v"`
	TimeTo       string `json:"time_to_v"`
}

type charItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (l lessonItem) getRoomID() int {
	return l.RoomId
}

func (l regularLessonItem) getRoomID() int {
	return l.RoomId
}

func (req *lessonRequest) setPage(page int) {
	req.Page = page
}

func (req *regularLessonRequest) setPage(page int) {
	req.Page = page
}

func (req *charRequest) setPage(page int) {
	req.Page = page
}

func (req *lessonRequest) getPage() int {
	return req.Page
}

func (req *regularLessonRequest) getPage() int {
	return req.Page
}

func (req *charRequest) getPage() int {
	return req.Page
}

func (resp *lessonResponse) getTotal() int {
	return resp.Total
}

func (resp *lessonResponse) getCount() int {
	return resp.Count
}

func (resp *lessonResponse) getItems() []interface{} {
	var res = make([]interface{}, 0, len(resp.Items))
	for _, item := range resp.Items {
		res = append(res, item)
	}
	return res
}

func (resp *regularLessonResponse) getTotal() int {
	return resp.Total
}

func (resp *regularLessonResponse) getCount() int {
	return resp.Count
}

func (resp *regularLessonResponse) getItems() []interface{} {
	var res = make([]interface{}, 0, len(resp.Items))
	for _, item := range resp.Items {
		res = append(res, item)
	}
	return res
}

func (resp *charResponse) getTotal() int {
	return resp.Total
}

func (resp *charResponse) getCount() int {
	return resp.Count
}

func (resp *charResponse) getItems() []interface{} {
	var res = make([]interface{}, 0, len(resp.Items))
	for _, item := range resp.Items {
		res = append(res, item)
	}
	return res
}
