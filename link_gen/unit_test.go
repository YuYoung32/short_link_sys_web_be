/**
 * Created by YuYoung on 2023/4/17
 * Description: 短链生成函数测试
 */

package link_gen

import "testing"

func generateLongLink() string {
	var longLink string
	for i := 0; i < 100; i++ {
		longLink = longLink + "https://www.baidu.com"
	}
	return longLink
}

var testCases = []struct {
	longLink string
}{{
	longLink: "https://www.baidu.com",
}, {
	longLink: "1",
}, {
	longLink: "",
}, {
	longLink: generateLongLink(),
},
}

func TestMurmurHash(t *testing.T) {
	var murmur MurmurHash
	for _, testCase := range testCases {
		l1 := murmur.GenLink(testCase.longLink)
		l2 := murmur.GenLink(testCase.longLink)
		if l1 != l2 {
			t.Log("murmurhash generate link warning", l1, l2)
		} else {
			t.Log(l1)
		}
	}
}

func TestXXHash(t *testing.T) {
	var xx XXHash
	for _, testCase := range testCases {
		l1 := xx.GenLink(testCase.longLink)
		l2 := xx.GenLink(testCase.longLink)
		if l1 != l2 {
			t.Log("xxhash generate link warning", l1, l2)
		} else {
			t.Log(l1)
		}
	}
}

func TestFNVHash(t *testing.T) {
	var fnv FNVHash
	for _, testCase := range testCases {
		l1 := fnv.GenLink(testCase.longLink)
		l2 := fnv.GenLink(testCase.longLink)
		if l1 != l2 {
			t.Log("fnvhash generate link warning", l1, l2)
		} else {
			t.Log(l1)
		}
	}
}

func TestSimpleSequencer(t *testing.T) {
	var auto SimpleSequencer
	for _, testCase := range testCases {
		l1 := auto.GenLink(testCase.longLink)
		l2 := auto.GenLink(testCase.longLink)
		if l1 == l2 {
			t.Log("autoincrement generate link warning", l1, l2)
		} else {
			t.Log(l1, l2)
		}
	}
}

func TestSnowflakeSequencer(t *testing.T) {
	SnowflakeInit()
	var snow SnowflakeSequencer
	for _, testCase := range testCases {
		l1 := snow.GenLink(testCase.longLink)
		l2 := snow.GenLink(testCase.longLink)
		if l1 == l2 {
			t.Log("snowflake generate link warning", l1, l2)
		} else {
			t.Log(l1, l2)
		}
	}
}
