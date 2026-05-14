package ctx

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
)

const RequestID3 key2 = "RequestID"
type key2 string

// setter function 
func SetRequestID3(ctx context.Context, Ctxid string) context.Context{
	return context.WithValue(ctx, RequestID3, Ctxid)
}

// getter func 
func GetRequestIDv2(ctx context.Context) string {
	requestid, _ := ctx.Value(RequestID3). (string)
	return requestid
}


// generaing truly random keys 
func generateId() string{
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil{
		return "req-" + time.Now().Format("05.000")
	}
	return fmt.Sprintf("%x", b)
}