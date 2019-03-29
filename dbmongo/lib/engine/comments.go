package engine

import "github.com/globalsign/mgo/bson"

// Comment est le type des annotations attachées aux établissements
type Comment struct {
	ID struct {
		Siret   string `json:"key" bson:"key"`
		Author  string `json:"author" bson:"author"`
		Date    string `json:"date" bson:"date"`
		Version int    `json:"version" bson:"version"`
	} `json:"id" bson:"_id"`
	Comment string `json:"comment" bson:"date"`
	History bool   `json:"history" bson:"history"`
}

// GetComments retourne les comments pour un siret
func GetComments(siret string) ([]Comment, error) {
	var comments []Comment
	err := Db.DBStatus.C("Comment").Find(bson.M{"_id.siret": siret, "_id.version": 0}).All(&comments)
	return comments, err
}

// GetCommentHistory retourne les versions précédentes d'une annotation
func GetCommentHistory(comment Comment) ([]Comment, error) {
	var comments []Comment
	err := Db.DBStatus.C("Comment").Find(bson.M{
		"_id.siret":  comment.ID.Siret,
		"_id.author": comment.ID.Author,
		"_id.date":   comment.ID.Date}).All(&comments)
	return comments, err
}

// SetComment insert ou met à jour un commentaire pour un siret
func SetComment(comment Comment) error {
	var comments []Comment
	err := Db.DBStatus.C("Comment").Find(
		bson.M{
			"_id.siret":  comment.ID.Siret,
			"_id.date":   comment.ID.Date,
			"_id.author": comment.ID.Author,
		},
	).Sort("_id.version").All(&comments)

	if err != nil {
		return err
	}

	comment.ID.Version = 0
	comment.History = false
	newComments := []Comment{comment}
	for _, c := range comments {
		c.ID.Version++
		c.History = true
		newComments = append(newComments, c)
	}

	for _, c := range newComments {
		_, err := Db.DBStatus.C("Comment").Upsert(c.ID, c)
		if err != nil {
			return err
		}
	}

	err = Db.DBStatus.C("Comment").Insert(comment.ID, comment)
	return err
}
