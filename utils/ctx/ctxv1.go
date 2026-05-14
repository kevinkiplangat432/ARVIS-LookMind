package ctx

import "context"

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