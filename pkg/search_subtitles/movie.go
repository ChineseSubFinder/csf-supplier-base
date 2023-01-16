package search_subtitles

import (
	"github.com/ChineseSubFinder/csf-supplier/internal/dao"
	"github.com/ChineseSubFinder/csf-supplier/internal/models"
	"sort"
)

func Movie(imdbID string) ([]models.SubtitleMovie, error) {

	var movieSubs []models.SubtitleMovie
	dao.Get().Where("imdb_id = ?", imdbID).Find(&movieSubs)
	// 降序排列
	sort.Sort(sort.Reverse(models.OrderSubtitleMovie(movieSubs)))

	return movieSubs, nil
}
