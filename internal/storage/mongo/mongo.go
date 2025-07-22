package storagemongo

import (
	"context"
	"sso/internal/domain/models"
	"sso/internal/storage"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	_ "go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *mongo.Client
}

func New(ctx context.Context, connectionURI string) (*Storage, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Storage{db: client}, nil
}

func (s *Storage) Stop(ctx context.Context) {
	s.db.Disconnect(ctx)
}

func (s *Storage) CreateUser(ctx context.Context, email string, password string) (uid string, err error) {
	coll := s.db.Database("sso").Collection("users")

	var checkUser models.User
	coll.FindOne(ctx, bson.M{"email": email}).Decode(&checkUser)
	if !checkUser.ID.IsZero() {
		return "", storage.ErrUserExists
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := models.User{
		Email: email,
		Password: bson.Binary{
			Subtype: 0,
			Data:    passHash,
		}}

	res, err := coll.InsertOne(ctx, newUser)
	if err != nil {
		return "", err
	}
	uid = res.InsertedID.(bson.ObjectID).Hex()

	return uid, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	coll := s.db.Database("sso").Collection("users")

	var res models.User
	err = coll.FindOne(ctx, bson.M{"email": email}).Decode(&res)
	if res.ID.IsZero() {
		return models.User{}, storage.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	return res, nil
}


func (s *Storage) GetAppById(ctx context.Context, id string) (app models.App, err error) {
	coll := s.db.Database("sso").Collection("apps")
	
	appId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.App{}, err
	}

	var res models.App
	err = coll.FindOne(ctx, bson.M{"_id": appId}).Decode(&res)
	if res.ID.IsZero() {
		return models.App{}, storage.ErrAppNotFound
	}
	if err != nil {
		return models.App{}, err
	}

	return res, nil
}
