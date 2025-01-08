package prarthana_data

import (
	"context"
	"fmt"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	mongoCommons "github.com/Out-Of-India-Theory/oit-go-commons/mongo"
	"github.com/Out-Of-India-Theory/prarthana-automated-script/configuration"
	"github.com/Out-Of-India-Theory/prarthana-automated-script/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
)

const (
	prarthana_collection = "prarthanas"
	deity_collection     = "deities"
	shlok_collection     = "shloks"
	stotra_collection    = "stotras"
)

type PrarthanaDataMongoRepository struct {
	logger              *zap.Logger
	prarthanaCollection *mongo.Collection
	deityCollection     *mongo.Collection
	shlokCollection     *mongo.Collection
	stotraCollection    *mongo.Collection
}

func InitPrarthanaDataMongoRepository(ctx context.Context, mongoConfig configuration.MongoConfig) *PrarthanaDataMongoRepository {
	mongoClient := mongoCommons.InitMongoClient(ctx, mongoConfig.MongoConfig)
	return &PrarthanaDataMongoRepository{
		logger:              logging.WithContext(ctx),
		prarthanaCollection: mongoClient.Database(mongoConfig.Database).Collection(prarthana_collection),
		deityCollection:     mongoClient.Database(mongoConfig.Database).Collection(deity_collection),
		shlokCollection:     mongoClient.Database(mongoConfig.Database).Collection(shlok_collection),
		stotraCollection:    mongoClient.Database(mongoConfig.Database).Collection(stotra_collection),
	}
}

func (r *PrarthanaDataMongoRepository) InsertManyShloks(ctx context.Context, shloks []entity.Shlok) error {
	// Prepare the documents for insertion
	var documents []interface{}
	for _, shlok := range shloks {
		documents = append(documents, shlok)
	}

	result, err := r.shlokCollection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("error inserting documents: %w", err)
	}

	log.Printf("Inserted %d documents with IDs: %v\n", len(result.InsertedIDs), result.InsertedIDs)
	return nil
}

func (r *PrarthanaDataMongoRepository) InsertManyStotras(ctx context.Context, stotras []entity.Stotra) error {
	// Prepare the documents for insertion
	var documents []interface{}
	for _, stotra := range stotras {
		documents = append(documents, stotra)
	}

	result, err := r.stotraCollection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("error inserting documents: %w", err)
	}

	log.Printf("Inserted %d documents with IDs: %v\n", len(result.InsertedIDs), result.InsertedIDs)
	return nil
}

func (r *PrarthanaDataMongoRepository) InsertManyDeities(ctx context.Context, deities []entity.DeityDocument) error {
	// Prepare the documents for insertion
	var documents []interface{}
	for _, deity := range deities {
		documents = append(documents, deity)
	}

	result, err := r.deityCollection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("error inserting documents: %w", err)
	}

	log.Printf("Inserted %d documents with IDs: %v\n", len(result.InsertedIDs), result.InsertedIDs)
	return nil
}

func (r *PrarthanaDataMongoRepository) InsertManyPrarthanas(ctx context.Context, prarthanas []entity.Prarthana) error {
	// Prepare the documents for insertion
	var documents []interface{}
	for _, prarthana := range prarthanas {
		documents = append(documents, prarthana)
	}

	result, err := r.prarthanaCollection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("error inserting documents: %w", err)
	}

	log.Printf("Inserted %d documents with IDs: %v\n", len(result.InsertedIDs), result.InsertedIDs)
	return nil
}

func (r *PrarthanaDataMongoRepository) GetTmpIdToPrarthanaIds(ctx context.Context) (map[string]string, map[string]string, error) {
	// Define the filter and projection for the MongoDB query
	filter := bson.M{}
	projection := bson.M{
		"_id":                     1,
		"TmpId":                   1,
		"ui_info.template_number": 1,
	}

	// Query the collection
	cursor, err := r.prarthanaCollection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, nil, fmt.Errorf("error querying the collection: %w", err)
	}
	defer cursor.Close(ctx)

	// Initialize maps to hold the data
	idTemplateMap := make(map[string]string)
	tmpIdToIdMap := make(map[string]string)

	// Iterate through the cursor and decode the documents
	for cursor.Next(ctx) {
		var result struct {
			ID     string `bson:"_id"`
			TmpId  string `bson:"TmpId"`
			UiInfo struct {
				TemplateNumber string `bson:"template_number"`
			} `bson:"ui_info"`
		}

		// Decode the current document into the result struct
		if err := cursor.Decode(&result); err != nil {
			return nil, nil, fmt.Errorf("error decoding document: %w", err)
		}

		// Populate the maps
		idTemplateMap[result.TmpId] = result.UiInfo.TemplateNumber
		tmpIdToIdMap[result.TmpId] = result.ID
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	// Return the populated maps
	return idTemplateMap, tmpIdToIdMap, nil
}

func (r *PrarthanaDataMongoRepository) GetTmpIdToDeityIdMap(ctx context.Context) (map[string]string, error) {
	// Define the filter and projection for the MongoDB query
	filter := bson.M{}
	projection := bson.M{
		"_id":                     1,
		"TmpId":                   1,
		"ui_info.template_number": 1,
	}

	// Query the collection
	cursor, err := r.deityCollection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("error querying the collection: %w", err)
	}
	defer cursor.Close(ctx)

	// Initialize the map to hold the data
	tmpIdToDeityIdMap := make(map[string]string)

	// Iterate through the cursor and decode the documents
	for cursor.Next(ctx) {
		var result struct {
			ID     string `bson:"_id"`
			TmpId  string `bson:"TmpId"`
			UiInfo struct {
				TemplateNumber string `bson:"template_number"`
			} `bson:"ui_info"`
		}

		// Decode the current document into the result struct
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("error decoding document: %w", err)
		}

		// Populate the map
		tmpIdToDeityIdMap[result.TmpId] = result.ID
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor iteration error: %w", err)
	}

	// Return the populated map
	return tmpIdToDeityIdMap, nil
}

func (r *PrarthanaDataMongoRepository) GetAllStotras(ctx context.Context) (map[string]entity.Stotra, error) {
	cursor, err := r.stotraCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching stotras: %w", err)
	}
	defer cursor.Close(ctx)
	stotraMap := make(map[string]entity.Stotra)
	for cursor.Next(ctx) {
		var stotra entity.Stotra
		if err := cursor.Decode(&stotra); err != nil {
			return nil, fmt.Errorf("error decoding stotra: %w", err)
		}
		stotraMap[stotra.ID] = stotra
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return stotraMap, nil
}

func (r *PrarthanaDataMongoRepository) GetAllDeities(ctx context.Context) ([]entity.DeityDocument, error) {
	cursor, err := r.deityCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var deities []entity.DeityDocument
	if err = cursor.All(ctx, &deities); err != nil {
		return nil, err
	}

	return deities, nil
}
