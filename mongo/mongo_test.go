package mongo_test

import (
    "context"
    "fmt"
    "log"
    "testing"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func setClient() *mongo.Client {
    // create a new context
    ctx := context.Background()

    // create a mongo client
    client, err := mongo.Connect(
        ctx,
        options.Client().ApplyURI("mongodb://localhost:27017/"),
    )
    if err != nil {
        log.Fatal(err)
    }

    return client
}

func TestConnection(t *testing.T) {
    // create a new context
    ctx := context.Background()

    // create a mongo client
    client, err := mongo.NewClient(
        options.Client().ApplyURI("mongodb://localhost:27017/"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // connect to mongo
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }

    // disconnects from mongo
    defer client.Disconnect(ctx)
}

func TestConnectionSimple(t *testing.T) {
    // create a new context
    ctx := context.Background()

    // create a mongo client
    client, err := mongo.Connect(
        ctx,
        options.Client().ApplyURI("mongodb://localhost:27017/"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // disconnects from mongo
    defer client.Disconnect(ctx)
}

type Post struct {
    ID        primitive.ObjectID `bson:"_id"`
    Title     string             `bson:"title"`
    Body      string             `bson:"body"`
    Tags      []string           `bson:"tags"`
    Comments  uint64             `bson:"comments"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

func TestInsertOne(t *testing.T) {
    client := setClient()

    db := client.Database("blog")
    col := db.Collection("posts")

    // create a new context with a 10 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // insert one document
    res, err := col.InsertOne(ctx, bson.M{
        "title": "Go mongodb driver cookbook",
        "tags":  []string{"golang", "mongodb"},
        "body": `this is a long post
that goes on and on
and have many lines`,
        "comments":   1,
        "created_at": time.Now(),
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf(
        "new post created with id: %s",
        res.InsertedID.(primitive.ObjectID).Hex(),
    )
}

func TestInsertMany(t *testing.T) {
    client := setClient()

    db := client.Database("blog")
    col := db.Collection("posts")

    // create a new context with a 10 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    res, err := col.InsertMany(ctx, []interface{}{
        bson.M{
            "title":      "Post one",
            "tags":       []string{"golang"},
            "body":       "post one body",
            "comments":   14,
            "created_at": time.Date(2019, time.January, 10, 15, 30, 0, 0, time.UTC),
        },
        bson.M{
            "title":      "Post two",
            "tags":       []string{"nodejs"},
            "body":       "post two body",
            "comments":   2,
            "created_at": time.Now(),
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("inserted ids: %v\n", res.InsertedIDs)
}

func TestUpdateOne(t *testing.T) {
    client := setClient()

    db := client.Database("blog")
    col := db.Collection("posts")

    // create a new context with a 10 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // create ObjectID from string
    id, err := primitive.ObjectIDFromHex("5e4bac00b783266161f0c372")
    if err != nil {
        log.Fatal(err)
    }

    // set filters and updates
    filter := bson.M{"_id": id}
    update := bson.M{"$set": bson.M{"title": "post 2 (two)"}}

    // update document
    res, err := col.UpdateOne(ctx, filter, update)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("modified count: %d\n", res.ModifiedCount)
}

func TestUpdateMany(t *testing.T) {
    client := setClient()

    db := client.Database("blog")
    col := db.Collection("posts")

    // create a new context with a 10 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // set filters and updates
    filter := bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "golang"}}}
    update := bson.M{"$set": bson.M{"comments": 0, "updated_at": time.Now()}}

    // update documents
    res, err := col.UpdateMany(ctx, filter, update)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("modified count: %d\n", res.ModifiedCount)
}

func TestFindOne(t *testing.T) {
    client := setClient()

    db := client.Database("blog")
    col := db.Collection("posts")

    // create a new context with a 10 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // filter posts tagged as golang
    filter := bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "golang"}}}

    // find one document
    var p Post
    if err := col.FindOne(ctx, filter).Decode(&p); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("post: %+v\n", p)
}

func TestFindMany(t *testing.T) {
    client := setClient()

    db := client.Database("blog")
    col := db.Collection("posts")

    // create a new context with a 10 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // filter posts tagged as golang
    filter := bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "golang"}}}

    // find all documents
    cursor, err := col.Find(ctx, filter)
    if err != nil {
        log.Fatal(err)
    }

    // iterate through all documents
    for cursor.Next(ctx) {
        var p Post
        // decode the document
        if err := cursor.Decode(&p); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("post: %+v\n", p)
    }

    // check if the cursor encountered any errors while iterating
    if err := cursor.Err(); err != nil {
        log.Fatal(err)
    }
}