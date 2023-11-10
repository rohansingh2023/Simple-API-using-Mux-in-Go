package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Course struct {
	CourseId    string  `json:"cid"`
	CourseName  string  `json:"cname"`
	CoursePrice string  `json:"cprice"`
	Author      *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

// fakeDB
var courses []Course

// helper
func (c *Course) isEmpty() bool {
	return c.CourseId == "" && c.CourseName == ""
}

// Controllers
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>First Go API</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get All Courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get One Course")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("No course found with the given id")
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create One Course")
	w.Header().Set("Content-Type", "application/json")
	// What if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}
	// what about - {}
	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.isEmpty() {
		json.NewEncoder(w).Encode("No data inside")
		return
	}
	// Generate unique id, string and append it to course
	course.CourseId = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update One Course")
	w.Header().Set("Content-Type", "application/json")

	// Grab id from incoming req
	params := mux.Vars(r)

	// loop through the courses and update it
	for index, cval := range courses {
		if cval.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create One Course")
	w.Header().Set("Content-Type", "application/json")

	// Grab id from incoming req
	params := mux.Vars(r)

	// Loop through the course and delete it
	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode("Course deleted successfully")
}

func main() {
	fmt.Println("Starting the API Server")

	// Initialize
	r := mux.NewRouter()

	// Seeding initial values
	courses = append(courses, Course{
		CourseId:    "10",
		CourseName:  "Starting with Blockchain",
		CoursePrice: "99",
		Author: &Author{
			Fullname: "Rohan Singh",
			Website:  "google.com",
		},
	})

	courses = append(courses, Course{
		CourseId:    "20",
		CourseName:  "AI: The Beginning",
		CoursePrice: "199",
		Author: &Author{
			Fullname: "Hamiz Mathews",
			Website:  "facebook.com",
		},
	})

	// Routes
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/create", createCourse).Methods("POST")
	r.HandleFunc("/change/{id}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/del/{id}", deleteOneCourse).Methods("DELETE")

	// Listen
	log.Fatal(http.ListenAndServe(":4002", r))
}
