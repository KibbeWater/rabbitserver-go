package rabbit

import (
	"math/rand"
	"strconv"
)

func calculateChecksum(imeiWithoutChecksum string) int {
	imeiArray := []int{}
	for _, v := range imeiWithoutChecksum {
		imeiArray = append(imeiArray, int(v))
	}

	sum := 0
	double := false
	for _, digit := range imeiArray {
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	checksum := (10 - (sum % 10)) % 10
	return checksum
}

func GenerateIMEI() string {
	TAC := "35847631"
	serialNumberPrefix := "00"
	serialNumber := serialNumberPrefix
	for i := 0; i < 4; i++ {
		serialNumber += string(rune(48 + rand.Intn(10)))
	}

	imeiWithoutChecksum := TAC + serialNumber
	checksum := calculateChecksum(imeiWithoutChecksum)
	generatedIMEI := imeiWithoutChecksum + strconv.Itoa(checksum)

	return generatedIMEI
}
