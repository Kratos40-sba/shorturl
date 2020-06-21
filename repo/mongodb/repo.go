package mongodb

import (
	"context"
	"time"

	"github.com/Kratos40-sba/urlshort/shorturl"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepo struct {
	client   *mongo.Client
	database string
	timout   time.Duration
}

func NewMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, err
}
func NewMongoRepo(mongoURL, mongoDB string, mongoTimeout int) (shorturl.RedirectRepo, error) {
	repo := &mongoRepo{
		timout:   time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := NewMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repo.NewMongoRepo")
	}
	repo.client = client
	return repo, nil
}
func (r *mongoRepo) Find(code string) (*shorturl.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timout)
	defer cancel()
	redirect := &shorturl.Redirect{}
	collection := r.client.Database(r.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shorturl.ErrRedirectNotFound, "shorturl.logic")

		}
		return nil, errors.Wrap(err, "repo.Redirect.Find")

	}

	return redirect, nil
}
func (r *mongoRepo) Store(redirect *shorturl.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repo.Redirect.Store")
	}

	return nil
}
