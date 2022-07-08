package model

type DatabaseConnection struct {
	ID        string  `bson:"id"`
	Username1 string  `bson:"username1"`
	Username2 string  `bson:"username2"`
	Debt      float64 `bson:"debt"` // How much is owed by username1 to username2
	CreatedAt string  `bson:"created_at"`
}
