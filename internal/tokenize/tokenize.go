package tokenize

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/kevinkiplangat432/arvis/internal/identifiers"
)

// tokenTTL is deliberately short — long enough to cover a real round
// trip to a provider and back, no longer. The token map is the most
// sensitive thing ARVIS ever holds, even briefly: it's the literal
// de-anonymization key. It never touches Postgres, and it doesn't
// outlive the request it belongs to.
const tokenTTL = 2 * time.Minute

func tokenMapKey(requestID string) string {
	return "tokenize:" + requestID
}

// Tokenize replaces every matched identifier in body with a
// sequential placeholder, stores the token→real-value map in Redis
// under requestID, and returns the tokenized body. Repeated values
// reuse the same token — two mentions of the same account number get
// the same placeholder — since the model may need to recognize
// they're the same entity to reason correctly about the prompt.
func Tokenize(ctx context.Context, rdb *redis.Client, requestID string, body []byte) ([]byte, error) {
	matches := identifiers.FindMatches(body)
	if len(matches) == 0 {
		return body, nil
	}

	tokenMap := make(map[string]string)
	valueToToken := make(map[string]string)
	counters := make(map[string]int)

	text := string(body)
	for _, m := range matches {
		token, seen := valueToToken[m.Value]
		if !seen {
			counters[m.Key]++
			token = fmt.Sprintf("[%s_%d]", strings.ToUpper(m.Key), counters[m.Key])
			valueToToken[m.Value] = token
			tokenMap[token] = m.Value
		}
		text = strings.ReplaceAll(text, m.Value, token)
	}

	data, err := json.Marshal(tokenMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token map: %w", err)
	}
	if err := rdb.Set(ctx, tokenMapKey(requestID), data, tokenTTL).Err(); err != nil {
		return nil, fmt.Errorf("failed to store token map: %w", err)
	}

	return []byte(text), nil
}

// Detokenize reverses Tokenize: loads the map, substitutes every
// token back to its real value, and deletes the map from Redis
// regardless of whether every token was actually found in the
// response — a model that never echoes a token back just means
// nothing gets substituted for it, not an error.
func Detokenize(ctx context.Context, rdb *redis.Client, requestID string, body []byte) ([]byte, error) {
	key := tokenMapKey(requestID)
	data, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return body, nil // nothing was tokenized for this request
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load token map: %w", err)
	}
	defer rdb.Del(ctx, key)

	var tokenMap map[string]string
	if err := json.Unmarshal(data, &tokenMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token map: %w", err)
	}

	text := string(body)
	for token, real := range tokenMap {
		text = strings.ReplaceAll(text, token, real)
	}
	return []byte(text), nil
}