package db

import "testing"

func TestGetGormDB(t *testing.T) {
	_, err := New("test.db")
	if err != nil {
		t.Fatal("open db failed")
	}
}
