package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Course struct {
	ID         int     `json: "id"`
	Name       string  `json: "name"`
	Price      float64 `json: "price"`
	Instructor string  `json: "instructor"`
}

var courseList []Course

func init() {
	CourseJSON := `[
		{
			"id": 101,
			"name": "Python",
			"price": 2590,
			"instructor": "Tle"
		},
		{
			"id": 102,
			"name": "JavaScript",
			"price": 0,
			"instructor": "Bow"
		},
		{
			"id": 103,
			"name": "SQL",
			"price": 0,
			"instructor": "Men"
		}
	]`
	err := json.Unmarshal([]byte(CourseJSON), &courseList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextId() int {
	highestId := -1
	for _, course := range courseList {
		if highestId < course.ID {
			highestId = course.ID
		}
	}
	return highestId + 1
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	courseJson, err := json.Marshal(courseList)
	switch r.Method {
	case http.MethodGet:
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJson)
	case http.MethodPost:
		var newCourse Course
		Bodybyte, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(Bodybyte, &newCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// เพิ่มตรวจสอบ ID ที่ไม่ใช่ 0 หรือถูกตั้งค่าไว้ใน JSON
		if newCourse.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// กำหนดค่า ID ด้วย getNextId() และเพิ่มคอร์สลงใน courseList
		newCourse.ID = getNextId()
		courseList = append(courseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	}

}

func main() {
	http.HandleFunc("/course", courseHandler)
	http.ListenAndServe(":5000", nil)

}
