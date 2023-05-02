package reminder

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReminderRepositary struct {
	collection *mongo.Collection
}

func NewRepositary(mongodb *mongo.Database) *ReminderRepositary {
	return &ReminderRepositary{collection: mongodb.Collection("reminders")}
}

// retrieves reminder that are close to to remind time
func (repo *ReminderRepositary) GetRemindersToRemind(ctx context.Context) ([]Reminder, error) {
	var reminders []Reminder

	cursor, err := repo.collection.Find(ctx, bson.M{
		"remind_again_on": bson.M{
			"$gt": primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Minute)),
			"$lt": primitive.NewDateTimeFromTime(time.Now().Add(1 * time.Minute)),
		},
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &reminders); err != nil {
		return nil, err
	}
	return reminders, nil
}

func (repo *ReminderRepositary) Upsert(ctx context.Context, reminder Reminder) (Reminder, error) {
	if reminder.ID == primitive.NilObjectID {
		reminder.ID = primitive.NewObjectID()
	}

	filter := bson.M{
		"_id": reminder.ID,
	}
	update := bson.M{
		"$set": reminder,
	}
	var res Reminder
	err := repo.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetUpsert(true), options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (repo *ReminderRepositary) ListUserReminders(ctx context.Context, userDeviceID string) ([]Reminder, error) {
	var reminders []Reminder
	cursor, err := repo.collection.Find(ctx, bson.M{"user_device_id": userDeviceID})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &reminders); err != nil {
		return nil, err
	}
	return reminders, nil
}

func (repo *ReminderRepositary) UpsertMany(ctx context.Context, reminders []Reminder) error {
	operations := make([]mongo.WriteModel, len(reminders))
	for i, reminder := range reminders {
		operations[i] = mongo.NewReplaceOneModel().
			SetUpsert(true).
			SetFilter(bson.M{"_id": reminder.ID}).
			SetReplacement(reminder)
	}

	_, err := repo.collection.BulkWrite(ctx, operations)
	if err != nil {
		return err
	}
	return nil
}
