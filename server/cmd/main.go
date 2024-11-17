package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type LikeNotification struct {
    EventType    string `json:"event_type"`
    PostID       string `json:"post_id"`
    PosterUserID string `json:"poster_user_id"`
    LikerUserID  string `json:"liker_user_id"`
    LikerUsername string `json:"liker_username"`
    Timestamp    string `json:"timestamp"`
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var notification LikeNotification
    err := json.NewDecoder(r.Body).Decode(&notification)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    notification.Timestamp = time.Now().Format(time.RFC3339)
    log.Printf("Received like notification: %+v\n", notification)

    // Here you would typically send the notification to Kafka or another message queue

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Notification received"))
}

func main() {
    http.HandleFunc("/like", likeHandler)
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err.Error())
    }
}
