package repositories

// imports
import (
	"context"
	"log"
	"time"

	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository() domain.UserRepository {
	// setup mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)       // set timeout
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("taskmanager")
	userCol := db.Collection("users")  

	return &UserRepository{Collection: userCol}
}

//  register user in to database
func (userRepo *UserRepository) CreateUser(user *domain.User) error {
	
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	// generate new ObjectID if not set
	if user.ID == "" {
		user.ID = primitive.NewObjectID().Hex()
	}

	// save user to database
	_, err := userRepo.Collection.InsertOne(contx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserExists
		}
		return err
	}

	return nil        // success
}

// find user from database by username
func (userRepo *UserRepository) GetByUsername(username string) (*domain.User, error) {
	
	var user domain.User
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()
	
	// find user by username
	err := userRepo.Collection.FindOne(contx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil        // success
}

// find user from database by id
func (userRepo *UserRepository) GetUserById(userID string) (*domain.User, error) {
	
	objID, err := primitive.ObjectIDFromHex(userID)        // convert string id to ObjectID
	if err != nil {
		return nil,  domain.ErrInvalidUserID
	}

	var user domain.User
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()
	
	// find user by id
	err = userRepo.Collection.FindOne(contx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil         // success
}

// count users in the database currently
func (userRepo *UserRepository) GetUserCount() (int64, error) {
	
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	// count users in user collection currently
	count, err := userRepo.Collection.CountDocuments(contx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil        // success
}

// update user role to admin in database (only admins can perform this operation)
func (userRepo *UserRepository) UpdateRole(userID string, role string) error {
	
	objID, err := primitive.ObjectIDFromHex(userID)        // convert string id to ObjectID
	if err != nil {
		return domain.ErrInvalidUserID
	}

	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	// update user's role to admin
	result, err := userRepo.Collection.UpdateOne(
		contx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"role": role}},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil        // success
}