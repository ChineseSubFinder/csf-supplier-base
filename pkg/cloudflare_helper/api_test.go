package cloudflare_helper

import (
	"github.com/ChineseSubFinder/csf-supplier-base/db/dao"
	"github.com/ChineseSubFinder/csf-supplier-base/db/models"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"testing"
)

func TestCloudFlareHelper_uploadFile(t *testing.T) {

	c := NewCloudFlareHelper(settings.Get().CloudFlareConfig)

	//err := c.DeleteAllFile()
	//if err != nil {
	//	t.Fatal(err)
	//}

	// tt0304141
	// tt0373889
	//imdbID := "tt1228705"
	//var movieSubs []models.SubtitleMovie
	//dao.Get().Where("imdb_id = ?", imdbID).Find(&movieSubs)
	//// 降序排列
	//sort.Sort(sort.Reverse(models.OrderSubtitleMovie(movieSubs)))

	//var subs []models.SubtitleMovie
	//dao.Get().Where("imdb_id = ? AND sub_sha256 = ? AND language = ?",
	//	"tt1228705",
	//	"64360b7bee815e207d316e42e028aa45948ad0dab58b2cfecc287abd67bd680f",
	//	1,
	//).Find(&subs)

	var subs []models.SubtitleTV
	dao.Get().Where("full_season_sha256 = ? AND season = ? AND episode = ? AND imdb_id = ? AND sub_sha256 = ? AND language = ?",
		"819151d881c4110e545f72925e3962a67f9c38cc372f52b20f5e3e3ee291fe49",
		1, 1,
		"tt0903747",
		"17f9aab16df39cd166ee21a51ea803d7473db22cf5e0c599043b3f7177db616a",
		3,
	).Find(&subs)

	for _, movieSub := range subs {
		println("movieSub.SubSha256:", movieSub.SubSha256)

		orgDLUrl, err := c.GenerateDownloadUrl(&movieSub.SubtitleInfo)
		if err != nil {
			t.Fatal(err)
		}
		println(orgDLUrl)

		//err = c.UploadFile(settings.Get().HouseKeepingConfig, &movieSub.SubtitleInfo)
		//if err != nil {
		//	t.Fatal(err)
		//}
		println("UploadFile success")
	}

}
