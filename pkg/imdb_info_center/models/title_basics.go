package models

type TitleBasic struct {
	TConst         string `gorm:"column:tconst;primary_key;type:varchar(20);not null" json:"tconst"`
	TitleType      string `gorm:"column:title_type;type:varchar(255);not null" json:"title_type"`
	PrimaryTitle   string `gorm:"column:primary_title;type:varchar(255);not null" json:"primary_title"`
	OriginalTitle  string `gorm:"column:original_title;type:varchar(255);not null" json:"original_title"`
	IsAdult        bool   `gorm:"column:is_adult;type:tinyint;not null" json:"is_adult"`
	StartYear      string `gorm:"column:start_year;type:varchar(10);not null" json:"start_year"`
	EndYear        string `gorm:"column:end_year;type:varchar(10);not null" json:"end_year"`
	RuntimeMinutes string `gorm:"column:runtime_minutes;type:varchar(10);not null" json:"runtime_minutes"`
	Genres         string `gorm:"column:genres;type:varchar(255);not null" json:"genres"`
}

func NewTitleBasic(tconst string, titleType string, primaryTitle string, originalTitle string, isAdult bool, startYear string, endYear string, runtimeMinutes string, genres string) *TitleBasic {
	return &TitleBasic{TConst: tconst, TitleType: titleType, PrimaryTitle: primaryTitle, OriginalTitle: originalTitle, IsAdult: isAdult, StartYear: startYear, EndYear: endYear, RuntimeMinutes: runtimeMinutes, Genres: genres}
}
