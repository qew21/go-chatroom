package util

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

var (
	currentSeqID = 0
	MAX_SEQ      = 0xFFF
	CHAR_CODE    = []byte{'2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
)

func GetSeq(typeID int, sessionID int) string {
	timeSeq := fmt.Sprintf("%042b", time.Now().UnixMilli())
	seqID := fmt.Sprintf("%012b", getMessageSeqID())
	binaryCode := timeSeq + seqID + fmt.Sprintf("%04b", typeID) + fmt.Sprintf("%022b", sessionID&0x3FFF)
	return getMessageID(binaryCode)
}

func getMessageSeqID() int {
	if currentSeqID > MAX_SEQ {
		currentSeqID = 0
	}
	ret := currentSeqID
	currentSeqID++
	return ret
}

func getMessageID(code string) string {
	ret := ""
	for i := 0; i < 16; i++ {
		asc, err := strconv.ParseInt(code[i*5:i*5+5], 2, 32)
		if err != nil {
			fmt.Println(err)
		}
		ret += string(CHAR_CODE[asc])
	}
	return ret
}

func getMessageCode(id string) string {
	ret := ""
	for i := 0; i < 16; i++ {
		code := bytes.IndexByte(CHAR_CODE, id[i])
		ret += fmt.Sprintf("%05b", code)
	}
	return ret
}
