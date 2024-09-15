package alfa

import (
	"CalSync/internal/logic"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	apiLoginURL         = "https://musicalwave.s20.online/v2api/auth/login"
	apiLessonURL        = "https://musicalwave.s20.online/v2api/1/lesson/index"
	apiRegularLessonURL = "https://musicalwave.s20.online/v2api/1/regular-lesson/index"
	apiCustomerURL      = "https://musicalwave.s20.online/v2api/1/customer/index"
	apiGroupURL         = "https://musicalwave.s20.online/v2api/1/group/index"
)

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

func sendIndexRequest(token, url string, reqStruct indexRequest, lessonRespStruct indexResponse) ([]interface{}, error) {
	page := 0
	reqStruct.setPage(page)

	err := processIndexRequest(token, url, reqStruct, lessonRespStruct)
	if err != nil {
		return nil, fmt.Errorf("sendIndexRequest: %w", err)
	}

	items := make([]interface{}, 0, lessonRespStruct.getTotal())
	items = append(items, lessonRespStruct.getItems()...)
	total := lessonRespStruct.getTotal() - lessonRespStruct.getCount()

	for total > 0 {
		page++
		reqStruct.setPage(page)

		err = processIndexRequest(token, url, reqStruct, lessonRespStruct)
		if err != nil {
			return nil, fmt.Errorf("sendIndexRequest: %w", err)
		}

		items = append(items, lessonRespStruct.getItems()...)
		total -= lessonRespStruct.getCount()
	}

	return items, nil
}

func convertLessons(lessons []lessonItem) []logic.Lesson {
	var res = []logic.Lesson{}
	for _, l := range lessons {
		res = append(res, logic.Lesson{
			Name:     l.Name,
			Date:     l.Date,
			CalId:    getCalId(l.RoomId),
			TimeFrom: l.TimeFrom,
			TimeTo:   l.TimeTo,
		})
	}
	return res
}

func convertRegularLessons(lessons []regularLessonItem, names []string) []logic.Lesson {
	var res = []logic.Lesson{}
	for i, l := range lessons {
		res = append(res, logic.Lesson{
			Name: names[i],
			// Date:     l.Date,
			CalId:    getCalId(l.RoomId),
			TimeFrom: l.TimeFrom,
			TimeTo:   l.TimeTo,
		})
	}
	return res
}

// Получение занятий из календаря по его id (если calId == -1 -> все занятия)
func GetLessons(token string, calId int) ([]logic.Lesson, error) {
	log.Println("GetLessons: start")

	lessonReq := lessonRequest{Status: 1, Page: 0}
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
	}

	log.Println("GetLessons: end")

	return convertLessons(lessons), nil
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
	}

	charReq := charRequest{Page: 0}
	charResp := charResponse{}

	items, err = sendIndexRequest(token, apiCustomerURL, &charReq, &charResp)
	if err != nil {
		return nil, fmt.Errorf("GetRegularLessons: %w", err)
	}

	log.Println("GetRegularLessons: customers recieved successfully")

	var customers = make(map[int]string, len(items))
	for _, item := range items {
		if item, ok := item.(charItem); ok {
			customers[item.ID] = item.Name
		}
	}

	log.Println("GetRegularLessons: mapped customers")

	charReq = charRequest{Page: 0}
	charResp = charResponse{}

	items, err = sendIndexRequest(token, apiGroupURL, &charReq, &charResp)
	if err != nil {
		return nil, fmt.Errorf("GetRegularLessons: %w", err)
	}

	log.Println("GetRegularLessons: groups recieved successfully")

	var groups = make(map[int]string, len(items))
	for _, item := range items {
		if item, ok := item.(charItem); ok {
			groups[item.ID] = item.Name
		}
	}

	log.Println("GetRegularLessons: mapped groups")

	var names = []string{}
	for _, l := range lessons {
		if l.RelatedClass == customerClass {
			names = append(names, customers[l.RelatedId])
		} else if l.RelatedClass == groupClass {
			names = append(names, groups[l.RelatedId])
		} else {
			names = append(names, l.RelatedClass)
		}
	}

	log.Println("GetRegularLessons: end")

	return convertRegularLessons(lessons, names), nil
}
