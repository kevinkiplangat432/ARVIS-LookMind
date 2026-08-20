package detector

import (
	"context"
	"regexp"
)

// ContentRule scans request and response bodies for sensitive Kenyan
// identifiers. Regex-based on purpose for v1 — fast, and every match
// is explainable to an auditor as "this exact pattern fired," which
// matters for a compliance product. A pluggable classifier is a
// natural later upgrade since it satisfies the same AsyncRule
// interface, not a rebuild.
//
// Known limitation worth being honest about: the national ID pattern
// (any bare 7-8 digit number) will over-match — order IDs, phone
// number fragments, token counts, anything numeric of that length.
// Fine as a first pass to prove the mechanism end to end; needs real
// tuning (context around the digits, checksum rules if any exist)
// before this is trustworthy enough to act on rather than just log.
type ContentRule struct {
	patterns map[string]*regexp.Regexp
}

func NewContentRule() *ContentRule {
	return &ContentRule{
		patterns: map[string]*regexp.Regexp{
			"kenyan_national_id":     regexp.MustCompile(`\b\d{7,8}\b`),
			"kra_pin":                regexp.MustCompile(`\b[A-Za-z]\d{9}[A-Za-z]\b`),
			"mpesa_transaction_code": regexp.MustCompile(`\b[A-Z]{2}[A-Z0-9]{8}\b`),
		},
	}
}

func (c *ContentRule) CheckAsync(ctx context.Context, in AsyncInput) []Flag {
	var flags []Flag

	for name, pattern := range c.patterns {
		if pattern.Match(in.RequestBody) {
			flags = append(flags, Flag{
				Rule: name, Category: "content", Severity: "high",
				Detail: "prompt appears to contain a " + humanName(name),
			})
		}
		if pattern.Match(in.ResponseBody) {
			flags = append(flags, Flag{
				Rule: name, Category: "content", Severity: "high",
				Detail: "response appears to contain a " + humanName(name),
			})
		}
	}

	return flags
}

func humanName(rule string) string {
	switch rule {
	case "kenyan_national_id":
		return "Kenyan national ID number"
	case "kra_pin":
		return "KRA PIN"
	case "mpesa_transaction_code":
		return "M-PESA transaction code"
	default:
		return rule
	}
}