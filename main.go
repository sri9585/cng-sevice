package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func init() {
	// Connect to MongoDB
	session, err := mgo.Dial("mongodb+srv://srik090704:sk1234@cluster0.inohjsj.mongodb.net/") // Replace with your MongoDB connection string
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
}

func main() {
	http.HandleFunc("/signup", SignupHandler)
	http.HandleFunc("/login", LoginHandler)

	log.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert user into MongoDB
	session := session.Copy()
	defer session.Close()
	c := session.DB("mydb").C("users") // Replace with your database and collection name

	// Check if the user already exists
	count, err := c.Find(bson.M{"username": user.Username}).Count()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	user.Password = hashPassword(user.Password) // You should hash the password securely
	if err = c.Insert(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authenticate user
	session := session.Copy()
	defer session.Close()
	c := session.DB("mydb").C("users") // Replace with your database and collection name

	var dbUser User
	err := c.Find(bson.M{"username": user.Username}).One(&dbUser)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Verify the password (you should use a secure password hashing library)
	if dbUser.Password != hashPassword(user.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Handle successful login
	response := map[string]string{"message": "Login successful"}
	json.NewEncoder(w).Encode(response)
}

func hashPassword(password string) string {
	// You should implement a secure password hashing function here
	// Do not store plain passwords in the database
	return password
}
