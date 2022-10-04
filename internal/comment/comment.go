package comment

import "database/sql"

type CommentDataRepo struct {
	DB *sql.DB
}

type AUserComment struct {
	Name    string `db:"name"`
	Comment string `db:"comment"`
}

type AUserCommentJson struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}
