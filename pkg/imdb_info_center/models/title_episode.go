package models

type TitleEpisode struct {
	TConst        string `gorm:"column:tconst;primary_key;type:varchar(20);not null" json:"tconst"`
	ParentTConst  string `gorm:"column:parent_imdb_id;type:varchar(20);not null" json:"parent_tconst"`
	SeasonNumber  int    `gorm:"column:season_number;type:int;not null" json:"season_number"`
	EpisodeNumber int    `gorm:"column:episode_number;type:int;not null" json:"episode_number"`
}

func NewTitleEpisode(tconst string, parentTconst string, seasonNumber int, episodeNumber int) *TitleEpisode {
	return &TitleEpisode{TConst: tconst, ParentTConst: parentTconst, SeasonNumber: seasonNumber, EpisodeNumber: episodeNumber}
}
