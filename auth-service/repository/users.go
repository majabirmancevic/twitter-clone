package repository

import (
	"auth-service/model"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	DATABASE   = "users"
	COLLECTION = "users"
)

type UserRepository struct {
	users *mongo.Collection
}

func NewUserMongoDBStore(client *mongo.Client) model.UserStore {
	users := client.Database(DATABASE).Collection(COLLECTION)
	return &UserRepository{
		users: users,
	}
}

func (r *UserRepository) Get(id primitive.ObjectID) (*model.RegularProfile, error) {
	filter := bson.M{"_id": id}
	return r.filterOne(filter)
}

func (r *UserRepository) Insert(order *model.RegularProfile) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) filterOne(filter interface{}) (User *model.RegularProfile, err error) {
	result := r.users.FindOne(context.TODO(), filter)
	err = result.Decode(&User)
	return
}

func decode(cursor *mongo.Cursor) (users []*model.RegularProfile, err error) {
	for cursor.Next(context.TODO()) {
		var User model.RegularProfile
		err = cursor.Decode(&User)
		if err != nil {
			return
		}
		users = append(users, &User)
	}
	err = cursor.Err()
	return
}
