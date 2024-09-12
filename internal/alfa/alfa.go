package alfa

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	apiLoginURL         = "https://musicwave.s20.online/v2api/auth/login"
	apiLessonURL        = "https://musicwave.s20.online/v2api/1/lesson/index"
	apiRegularLessonURL = "https://musicwave.s20.online/v2api/1/regular-lesson/index"
)

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

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return "", err
	}

	log.Println("GetAPIToken: token received successfully")

	return authResp.Token, nil
}

func GetLessons(token string) ([]LessonItem, error) {
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

	var lessonResp LessonResponse
	err = json.NewDecoder(resp.Body).Decode(&lessonResp)
	if err != nil {
		return nil, err
	}

	log.Println("GetLessons: lessons recieved successfully")

	return lessonResp.Items, nil
}

func GetRegularLessons(token string) ([]RegularLessonItem, error) {
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

	var regularLessonResp RegularLessonResponse
	err = json.NewDecoder(resp.Body).Decode(&regularLessonResp)
	if err != nil {
		return nil, err
	}

	log.Println("GetRegularLessons: lessons recieved successfully")

	return regularLessonResp.Items, nil
}
