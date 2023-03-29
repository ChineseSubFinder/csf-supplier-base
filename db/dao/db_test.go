package dao

import "testing"

func Test_initDB(t *testing.T) {

	err := initDB()
	if err != nil {
		t.Error(err)
	}
}
