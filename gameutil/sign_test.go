package gameutil

import "testing"

func TestHmacSHA256(t *testing.T) {
	stringA := "appid=wx1234567&offer_id=12345678&openid=odkx20ENSNa2w5y3g_qOkOvBNM1g&pf=android&ts=1507530737&zone_id=1"
	stringSignTemp := stringA + "&org_loc=/cgi-bin/midas/getbalance&method=POST&secret=zNLgAGgqsEWJOg1nFVaO5r7fAlIQxr1u"

	sig := HmacSHA256(stringSignTemp, "zNLgAGgqsEWJOg1nFVaO5r7fAlIQxr1u")
	if sig != "1ad64e8dcb2ec1dc486b7fdf01f4a15159fc623dc3422470e51cf6870734726b" {
		t.Fail()
	}
}
