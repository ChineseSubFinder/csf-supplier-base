package cloudflare_helper

import (
	"github.com/ChineseSubFinder/csf-supplier/internal/dao"
	"github.com/ChineseSubFinder/csf-supplier/internal/models"
	"sort"
	"testing"
)

func TestCloudFlareHelper_uploadFile(t *testing.T) {

	c := NewCloudFlareHelper()

	//err := c.DeleteAllFile()
	//if err != nil {
	//	t.Fatal(err)
	//}

	// tt0304141
	// tt0373889
	imdbID := "tt1228705"
	var movieSubs []models.SubtitleMovie
	dao.Get().Where("imdb_id = ?", imdbID).Find(&movieSubs)
	// 降序排列
	sort.Sort(sort.Reverse(models.OrderSubtitleMovie(movieSubs)))

	for _, movieSub := range movieSubs {
		println("movieSub.SubSha256:", movieSub.SubSha256)

		orgDLUrl, err := c.GenerateDownloadUrl(&movieSub.SubtitleInfo)
		if err != nil {
			t.Fatal(err)
		}
		println(orgDLUrl)

		err = c.UploadFile(&movieSub.SubtitleInfo)
		if err != nil {
			t.Fatal(err)
		}
		println("UploadFile success")
	}

}
