package logic

import "time"

type Lesson struct {
	Name     string
	Date     time.Time
	CalId    int
	TimeFrom time.Time
	TimeTo   time.Time
}

//  primaryKey
// 	Date     time.Time
// 	CalId    int
// 	TimeFrom time.Time
// 	TimeTo   time.Time

func RemoveCommonElements(arr1, arr2 []Lesson) ([]Lesson, []Lesson) {
	resultArr1 := []Lesson{}
	resultArr2 := []Lesson{}

	contains := func(slice []Lesson, value Lesson) bool {
		for _, v := range slice {
			if v.Date.Equal(value.Date) && v.CalId == value.CalId &&
				v.TimeFrom.Equal(value.TimeFrom) && v.TimeTo.Equal(value.TimeTo) {
				return true
			}
		}
		return false
	}

	commonElements := []Lesson{}
	for _, v1 := range arr1 {
		if contains(arr2, v1) && !contains(commonElements, v1) {
			commonElements = append(commonElements, v1)
		}
	}

	for _, v1 := range arr1 {
		if !contains(commonElements, v1) {
			resultArr1 = append(resultArr1, v1)
		}
	}

	for _, v2 := range arr2 {
		if !contains(commonElements, v2) {
			resultArr2 = append(resultArr2, v2)
		}
	}

	return resultArr1, resultArr2
}
