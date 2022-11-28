package repository

import (
	"auth-service/model"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type AuthRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*AuthRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &AuthRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// Disconnect from database
func (pr *AuthRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (pr *AuthRepo) Ping() {
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

func (pr *AuthRepo) Insert(user *model.RegularProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("-----------Ulazak u bazu")
	usersCollection := pr.getCollection()

	log.Println("----------KORISNICI--------")
	log.Println(usersCollection)
	result, err := usersCollection.InsertOne(ctx, &user)
	log.Println("----rezultat---- ", result)
	log.Println("----eror---- ", err)

	log.Println("upisan korisnik sa ID-om : ", result.InsertedID)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return errors.New("User with that username already exist")
		}
		return err
	}
	pr.logger.Printf("Documents ID: %v\n", result.InsertedID)
	pr.logger.Println(" --------- Kreiran korisnik sa korisnickim imenom ", user.Username)
	return nil
}

func (pr *AuthRepo) GetByUsername(username string) (*model.RegularProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := pr.getCollection()

	var user model.RegularProfile
	err := usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		pr.logger.Println(" -------- Ovaj korisnik ne postoji")
		pr.logger.Println(err)
		return nil, err
	}
	return &user, nil
}

func (pr *AuthRepo) GetByVerificationCode(code string) (*model.RegularProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := pr.getCollection()

	var user model.RegularProfile
	err := usersCollection.FindOne(ctx, bson.M{"verificationCode": code}).Decode(&user)
	if err != nil {
		pr.logger.Println(" -------- Ovaj korisnik ne postoji")
		pr.logger.Println(err)
		return nil, err
	}
	return &user, nil
}

func (pr *AuthRepo) GetAll() (model.RegularProfiles, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := pr.getCollection()

	var users model.RegularProfiles
	usersCursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	if err = usersCursor.All(ctx, &users); err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return users, nil
}

func (pr *AuthRepo) Update(id string, user *model.RegularProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usersCollection := pr.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"verificationCode": user.VerificationCode,
		"verified":         user.Verified,
	}}
	result, err := usersCollection.UpdateOne(ctx, filter, update)
	pr.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	pr.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		pr.logger.Println(err)
		return err
	}
	return nil
}

func (pr *AuthRepo) getCollection() *mongo.Collection {
	userDatabase := pr.cli.Database("twitter")
	userCollection := userDatabase.Collection("users")
	return userCollection
}
