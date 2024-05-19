package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
	"time"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i any) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes(i any, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i))
}

func Hash(i any) string {
	s := fmt.Sprint(i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Splitter(s string, sep string, i int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}

func Print(total, start int, fst, snd string) string {
	s := "| "
	s1 := fst + strings.Repeat(" ", start-len(fst)-2) + snd
	s += s1
	s += strings.Repeat(" ", total-len(s1)-2) + " |\n"
	return s
}

func ConvertTime(unixTime int64) string {
	// Unix 타임스탬프를 time.Time 객체로 변환
	t := time.Unix(unixTime, 0)

	// 년, 월, 일, 시간, 분을 추출
	year, month, day := t.Date()
	hour, minute, sec := t.Clock()

	return fmt.Sprintf("%d/%d/%d/%d:%d:%d", year, month, day, hour, minute, sec)
}
