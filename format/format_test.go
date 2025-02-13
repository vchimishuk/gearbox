package format

import "testing"

func TestFileSize(t *testing.T) {
	assert(t, "0 B", FileSize(0))
	assert(t, "5 B", FileSize(5))
	assert(t, "1 kB", FileSize(1000))
	assert(t, "1 kB", FileSize(1023))
	assert(t, "1 kB", FileSize(1024))
	assert(t, "5 kB", FileSize(KB*5))
	assert(t, "9.8 kB", FileSize(KB*9+900))
	assert(t, "9.9 kB", FileSize(KB*9+1020))
	assert(t, "10 kB", FileSize(KB*10+512))
	assert(t, "15 kB", FileSize(KB*15+512))
	assert(t, "158 kB", FileSize(KB*158+512))
	assert(t, "1 MB", FileSize(KB*1000+512))
	assert(t, "1 MB", FileSize(KB*1023+512))
	assert(t, "1 MB", FileSize(MB*1))
	assert(t, "1 MB", FileSize(MB*1+KB))
	assert(t, "1.4 MB", FileSize(MB*1+KB*500))
	assert(t, "1.9 MB", FileSize(MB*1+KB*1023))
	assert(t, "9.9 MB", FileSize(MB*9+KB*1023))
	assert(t, "10 MB", FileSize(MB*10))
	assert(t, "10 MB", FileSize(MB*10+KB*512))
	assert(t, "124 MB", FileSize(MB*124))
	assert(t, "1 GB", FileSize(MB*1020))
	assert(t, "1 GB", FileSize(MB*1024))
	assert(t, "1 GB", FileSize(GB+MB*1))
	assert(t, "1.1 GB", FileSize(GB+MB*120))
	assert(t, "1 TB", FileSize(TB))
	assert(t, "1 TB", FileSize(TB+MB*512))
	assert(t, "1.5 TB", FileSize(TB+GB*512))
}

func assert(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Fatalf("%s != %s", expected, actual)
	}
}
