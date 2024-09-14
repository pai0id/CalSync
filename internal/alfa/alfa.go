package alfa

import (
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
)

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

func filterByRoomID[T Lesson](items []T, roomID int) []T {
	var filteredItems []T
	for _, item := range items {
		if item.GetRoomID() == roomID {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func GetAPIToken(email, apiKey string) (string, error) {
	authReq := AuthRequest{Email: email, APIKey: apiKey}
	reqBody, err := json.Marshal(authReq)
	if err != nil {
		return "", err
	}

	log.Println("GetAPIToken: sending request")

	resp, err := http.Post(apiLoginURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("GetAPIToken: request response code %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GetAPIToken: response code %d", resp.StatusCode)
	}

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return "", err
	}

	log.Println("GetAPIToken: token received successfully")

	return authResp.Token, nil
}

// Получение занятий из календаря по его id (если calId == -1 -> все занятия)
func GetLessons(token string, calId int) ([]LessonItem, error) {
	lessonReq := LessonRequest{Status: 1, Page: 0}
	reqBody, err := json.Marshal(lessonReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiLessonURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ALFACRM-TOKEN", token)

	log.Println("GetLessons: sending request")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("GetLessons: request response code %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetAPIToken: response code %d", resp.StatusCode)
	}

	var lessonResp LessonResponse
	err = json.NewDecoder(resp.Body).Decode(&lessonResp)
	if err != nil {
		return nil, err
	}

	log.Println("GetLessons: lessons recieved successfully")

	if calId != -1 {
		return filterByRoomID(lessonResp.Items, roomIds[calId]), nil
	}

	return lessonResp.Items, nil
}

// Получение регулярных занятий из календаря по его id (если calId == -1 -> все занятия)
func GetRegularLessons(token string, calId int) ([]RegularLessonItem, error) {
	regLessonReq := RegularLessonRequest{Page: 0}
	reqBody, err := json.Marshal(regLessonReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiRegularLessonURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ALFACRM-TOKEN", token)

	log.Println("GetRegularLessons: sending request")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("GetLessons: request response code %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetAPIToken: response code %d", resp.StatusCode)
	}

	var regularLessonResp RegularLessonResponse
	err = json.NewDecoder(resp.Body).Decode(&regularLessonResp)
	if err != nil {
		return nil, err
	}

	log.Println("GetRegularLessons: lessons recieved successfully")

	if calId != -1 {
		return filterByRoomID(regularLessonResp.Items, roomIds[calId]), nil
	}

	return regularLessonResp.Items, nil
}
