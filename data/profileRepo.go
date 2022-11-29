package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	// NoSQL: module containing Mongo api client
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NoSQL: ProductRepo struct encapsulating Mongo api client
type ProfileRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment
func New(ctx context.Context, logger *log.Logger) (*ProfileRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &ProfileRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// Disconnect from database
func (pr *ProfileRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (pr *ProfileRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := pr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err)
	}

	// Print available databases
	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
	}
	fmt.Println(databases)
}

func (pr *ProfileRepo) GetAll() (Profiles, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profilesCollection := pr.getCollection()

	var profiles Profiles
	patientsCursor, err := profilesCollection.Find(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	if err = patientsCursor.All(ctx, &profiles); err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return profiles, nil
}

func (pr *ProfileRepo) GetById(id string) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profilesCollection := pr.getCollection()

	var profile Profile
	objID, _ := primitive.ObjectIDFromHex(id)
	err := profilesCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&profile)
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return &profile, nil
}

func (pr *ProfileRepo) GetByName(name string) (Profiles, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profilesCollection := pr.getCollection()

	var profiles Profiles
	profilesCursor, err := profilesCollection.Find(ctx, bson.M{"name": name})
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	if err = profilesCursor.All(ctx, &profiles); err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return profiles, nil
}

func (pr *ProfileRepo) Insert(profile *Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	profilesCollection := pr.getCollection()

	result, err := profilesCollection.InsertOne(ctx, &profile)
	if err != nil {
		pr.logger.Println(err)
		return err
	}
	pr.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (pr *ProfileRepo) Update(id string, profile *Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	profilesCollection := pr.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"name":    profile.Name,
		"surname": profile.Surname,
	}}
	result, err := profilesCollection.UpdateOne(ctx, filter, update)
	pr.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	pr.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		pr.logger.Println(err)
		return err
	}
	return nil
}

func (pr *ProfileRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	profilesCollection := pr.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	result, err := profilesCollection.DeleteOne(ctx, filter)
	if err != nil {
		pr.logger.Println(err)
		return err
	}
	pr.logger.Printf("Documents deleted: %v\n", result.DeletedCount)
	return nil
}

func (pr *ProfileRepo) getCollection() *mongo.Collection {
	profileDatabase := pr.cli.Database("mongoDemo")
	profilesCollectiom := profileDatabase.Collection("patients")
	return profilesCollectiom
}
