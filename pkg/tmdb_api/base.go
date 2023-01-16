package tmdb_api

type IdType int

const (
	Imdb IdType = iota + 1
	Tmdb
)

func (s IdType) String() string {
	switch s {
	case Imdb:
		return "imdb_id"
	case Tmdb:
		return "tmdb_id"
	default:
		return "Unknown"
	}
}

type ConvertIdResult struct {
	ImdbID string `json:"imdb_id"`
	TmdbID string `json:"tmdb_id"`
	TvdbID string `json:"tvdb_id"`
}

type TimeWindow int

const (
	Day TimeWindow = iota + 1
	Week
)

func (t TimeWindow) String() string {
	switch t {
	case Day:
		return "day"
	case Week:
		return "week"
	default:
		return "Unknown"
	}
}
