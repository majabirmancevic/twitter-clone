package repository

import (
	"auth-service/model"
	"context"
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

const (
	DATABASE   = "twitter"
	COLLECTION = "users"
)

type UserRepository struct {
	users  *mongo.Collection
	cli    *mongo.Client
	logger *log.Logger
}

func NewUserMongoDBStore(client *mongo.Client) model.UserStore {
	users := client.Database(DATABASE).Collection(COLLECTION)
	return &UserRepository{
		users: users,
	}
}

// Disconnect from database
func (r *UserRepository) Disconnect(ctx context.Context) error {
	err := r.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (r *UserRepository) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := r.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		r.logger.Println(err)
	}

	// Print available databases
	databases, err := r.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		r.logger.Println(err)
	}
	fmt.Println(databases)
}

// NoSQL: Constructor which reads db configuration from environment
func New(ctx context.Context, logger *log.Logger) (*UserRepository, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &UserRepository{
		cli:    client,
		logger: logger,
	}, nil
}

func (r *UserRepository) getCollection() *mongo.Collection {
	userDatabase := r.cli.Database("twitter")
	usersCollection := userDatabase.Collection("users")
	return usersCollection
}

func (r *UserRepository) GetByUsername(username string) (*model.RegularProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := r.getCollection()
	var user model.RegularProfile
	err := usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		r.logger.Println(err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetById(id string) (*model.RegularProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := r.getCollection()

	var user model.RegularProfile
	objID, _ := primitive.ObjectIDFromHex(id)
	err := usersCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		r.logger.Println(err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Insert(user *model.RegularProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	usersCollection := r.getCollection()

	result, err := usersCollection.InsertOne(ctx, &user)
	if err != nil {
		r.logger.Println(err)
		return err
	}
	r.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}
