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
	var murmur murmurHash
	for _, testCase := range testCases {
		l1 := murmur.GenLink(testCase.longLink)
		l2 := murmur.GenLink(testCase.longLink)
		if l1 != l2 {
			t.Error("murmurhash generate link error", l1, l2)
		} else {
			t.Log(l1)
		}
	}
}

func TestXXHash(t *testing.T) {
	var xx xxHash
	for _, testCase := range testCases {
		l1 := xx.GenLink(testCase.longLink)
		l2 := xx.GenLink(testCase.longLink)
		if l1 != l2 {
			t.Error("xxhash generate link error", l1, l2)
		} else {
			t.Log(l1)
		}
	}
}

func TestFNVHash(t *testing.T) {
	var fnv fnvHash
	for _, testCase := range testCases {
		l1 := fnv.GenLink(testCase.longLink)
		l2 := fnv.GenLink(testCase.longLink)
		if l1 != l2 {
			t.Error("fnvhash generate link error", l1, l2)
		} else {
			t.Log(l1)
		}
	}
}

func TestAutoIncrement(t *testing.T) {
	var auto autoIncrement
	for _, testCase := range testCases {
		l1 := auto.GenLink(testCase.longLink)
		l2 := auto.GenLink(testCase.longLink)
		if l1 == l2 {
			t.Error("autoincrement generate link error", l1, l2)
		} else {
			t.Log(l1, l2)
		}
	}
}
