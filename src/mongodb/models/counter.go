package models

type Counter struct {
	ID            string `bson:"_id"`
	SequenceValue int    `bson:"sequence_value"`
}
