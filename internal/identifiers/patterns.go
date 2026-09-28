package identifiers

import (
	"regexp"
	"strings"
)

// contextWindow is how many characters before a candidate match are
// searched for a cue word. This is specifically the fix for the old
// over-matching bug: a bare 7-8 digit pattern matched order IDs,
// phone fragments, token counts — anything numeric of that length.
// Requiring a nearby cue like "ID No" or "national id" cuts that down
// sharply without needing a full NER model.
const contextWindow = 40

type Identifier struct {
	Key         string
	Name        string
	Pattern     *regexp.Regexp
	ContextCues []string // empty means the pattern's own shape is distinctive enough alone
}

var Catalog = []Identifier{
	{
    Key:         "kenyan_national_id", 
    Name:        "Kenyan National ID number",
    Pattern:     regexp.MustCompile(`\b\d{7,9}\b`),
    ContextCues: []string{"id no", "id number", "national id", "id:", "identification", "maisha", "maisha namba"},
	},
	{
    Key:         "kra_pin", 
    Name:        "Kenyan KRA PIN",
    Pattern:     regexp.MustCompile(`(?i)\b[AP]\d{9}[A-Z]\b`),
	},
	{
    Key:         "mpesa_transaction_code",
    Name:        "M-PESA transaction code",
    Pattern:     regexp.MustCompile(`(?i)\b[A-Z]{2}[A-Z0-9]{8}\b`),
	},
	{
    Key:         "kenyan_phone_number",
    Name:        "Kenyan phone number",
    Pattern:     regexp.MustCompile(`\b(?:\+?254|0)[17]\d{8}\b`),
    ContextCues: []string{"phone", "mobile", "tel", "contact", "call", "sms", "whatsapp", "no."},
	},
}

type Match struct {
	Key   string
	Name  string
	Value string
	Start int
	End   int
}

// FindMatches scans body once against the whole catalog. This is now
// the single source of truth for "what counts as sensitive" — both
// the detector's flagging and the tokenizer's redaction call this
// same function, so a pattern tightened here fixes both at once
// instead of two copies silently drifting apart.
func FindMatches(body []byte) []Match {
	text := string(body)
	lower := strings.ToLower(text)
	var matches []Match

	for _, id := range Catalog {
		// FindAllStringIndex operates efficiently on the text
		for _, loc := range id.Pattern.FindAllStringIndex(text, -1) {
			start, end := loc[0], loc[1]
			if len(id.ContextCues) > 0 && !hasNearbyCue(lower, start, id.ContextCues) {
				continue
			}
			matches = append(matches, Match{
				Key:   id.Key,
				Name:  id.Name,
				Value: text[start:end],
				Start: start,
				End:   end,
			})
		}
	}
	return matches
}

func hasNearbyCue(lower string, matchStart int, cues []string) bool {
	windowStart := matchStart - contextWindow
	if windowStart < 0 {
		windowStart = 0
	}

	// Protect against splitting multi-byte UTF-8 characters
	for windowStart < matchStart && (lower[windowStart]&0xC0) == 0x80 {
		windowStart++
	}

	window := lower[windowStart:matchStart]
	for _, cue := range cues {
		if strings.Contains(window, cue) {
			return true
		}
	}
	return false
}