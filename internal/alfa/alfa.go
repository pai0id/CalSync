package alfa

import (
	"CalSync/internal/logic"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	apiLoginURL         = "https://musicalwave.s20.online/v2api/auth/login"
	apiLessonURL        = "https://musicalwave.s20.online/v2api/1/lesson/index"
	apiRegularLessonURL = "https://musicalwave.s20.online/v2api/1/regular-lesson/index"
	apiCustomerURL      = "https://musicalwave.s20.online/v2api/1/customer/index"
	apiGroupURL         = "https://musicalwave.s20.online/v2api/1/group/index"
)

const apiLessonCreateURL = "https://musicalwave.s20.online/v2api/1/lesson/create"

const tokenHeader = "X-ALFACRM-TOKEN"

var roomIds = []int{
	1,  // 01 - Pro
	6,  // 02 - G Labs
	11, // 03 - V Labs
	19, // 04 - V+Pro
	5,  // 06 - VP+Pro
	12, // 07 - VGP
	7,  // 08 - Vocal + G
	16, // Зал loft
	37, // ОПЛОТ - игровая
	38, // ОПЛОТ - мастеровая
}

var customers = map[int]string{}
var groups = map[int]string{}
var types = map[int]string{
	1: "Индивидувльный",
	2: "Групповой",
	3: "Пробный",
	4: "Группа online",
	5: "Аренда",
	6: "Мероприятие",
}

func getCalId(roomId int) int {
	for i, v := range roomIds {
		if v == roomId {
			return i
		}
	}
	return -1
}

func filterByRoomID[T lesson](items []T, roomID int) []T {
	var filteredItems []T
	for _, item := range items {
		if item.getRoomID() == roomID {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func filter[T lesson](items []T) []T {
	var filteredItems []T
	for _, item := range items {
		if getCalId(item.getRoomID()) != -1 {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func GetAPIToken(email, apiKey string) (string, error) {
	authReq := authRequest{Email: email, APIKey: apiKey}
	reqBody, err := json.Marshal(authReq)
	if err != nil {
		return "", fmt.Errorf("GetAPIToken: %w", err)
	}

	log.Println("GetAPIToken: sending request")

	resp, err := http.Post(apiLoginURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("GetAPIToken: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("GetAPIToken: request response code %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GetAPIToken: response code %d", resp.StatusCode)
	}

	var authResp authResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return "", fmt.Errorf("GetAPIToken: %w", err)
	}

	log.Println("GetAPIToken: token received successfully")

	return authResp.Token, nil
}

func processIndexRequest(token, url string, reqStruct indexRequest, respStruct indexResponse) error {
	reqBody, err := json.Marshal(reqStruct)
	if err != nil {
		return fmt.Errorf("processIndexRequest: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("processIndexRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(tokenHeader, token)

	log.Printf("processIndexRequest: sending request, page: %d", reqStruct.getPage())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("processIndexRequest: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("processIndexRequest: request response code %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("processIndexRequest: response code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&respStruct)
	if err != nil {
		return fmt.Errorf("processIndexRequest: %w", err)
	}

	return nil
}

func sendIndexRequest(token, url string, reqStruct indexRequest, respStruct indexResponse) ([]interface{}, error) {
	page := 0
	reqStruct.setPage(page)

	err := processIndexRequest(token, url, reqStruct, respStruct)
	if err != nil {
		return nil, fmt.Errorf("sendIndexRequest: %w", err)
	}

	items := make([]interface{}, 0, respStruct.getTotal())
	items = append(items, respStruct.getItems()...)
	total := respStruct.getTotal() - respStruct.getCount()

	for total > 0 {
		page++
		reqStruct.setPage(page)

		err = processIndexRequest(token, url, reqStruct, respStruct)
		if err != nil {
			return nil, fmt.Errorf("sendIndexRequest: %w", err)
		}

		items = append(items, respStruct.getItems()...)
		total -= respStruct.getCount()
	}

	return items, nil
}

func convertLessons(lessons []lessonItem, names []string) ([]logic.Lesson, error) {
	var res = []logic.Lesson{}
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, fmt.Errorf("error loading location: %w", err)
	}
	for i, l := range lessons {
		date, err := time.Parse("2006-01-02", l.Date)
		if err != nil {
			return nil, fmt.Errorf("convertLessons: %w", err)
		}
		timeFrom, err := time.Parse("2006-01-02 15:04:05", l.TimeFrom)
		if err != nil {
			return nil, fmt.Errorf("convertLessons: %w", err)
		}
		timeTo, err := time.Parse("2006-01-02 15:04:05", l.TimeTo)
		if err != nil {
			return nil, fmt.Errorf("convertLessons: %w", err)
		}
		newLesson := logic.Lesson{
			Name:     names[i],
			Date:     date.Truncate(24 * time.Hour).In(location).Add(time.Hour * -3),
			CalId:    getCalId(l.RoomId),
			TimeFrom: timeFrom.In(location).Add(time.Hour * -3).Truncate(time.Minute),
			TimeTo:   timeTo.In(location).Add(time.Hour * -3).Truncate(time.Minute),
		}
		if strings.HasPrefix(l.Note, "Google") {
			newLesson.ToAdd = false
		} else {
			newLesson.ToAdd = true
		}

		res = append(res, newLesson)
	}
	return res, nil
}

func getDatesForDayOfWeek(start, end time.Time, targetDay time.Weekday) []time.Time {
	var dates []time.Time
	var t = time.Now().Truncate(24 * time.Hour)

	if start.Before(t) {
		start = t
	}
	if end.After(t.AddDate(0, 1, 0)) {
		end = t.AddDate(0, 1, 0)
	}

	if start.After(end) {
		return dates
	}

	for start.Weekday() != targetDay {
		start = start.AddDate(0, 0, 1)
	}

	for !start.After(end) {
		dates = append(dates, start)
		start = start.AddDate(0, 1, 0)
	}

	return dates
}

func convertRegularLessons(lessons []regularLessonItem, names []string) ([]logic.Lesson, error) {
	var res = []logic.Lesson{}
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, fmt.Errorf("error loading location: %w", err)
	}
	for i, l := range lessons {
		dateFrom, err := time.Parse("2006-01-02", l.BegDate)
		if err != nil {
			return nil, fmt.Errorf("convertLessons: %w", err)
		}
		dateTo, err := time.Parse("2006-01-02", l.EndDate)
		if err != nil {
			return nil, fmt.Errorf("convertLessons: %w", err)
		}
		for _, date := range getDatesForDayOfWeek(dateFrom, dateTo, time.Weekday(l.Day%7)) {
			timeFrom, err := time.Parse("2006-01-02 15:04", date.Format("2006-01-02")+" "+l.TimeFrom)
			if err != nil {
				return nil, fmt.Errorf("convertLessons: %w", err)
			}
			timeTo, err := time.Parse("2006-01-02 15:04", date.Format("2006-01-02")+" "+l.TimeTo)
			if err != nil {
				return nil, fmt.Errorf("convertLessons: %w", err)
			}
			res = append(res, logic.Lesson{
				Name:     names[i],
				Date:     date.Truncate(24 * time.Hour).In(location).Add(time.Hour * -3),
				CalId:    getCalId(l.RoomId),
				TimeFrom: timeFrom.In(location).Add(time.Hour * -3).Truncate(time.Minute),
				TimeTo:   timeTo.In(location).Add(time.Hour * -3).Truncate(time.Minute),
			})
		}
	}
	return res, nil
}

// Получение занятий из календаря по его id (если calId == -1 -> все занятия)
func GetLessons(token string, calId int) ([]logic.Lesson, error) {
	t := time.Now().Truncate(24 * time.Hour)
	log.Println("GetLessons: start")
	var res []logic.Lesson

	lessonReq := lessonRequest{Status: 1, Page: 0,
		DateTo:   t.AddDate(0, 1, 0).Format("2006-01-02"),
		DateFrom: t.Format("2006-01-02"),
	}
	lessonResp := lessonResponse{}

	items, err := sendIndexRequest(token, apiLessonURL, &lessonReq, &lessonResp)
	if err != nil {
		return nil, fmt.Errorf("GetLessons: %w", err)
	}

	log.Println("GetLessons: lessons recieved successfully")

	var lessons = make([]lessonItem, 0, len(items))
	for _, item := range items {
		if item, ok := item.(lessonItem); ok {
			lessons = append(lessons, item)
		}
	}

	if calId != -1 {
		lessons = filterByRoomID(lessons, roomIds[calId])
	} else {
		lessons = filter(lessons)
	}

	var names = make([]string, 0, len(lessons))
	for _, l := range lessons {
		// name := ""
		// for _, c := range l.Customers {
		// 	name = fmt.Sprintf("%s %s", name, customers[c])
		// }
		// for _, g := range l.Groups {
		// 	name = fmt.Sprintf("%s %s", name, groups[g])
		// }
		// if name == "" {
		// 	name = "Занятие"
		// }
		name := fmt.Sprintf("Alfa: %d, %s", l.ID, l.Note)
		if v, ok := types[l.TypeID]; ok {
			name = fmt.Sprintf("%s (%s)", name, v)
		}
		names = append(names, name)
	}

	log.Println("GetLessons: convert")

	res, err = convertLessons(lessons, names)
	if err != nil {
		return nil, fmt.Errorf("GetLessons: %w", err)
	}

	log.Println("GetLessons: end")

	return res, nil
}

// Получение регулярных занятий из календаря по его id (если calId == -1 -> все занятия)
func GetRegularLessons(token string, calId int) ([]logic.Lesson, error) {
	log.Println("GetRegularLessons: start")

	regLessonReq := regularLessonRequest{Page: 0}
	regLessonResp := regularLessonResponse{}

	items, err := sendIndexRequest(token, apiRegularLessonURL, &regLessonReq, &regLessonResp)
	if err != nil {
		return nil, fmt.Errorf("GetRegularLessons: %w", err)
	}

	log.Println("GetRegularLessons: lessons recieved successfully")

	var lessons = make([]regularLessonItem, 0, len(items))
	for _, item := range items {
		if item, ok := item.(regularLessonItem); ok {
			lessons = append(lessons, item)
		}
	}

	if calId != -1 {
		lessons = filterByRoomID(lessons, roomIds[calId])
	} else {
		lessons = filter(lessons)
	}

	var names = make([]string, 0, len(lessons))
	for _, l := range lessons {
		if l.RelatedClass == customerClass {
			names = append(names, customers[l.RelatedId])
		} else if l.RelatedClass == groupClass {
			names = append(names, groups[l.RelatedId])
		} else {
			names = append(names, l.RelatedClass)
		}
	}

	log.Println("GetRegularLessons: convert")

	res, err := convertRegularLessons(lessons, names)
	if err != nil {
		return nil, fmt.Errorf("GetRegularLessons: %w", err)
	}

	log.Println("GetRegularLessons: end")

	return res, nil
}

func UpdateCustomers(token string) error {
	charReq := charRequest{Page: 0}
	charResp := charResponse{}

	items, err := sendIndexRequest(token, apiCustomerURL, &charReq, &charResp)
	if err != nil {
		return fmt.Errorf("UpdateCustomers: %w", err)
	}

	log.Println("UpdateCustomers: customers recieved successfully")

	for _, item := range items {
		if item, ok := item.(charItem); ok {
			if _, exists := customers[item.ID]; !exists {
				customers[item.ID] = item.Name
			}
		}
	}

	log.Println("UpdateCustomers: mapped customers")
	return nil
}

func UpdateGroups(token string) error {
	charReq := charRequest{Page: 0}
	charResp := charResponse{}

	items, err := sendIndexRequest(token, apiGroupURL, &charReq, &charResp)
	if err != nil {
		return fmt.Errorf("UpdateGroups: %w", err)
	}

	log.Println("UpdateGroups: groups recieved successfully")

	for _, item := range items {
		if item, ok := item.(charItem); ok {
			if _, exists := groups[item.ID]; !exists {
				groups[item.ID] = item.Name
			}
		}
	}

	log.Println("UpdateGroups: mapped groups")
	return nil
}

func createLesson(token string, lesson logic.Lesson) error {
	reqStruct := createLessonRequest{
		Topic:        strings.TrimSpace(lesson.Name),
		LessonDate:   lesson.Date.Format("02.01.2006"),
		RoomId:       roomIds[lesson.CalId],
		TimeFrom:     lesson.TimeFrom.Format("15:04"),
		Duration:     int(lesson.TimeTo.Sub(lesson.TimeFrom).Minutes()) - 1,
		TeacherIds:   []int{55},
		LessonTypeId: 5,
		SubjectId:    118,
	}

	reqBody, err := json.Marshal(reqStruct)
	if err != nil {
		return fmt.Errorf("createLesson: %w", err)
	}

	req, err := http.NewRequest("POST", apiLessonCreateURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("createLesson: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(tokenHeader, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("createLesson: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Created: %v\n", reqStruct)
	log.Printf("createLesson: request response code %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("createLesson: response code %d", resp.StatusCode)
	}

	return nil
}

func AddEvents(token string, lessons []logic.Lesson) error {
	for _, lesson := range lessons {
		if !lesson.ToAdd {
			continue
		}
		err := createLesson(token, lesson)
		if err != nil {
			fmt.Printf("AddEvents: %v", err)
		}
	}
	return nil
}
