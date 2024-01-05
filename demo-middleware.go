package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func findID(ID int) (*Course, int) {
	for i, course := range courseList {
		if course.ID == ID {
			return &course, i
		}

	}
	return nil, 0

}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
	}
	course, listItemIndex := findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("no Course with id %d", ID), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		courseJson, err := json.Marshal(course)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content Type", "application/json")
		w.Write(courseJson)

	case http.MethodPut:
		var updatedCourse Course
		byteBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(byteBody, &updatedCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedCourse.ID != ID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		course = &updatedCourse
		courseList[listItemIndex] = *course
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
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

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before hadler middle start")

		handler.ServeHTTP(w, r)
		fmt.Println("middle finish")

	})
}

func main() {

	courseItemHandler := http.HandlerFunc(courseHandler)
	courseListHandler := http.HandlerFunc(coursesHandler)
	http.Handle("/course/", middlewareHandler(courseItemHandler))
	http.Handle("/course", middlewareHandler(courseListHandler))
	http.ListenAndServe(":5000", nil)

}
