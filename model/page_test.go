package model

import "testing"

func TestTruncatedText(t *testing.T) {
	str := `一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十`
	tmp := TruncatedText(str)

	if len([]rune(tmp))-len(`...`) != TRUNCNUM {
		t.Errorf("trunc text error")
	}
}
