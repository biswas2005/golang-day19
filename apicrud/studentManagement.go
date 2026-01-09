package apicrud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Student struct {
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Age   int    `json:"age"`
	Class string `json:"class"`
}

// slice to store
var students []Student

// creates a new student
func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newStudent Student
	json.NewDecoder(r.Body).Decode(&newStudent)
	newStudent.ID = len(students) + 1
	students = append(students, newStudent)
	json.NewEncoder(w).Encode(newStudent)
}

// Reads all the students
func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// Read by ID
func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Id not found")
		return
	}

	for _, student := range students {
		if student.ID == id {
			json.NewEncoder(w).Encode(student)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode("Student not found.")

}

// Update Student
func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//mux.Vars(r) is like map[string]string
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Could not find ID")
		return
	}
	for i, student := range students {
		if student.ID == id {
			students = append(students[:i], students[i+1:]...)

			var updated Student
			err := json.NewDecoder(r.Body).Decode(&updated)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Error updating")
			}
			updated.ID = id
			students = append(students, updated)
			json.NewEncoder(w).Encode(updated)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode("Could not update.")
}

// Delete student
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid ID")
		return
	}
	for i, student := range students {
		if student.ID == id {
			students = append(students[:i], students[i+1:]...)
			json.NewEncoder(w).Encode("Student Deleted")
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode("Failed to Delete.")
}

func StudentManagement() {
	// mux.NewRouter() is like a trafic signal
	//It controls when to execute which method
	r := mux.NewRouter()

	students = append(students, Student{Name: "abc", Age: 12, ID: 1, Class: "CS"})

	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students/{id}", getStudent).Methods("GET")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
