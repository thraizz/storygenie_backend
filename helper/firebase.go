package helper

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
)

var app *firebase.App

func GetFirebaseApp() *firebase.App {
	if app == nil {
		firebaseApp, err := firebase.NewApp(context.Background(), nil)
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
		}
		app = firebaseApp
	}
	log.Default().Println("Firebase app initialized")
	return app

}
