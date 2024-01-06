package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Course struct {
	CourseID   int     `json: "courseid"`
	CourseName string  `json: "coursename"`
	Price      float64 `json: "price"`
	ImageURL   string  `json: "imageurl"`
}

const coursePath = "courses"

var courseList []Course

const basePath = "/api"

var Db *sql.DB

// ** แสดงข้อมูลทั้งหมด
func getCourseList() ([]Course, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()
	results, err := Db.QueryContext(ctx, `SELECT 
	courseid, 
	coursename, 
	price, 
	image_url 
	FROM courseonline`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	courses := make([]Course, 0)
	for results.Next() {
		var course Course
		results.Scan(&course.CourseID,
			&course.CourseName,
			&course.Price,
			&course.ImageURL)

		courses = append(courses, course)
	}
	return courses, nil

}

// ** เพิ่มข้อมูล
func insertProduct(course Course) (Course, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()
	result, err := Db.ExecContext(ctx, `INSERT INTO courseonline 
	(courseid, 
	coursename, 
	price, 
	image_url 
	) VALUES (?, ?, ?, ?)`,
		course.CourseID,
		course.CourseName,
		course.Price,
		course.ImageURL)
	if err != nil {
		log.Println(err.Error())
		return course, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return course, err
	}
	course.CourseID = int(insertID)
	return course, nil

}

// ** Handle Courses
func handleCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		courseList, err := getCourseList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(courseList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var course Course
		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		insertedCourse, err := insertProduct(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		j, err := json.Marshal(insertedCourse)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// ** getCourse() เอาเฉพาะข้อมูลไอดีนั้นๆมา
func getCourse(courseid int) (*Course, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()
	row := Db.QueryRowContext(ctx, `SELECT 
	courseid, 
	coursename, 
	price, 
	image_url 
	FROM courseonline
	WHERE courseid = ?`, courseid)

	course := &Course{}
	err := row.Scan(
		&course.CourseID,
		&course.CourseName,
		&course.Price,
		&course.ImageURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil

}

// ** removeCourse() ลบเฉพาะข้อมูลไอดี
func removeCourse(courseID int) error {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()
	_, err := Db.ExecContext(ctx, `DELETE FROM courseonline WHERE courseid = ?`, courseID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil

}

// ** Handle Course
func handleCourse(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	courseID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		course, err := getCourse(courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(course)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodDelete:
		err := removeCourse(courseID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// ** CORS เพื่อให้ URL ใช้จากแหล่งที่มาอื่นๆได้
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
		handler.ServeHTTP(w, r)

	})
}

// ** Setup Route Set PATH
func SetupRoutes(apiBasePath string) {
	courseHandler := http.HandlerFunc(handleCourse)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler))
	coursesHandler := http.HandlerFunc(handleCourses)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler))

}

// ** connect to DB "mysql", "id:password@(localhost)/dbName"
func SetupDB() {
	var err error
	Db, err = sql.Open("mysql", "root:@(127.0.0.1:3306)/coursedb")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Db)
	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)

}

func main() {

	SetupDB()
	SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":5000", nil))

}
