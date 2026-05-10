package ctx
import (
	"context"
)

// assign a name and a specific type to a constant
const keyRequestID key = "requestID" // here key is a custom type

type key string // create a provate type called key 

// getter
func RequestID(ctx context.Context) string {
	// get the value from the context
	requestId, _ := ctx.Value(keyRequestID) .(string)
	return requestId
}

//setter
func SetRequestID(ctx context.Context, requestID string) context.Context{
	return  context.WithValue(ctx, keyRequestID, requestID)
}



// context in go are immutable you cannot change them 
// note : in go you capitalizee so that other function or so that the function itself can  call or be called.
// the setter pattern for coding it 
// signature:  func SetX(ctx, value)
// logic: context.WithValue(....)
// return: return the new ctx
// remember that the context.WithValue(///)..creates a copy.



// pattern to follow in the getter func 
// search --- convert-- return
// signature: func GetX(ctx)
// Logic ctx.Value(key)
// conversion .type