package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User model
type User struct {
	Username string
	Password string
}

var client *mongo.Client

func init() {
	// Initialize MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb+srv://srik090704:sk1234@cluster0.inohjsj.mongodb.net/")
	client, _ = mongo.Connect(context.TODO(), clientOptions)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signup", SignupHandler).Methods("GET", "POST")
	r.HandleFunc("/login", LoginHandler).Methods("GET", "POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse the form data
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Insert the user into the database
		user := User{Username: username, Password: password}
		collection := client.Database("test").Collection("users")
		_, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(w, "Signup Successful!")
	} else {
		// Display the signup form
		tmpl := `<html>
                    <body>
                        <h1>Signup</h1>
                        <form method="post" action="/signup">
                            Username: <input type="text" name="username"><br>
                            Password: <input type="password" name="password"><br>
                            <input type="submit" value="Signup">
                        </form>
                    </body>
                </html>`
		t, _ := template.New("signup").Parse(tmpl)
		t.Execute(w, nil)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse the form data
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Check if the user exists in the database
		collection := client.Database("test").Collection("users")
		var user User
		err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
		if err != nil {
			fmt.Fprintln(w, "Login Failed")
			return
		}

		// Check if the provided password matches the stored password
		if user.Password == password {
			fmt.Fprintln(w, "Login Successful!")
		} else {
			fmt.Fprintln(w, "Login Failed")
		}
	} else {
		// Display the login form
		tmpl := `<html>
                    <body>
                        <h1>Login</h1>
                        <form method="post" action="/login">
                            Username: <input type="text" name="username"><br>
                            Password: <input type="password" name="password"><br>
                            <input type="submit" value="Login">
                        </form>
                    </body>
                </html>`
		t, _ := template.New("login").Parse(tmpl)
		t.Execute(w, nil)
	}
}
