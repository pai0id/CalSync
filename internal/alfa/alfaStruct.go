package alfa

type Lesson interface {
	GetRoomID() int
}

type AuthRequest struct {
	Email  string `json:"email"`
	APIKey string `json:"api_key"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type LessonRequest struct {
	Status int `json:"status"`
	Page   int `json:"page"`
}

type LessonItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	RoomId   int    `json:"room_id"`
	TimeFrom string `json:"time_from"`
	TimeTo   string `json:"time_to"`
}

type LessonResponse struct {
	Total int          `json:"total"`
	Count int          `json:"count"`
	Page  int          `json:"page"`
	Items []LessonItem `json:"items"`
}

type RegularLessonRequest struct {
	Status int `json:"status"`
	Page   int `json:"page"`
}

type RegularLessonItem struct {
	ID       int    `json:"id"`
	BegDate  string `json:"b_date"`
	EndDate  string `json:"e_date"`
	RoomId   int    `json:"room_id"`
	Day      int    `json:"day"`
	TimeFrom string `json:"time_from_v"`
	TimeTo   string `json:"time_to_v"`
}

type RegularLessonResponse struct {
	Total int                 `json:"total"`
	Count int                 `json:"count"`
	Page  int                 `json:"page"`
	Items []RegularLessonItem `json:"items"`
}

func (l LessonItem) GetRoomID() int {
	return l.RoomId
}

func (l RegularLessonItem) GetRoomID() int {
	return l.RoomId
}
