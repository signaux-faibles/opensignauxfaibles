package engine

import "github.com/globalsign/mgo/bson"

type Comment struct {
	ID struct {
		Siret  string `json:"key" bson:"key"`
		Author string `json:"author" bson:"author"`
		Date   string `json:"date" bson:"date"`
	} `json:"id" bson:"_id"`
	Comment string `json:"comment" bson:"date"`
}

// GetComments retourne les comments pour un siret
func GetComments(siret string) []Comment {
	var comments []Comment
	Db.DBStatus.C("Comment").Find(bson.M{"_id.key": siret}).All(&comments)
	return comments
}
