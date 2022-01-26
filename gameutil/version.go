package gameutil

import (
	"strings"
)

//格式 apkVersion;resVersion 兼容旧的只有resVersion
func ParseVersion(allVersion string) (apkVersion string, resVersion string) {
	apkVersion = "2.0.2"
	resVersion = allVersion
	tmp := strings.Split(allVersion, ";")
	if len(tmp) == 2 {
		apkVersion = tmp[0]
		resVersion = tmp[1]
	}
	return
}

//版本号v1是否大于等于版本v2
// 1,如果v1为"",返回false 如果v2为"",返回true
func IsVersionGE(v1, v2 string) bool {
	return CompareStringVersion(v1, v2) >= 0
}

//版本号v1是否大于版本v2
// 1,如果v1为"",返回false 如果v2为"",返回true
func IsVersionGreater(v1, v2 string) bool {
	return CompareStringVersion(v1, v2) > 0
}

func IsVersionSmaller(v1, v2 string) bool {
	return CompareStringVersion(v1, v2) < 0
}

//版本号v1是否大于版本v2
// 1,如果v1为"",返回false 如果v2为"",返回true
func CompareStringVersion(v1, v2 string) int {
	if v2 == "" {
		return 1
	}
	if v1 == "" {
		return -1
	}
	s1 := strings.Split(v1, ".")
	s2 := strings.Split(v2, ".")

	for len(s1) < len(s2) {
		s1 = append(s1, "0")
	}
	for len(s1) > len(s2) {
		s2 = append(s2, "0")
	}

	ret := CompareVersion(s1, s2)
	return ret
}

func CompareVersion(verA, verB []string) int {
	for index, _ := range verA {
		if ret := compareLittleVer(verA[index], verB[index]); ret != 0 {
			return ret
		}
	}

	return 0
}

func compareLittleVer(verA, verB string) int {

	bytesA := []byte(verA)
	bytesB := []byte(verB)

	lenA := len(bytesA)
	lenB := len(bytesB)
	if lenA > lenB {
		return 1
	}

	if lenA < lenB {
		return -1
	}

	//如果长度相等则按byte位进行比较
	return compareByBytes(bytesA, bytesB)
}

// 按byte位进行比较小版本号
func compareByBytes(verA, verB []byte) int {

	for index, _ := range verA {
		if verA[index] > verB[index] {
			return 1
		}
		if verA[index] < verB[index] {
			return -1
		}
	}

	return 0
}
