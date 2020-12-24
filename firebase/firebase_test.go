package firebase_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"cloud.google.com/go/storage"
	"firebase.google.com/go"
	"google.golang.org/api/option"
)

func TestFireBase(t *testing.T) {
	config := &firebase.Config{
		StorageBucket: "won-beta.appspot.com",
	}

	opt := option.WithCredentialsFile("./firebase_beta_key.json")
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		t.Fatal(err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		t.Fatal(err)
	}

	//fileName := fmt.Sprintf("images/%d/%s", 1797, "id_confirmation")
	fileName := fmt.Sprintf("images/%d/%s", 1797, "driver_license")
	obj := bucket.Object(fileName)
	objAttrs, err := obj.Attrs(context.Background())
	if err != nil {
		if err == storage.ErrObjectNotExist {
			t.Fatal("firebase file missing")
		} else {
			t.Fatal(err)
		}
	}

	t.Log(base64.StdEncoding.EncodeToString(objAttrs.MD5))
}
