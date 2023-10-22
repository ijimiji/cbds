package crypt

import (
	"bytes"
	"strings"
)

var letterFrequencies = map[byte]float64{
	'a': 0.08167, 'b': 0.01492, 'c': 0.02782, 'd': 0.04253, 'e': 0.12702, 'f': 0.0228,
	'g': 0.02015, 'h': 0.06094, 'i': 0.06966, 'j': 0.00153, 'k': 0.00772, 'l': 0.04025,
	'm': 0.02406, 'n': 0.06749, 'o': 0.07507, 'p': 0.01929, 'q': 0.00095, 'r': 0.05987,
	's': 0.06327, 't': 0.09056, 'u': 0.02758, 'v': 0.00978, 'w': 0.0236, 'x': 0.0015,
	'y': 0.01974, 'z': 0.00074,
}

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz")

func Matrix(n int, keyWord string) [][]byte {
	m := len(alphabet)
	matrix := make([][]byte, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]byte, m)
	}

	for i := 1; i < n; i++ {
		matrix[i][0] = keyWord[i-1]
	}

	for j := 0; j < m; j++ {
		matrix[0][j] = alphabet[j]
	}

	for i := 1; i < n; i++ {
		startByte := matrix[i][0]
		startBytePos := bytes.IndexByte(alphabet, startByte)

		for j := 1; j < m-startBytePos; j++ {
			matrix[i][j] = alphabet[startBytePos+j]
		}

		l := 0
		for j := m - startBytePos; j < m; j++ {
			matrix[i][j] = alphabet[l]
			l++
		}
	}

	return matrix
}

func KeyWordLine(seq string, keyWord string) string {
	n := len(seq)
	m := len(keyWord)
	if n%m == 0 {
		return strings.Repeat(keyWord, n/m)
	}
	return strings.Repeat(keyWord, n/m) + keyWord[:n%m]
}

func EncodeLine(line string, keyWordLine string, matrix [][]byte, keyWord string) string {
	encodedLine := ""

	for i := range line {
		current := line[i]
		shiftByte := keyWordLine[i]

		if !bytes.Contains(alphabet, []byte{current}) {
			encodedLine += string(current)
			continue
		}

		currentPos := bytes.IndexByte([]byte(alphabet), current)
		shiftBytePos := bytes.IndexByte([]byte(keyWord), shiftByte) + 1
		encodedLine += string(matrix[shiftBytePos][currentPos])
	}

	return encodedLine
}

func Deltas(text string, l_gramm_length int) {

	dictionary := map[string]int{}
	entries := map[string][]int{}
	deltas := map[string][]int{}
	global_deltas := []byte{}

	for i := 0; i < (len(text) - l_gramm_length + 1); i++ {
		curr_l_gramm := text[i : i+l_gramm_length+1]
		if _, ok := dictionary[curr_l_gramm]; !ok {
			dictionary[curr_l_gramm] = 1
			entries[curr_l_gramm] = []int{i}
			continue
		}

		dictionary[curr_l_gramm] += 1
		deltas[curr_l_gramm] = append(deltas[curr_l_gramm], i-entries[curr_l_gramm][-1])
		global_deltas = append(global_deltas, i-entries[curr_l_gramm][-1])
		entries[curr_l_gramm] = append(entries[curr_l_gramm], i)
	}
	return deltas
}
