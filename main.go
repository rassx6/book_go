package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Author      string             `json:"author" bson:"author"`
	Genre       string             `json:"genre" bson:"genre"`
	Description string             `json:"description" bson:"description"`
}

func connectDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB")
	return nil
}

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("bookstoredb").Collection("books")

	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println("Error fetching books:", err)
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())

	var fetchedBooks []Book
	for cur.Next(context.Background()) {
		var book Book
		err := cur.Decode(&book)
		if err != nil {
			log.Println("Error decoding book:", err)
			http.Error(w, "Error decoding book", http.StatusInternalServerError)
			return
		}
		fetchedBooks = append(fetchedBooks, book)
	}

	if err := cur.Err(); err != nil {
		log.Println("Error iterating over books:", err)
		http.Error(w, "Error iterating over books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fetchedBooks)
}

func main() {
	if err := connectDB(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/books", getBooksHandler).Methods("GET")

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.Handle("/api/", http.StripPrefix("/api", router))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
