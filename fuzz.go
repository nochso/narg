package narg

// Fuzz is exposed only for internal use with go-fuzz: https://github.com/dvyukov/go-fuzz
//
// Corpus resides in test-fixtures/fuzz/corpus
//
//  go get github.com/dvyukov/go-fuzz/go-fuzz
//  go get github.com/dvyukov/go-fuzz/go-fuzz-build
//  go-fuzz-build github.com/nochso/narg
//  go-fuzz -bin narg-fuzz.zip -workdir test-fixtures/fuzz
func Fuzz(data []byte) int {
	_, err := Parse(string(data))
	if err != nil {
		return 0
	}
	return 1
}
