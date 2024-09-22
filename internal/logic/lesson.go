package logic

import "time"

type Lesson struct {
	Name     string
	Date     time.Time
	CalId    int
	TimeFrom time.Time
	TimeTo   time.Time
	ToAdd    bool
}

//  primaryKey
// 	Date     time.Time
// 	CalId    int
// 	TimeFrom time.Time

func RemoveCommonElements(arr1, arr2 []Lesson) ([]Lesson, []Lesson) {
	resultArr1 := []Lesson{}
	resultArr2 := []Lesson{}

	contains := func(slice []Lesson, value Lesson) bool {
		for _, v := range slice {
			if v.Date.Equal(value.Date) && v.CalId == value.CalId &&
				v.TimeFrom.Equal(value.TimeFrom) {
				return true
			}
		}
		return false
	}

	for _, v1 := range arr1 {
		if !contains(arr2, v1) {
			resultArr1 = append(resultArr1, v1)
		}
	}

	for _, v2 := range arr2 {
		if !contains(arr1, v2) {
			resultArr2 = append(resultArr2, v2)
		}
	}

	return resultArr1, resultArr2
}

func FilterLessons(lessons []Lesson, maxDate time.Time, minDate time.Time) []Lesson {
	var filteredLessons1 []Lesson
	for _, lesson := range lessons {
		if lesson.Date.Before(maxDate) || lesson.Date.Equal(maxDate) {
			filteredLessons1 = append(filteredLessons1, lesson)
		}
	}
	var filteredLessons2 []Lesson
	for _, lesson := range filteredLessons1 {
		if lesson.Date.After(minDate) || lesson.Date.Equal(minDate) {
			filteredLessons2 = append(filteredLessons2, lesson)
		}
	}
	return filteredLessons2
}
