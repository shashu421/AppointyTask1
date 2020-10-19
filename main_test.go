//UNIT TESTING

package main

import(
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
)

func TestGetMeetings(t *testing.T) {
	req, err := http.NewRequest("GET", "/meetings", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(meetings)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[
		{
		  "title": "Rreifdf",
		  "start": "No",
		  "end": "shashank.d2018@vitstudent.ac.in",
		  "participants": [
			{
			  "name": "shashu",
			  "email": "shashank.d2018@vitstudent.ac.in",
			  "rsvp": "YES"
			},
			{
			  "name": "shashu",
			  "email": "shashank.d2018@vitstudent.ac.in",
			  "rsvp": "YES"
			}
		  ]
		},
		{
		  "title": "Rreifdf",
		  "start": "No",
		  "end": "shashank.d2018@vitstudent.ac.in",
		  "participants": [
			{
			  "name": "shaskhan",
			  "email": "shashank@vitstudent.ac.in",
			  "rsvp": "YES"
			},
			{
			  "name": "sasaas",
			  "email": "shashank@vitstudent.ac.in",
			  "rsvp": "YES"
			}
		  ]
		}
	  ]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetMeetingByID(t *testing.T) {

	req, err := http.NewRequest("GET", "/meetings", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("5f8d430cc7f8bdbfafbb99b6")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getmeeting)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{
		"title": "Hanuman",
		"start": "No",
		"end": "shashank2018@vitstudent.ac.in",
		"participants": [
		  {
			"name": "Punjabi",
			"email": "shashank@vitstudent.ac.in",
			"rsvp": "YES"
		  },
		  {
			"name": "Shashank",
			"email": "shashank@vitstudent.ac.in",
			"rsvp": "YES"
		  }
		]
	  }`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}


func TestCreateMeeting(t *testing.T) {

	var jsonStr = []byte(`{
    
        "title":"Hanuman",
        "end":"hanuman2018@vitstudent.ac.in",
        "start":"No",
        "participants":[{
            "name":"Punjabi",
            "email":"shashank@vitstudent.ac.in",
            "rsvp":"YES"
        },{
            "name":"Shashank",
            "email":"shashank@vitstudent.ac.in",
            "rsvp":"YES"
        }
        ]
}`)

	req, err := http.NewRequest("POST", "/meetings", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(meetings)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{
		"InsertedID": "5f8d4507c7f8bdbfafbb99b7"
	  }`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
