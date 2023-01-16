package imdb_info_center

import "testing"

func TestEpsId2TvId(t *testing.T) {

	// 黄石
	epsId := "tt4236770"
	// 黄石 S05E01
	epsId = "tt17663758"

	found, tvId, season, eps := EpsId2TvId(epsId)
	println(found, tvId, season, eps)
}
