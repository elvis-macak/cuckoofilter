package cuckoofilter

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestInsertion(t *testing.T) {
	cf, capacity := NewCuckooFilter(1000000)
	if float32(1000000)/float32(capacity) < 0.95 {
		t.Errorf("Expected capacity usage must larger than 0.95")
	}
	fd, err := os.Open("/usr/share/dict/words")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)

	var values [][]byte
	var lineCount uint
	for scanner.Scan() {
		s := []byte(scanner.Text())
		if cf.InsertUnique(s) {
			lineCount++
		}
		values = append(values, s)
	}

	count := cf.Count()
	if count != lineCount {
		t.Errorf("Expected count = %d, instead count = %d", lineCount, count)
	}

	for _, v := range values {
		cf.Delete(v)
	}

	count = cf.Count()
	if count != 0 {
		t.Errorf("Expected count = 0, instead count == %d", count)
	}
}

func TestEncodeDecode(t *testing.T) {
	cf, capacity := NewCuckooFilter(8)
	if capacity != 8 {
		t.Errorf("Expected capacity is 8")
	}
	cf.buckets = []bucket{
		[bucketSize]FingerprintType{1, 2, 3, 4},
		[bucketSize]FingerprintType{5, 6, 7, 8},
	}
	cf.count = 8
	bytes := cf.Encode()
	ncf, err := Decode(bytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(cf, ncf) {
		t.Errorf("Expected %v, got %v", cf, ncf)
	}
}
