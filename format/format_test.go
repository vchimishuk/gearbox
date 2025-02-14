package format

import "testing"

func TestSize(t *testing.T) {
	assert(t, "0 B", Size(0))
	assert(t, "5 B", Size(5))
	assert(t, "1 kB", Size(1000))
	assert(t, "1 kB", Size(1023))
	assert(t, "1 kB", Size(1024))
	assert(t, "5 kB", Size(KB*5))
	assert(t, "9.8 kB", Size(KB*9+900))
	assert(t, "9.9 kB", Size(KB*9+1020))
	assert(t, "10 kB", Size(KB*10+512))
	assert(t, "15 kB", Size(KB*15+512))
	assert(t, "158 kB", Size(KB*158+512))
	assert(t, "1 MB", Size(KB*1000+512))
	assert(t, "1 MB", Size(KB*1023+512))
	assert(t, "1 MB", Size(MB*1))
	assert(t, "1 MB", Size(MB*1+KB))
	assert(t, "1.4 MB", Size(MB*1+KB*500))
	assert(t, "1.9 MB", Size(MB*1+KB*1023))
	assert(t, "9.9 MB", Size(MB*9+KB*1023))
	assert(t, "10 MB", Size(MB*10))
	assert(t, "10 MB", Size(MB*10+KB*512))
	assert(t, "124 MB", Size(MB*124))
	assert(t, "1 GB", Size(MB*1020))
	assert(t, "1 GB", Size(MB*1024))
	assert(t, "1 GB", Size(GB+MB*1))
	assert(t, "1.1 GB", Size(GB+MB*120))
	assert(t, "1 TB", Size(TB))
	assert(t, "1 TB", Size(TB+MB*512))
	assert(t, "1.5 TB", Size(TB+GB*512))
}

func assert(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Fatalf("%s != %s", expected, actual)
	}
}
