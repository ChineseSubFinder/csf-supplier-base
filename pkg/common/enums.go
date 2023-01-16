package common

import "github.com/WQGroup/logger"

type MediaIdType int

const (
	IMDB MediaIdType = iota + 1
	TMDB
	TVDB
	DouBan
)

func (s MediaIdType) QueryString() string {
	switch s {
	case IMDB:
		return "imdb_id = ?"
	case TMDB:
		return "tmdb_id = ?"
	case TVDB:
		return "tvdb_id = ?"
	case DouBan:
		return "douban_id = ?"
	default:
		logger.Panicln("unknown media id type")
		return "unknown media id type"
	}
}

type MediaType int

const (
	Movie MediaType = iota + 1
	TV
)

func (m MediaType) String() string {
	switch m {
	case Movie:
		return "movie"
	case TV:
		return "tv"
	default:
		return "unknown"
	}
}
