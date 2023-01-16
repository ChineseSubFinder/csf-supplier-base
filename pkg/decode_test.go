package pkg

import "testing"

func TestGetIMDBIdFromIMDBUrl(t *testing.T) {

	url := "http://www.imdb.com/title/tt0583435/"
	if GetIMDBIdFromIMDBUrl(url) != "tt0583435" {
		t.Fatal("GetIMDBIdFromIMDBUrl error")
	}
}

func TestGetDouBanIdFromDouBanUrl(t *testing.T) {

	url := "https://movie.douban.com/subject/27605669/"
	if GetDouBanIdFromDouBanUrl(url) != "27605669" {
		t.Fatal("GetDouBanIdFromDouBanUrl error")
	}
}
