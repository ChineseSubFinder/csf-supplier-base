package models

type TitleRatings struct {
	TConst        string  `gorm:"column:tconst;primary_key;type:varchar(20);not null" json:"tconst"`
	AverageRating float32 `gorm:"column:average_rating;type:float;not null" json:"average_rating"`
	NumVotes      int     `gorm:"column:num_votes;type:int;not null" json:"num_votes"`
}

func NewTitleRatings(tconst string, averageRating float32, numVotes int) *TitleRatings {
	return &TitleRatings{TConst: tconst, AverageRating: averageRating, NumVotes: numVotes}
}
