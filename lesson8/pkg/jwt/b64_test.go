package jwt

// import (
// 	"encoding/base64"
// 	"testing"
// )

var (
	in []byte = []byte("aaaaaaaaaaaaaaaaaa")
)

// func BenchmarkB64Str(b *testing.B) {

// 	for i := 0; i <= b.N; i++ {
// 		enc := base64.RawURLEncoding
// 		enc.EncodeToString(in)
// 	}
// }

// func BenchmarkB64Bytes(b *testing.B) {

// 	for i := 0; i <= b.N; i++ {
// 		enc := base64.RawURLEncoding
// 		res := make([]byte, enc.EncodedLen(len(in)))
// 		enc.Encode(res, in)
// 	}
// }
