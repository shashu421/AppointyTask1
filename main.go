package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
 	"go.mongodb.org/mongo-driver/bson/primitive"
 	"go.mongodb.org/mongo-driver/mongo"
 	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

//Structure for Meeting and Participants
type Meeting struct {
	Title        string `json:"title"`
	Start_time string `json:"start"`
	End_time          string `json:"end"`
	Created_at      time.Time `json:"creation"`
	Participants  []Participants `json:"participants" bson:"participants"`
}
type Participants struct {
	Name string `json:"name"`
	Email string `json:"email"`
	RSVP string `json:"rsvp"`
}

//for /meetings route we are handling get and post both from this function
func meetings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		post(w,r)
		return
	case "GET":
		getmeetingrt(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func post(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
		fmt.Println(r.Body);
		var meeting Meeting
		_ = json.NewDecoder(r.Body).Decode(&meeting)
		fmt.Println(meeting)
		collection := client.Database("meeting_scheduler").Collection("Meeting")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		result, _ := collection.InsertOne(ctx, meeting)
		json.NewEncoder(w).Encode(result)
}

func getmeetingrt(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")
	start_time := r.FormValue("start")
	end_time := r.FormValue("end")
	page ,err:= strconv.Atoi(r.FormValue("page"))
	a :=6
	v := page*a

	fmt.Println(start_time)
	fmt.Println(end_time)
	var meetings []Meeting
	collection := client.Database("meeting_scheduler").Collection("Meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		meetings = append(meetings, meeting)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	//PAGINATION
	var output []Meeting
	fmt.Println(meetings)
	for i :=v;i<v+a;i++{
		if i>=len(meetings){
			break;
		} else{
		output=append(output,meetings[i])
		}
	}
	json.NewEncoder(w).Encode(output)

}

//function for handling /meeting/{id} get route
func getmeeting(w http.ResponseWriter,r *http.Request){
	switch r.Method {
	case "GET":
		get(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}


func get(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	fmt.Println(parts)
	fmt.Println(parts[2])
	id, _ := primitive.ObjectIDFromHex(parts[2])
	fmt.Println(id)
	var meeting Meeting
	collection := client.Database("meeting_scheduler").Collection("Meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&meeting)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(meeting)
}

//for handling our 4th route i.e /meeting/participants=<email_id> (get method)
func newfunction(w http.ResponseWriter,r *http.Request){
	switch r.Method {
	case "GET":
		getparticipants(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func getparticipants(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")
	var participants [][]Participants
	var meetings []Meeting
	collection := client.Database("meeting_scheduler").Collection("Meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	pemailid := r.FormValue("participant")
	fmt.Println(pemailid);
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		fmt.Println(meeting.Participants)
		for i := 0; i < len(meeting.Participants); i++ {
			x := meeting.Participants[i]
			if x.Email == pemailid {
				participants = append(participants,meeting.Participants)
				meetings = append(meetings, meeting)
			}
		}
		
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(meetings)
}

//our entry point our main function
func main() {
	fmt.Println("Server has started on port 3000")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	//(1) and (3)route for creating a meeting and retrieving all meetings with given constraints
	http.HandleFunc("/meetings", meetings)

	//(2)route for getting a meeting by ID
	http.HandleFunc("/meetings/",getmeeting)

	//(4)routes for getting all meetings of a participant using email id(List all meetings of Participant)
	http.HandleFunc("/meeting",newfunction)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}




