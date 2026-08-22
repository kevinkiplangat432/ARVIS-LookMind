package policy

// Topic is a blockable subject-matter category. Keyword-based
// matching for v1, same philosophy as detector.ContentRule — every
// match is explainable to an auditor as "this exact pattern fired,"
// which matters more for a compliance product than raw recall would.
// Source is cited on every topic so a policy decision can always be
// traced back to why it exists, not just that it exists.
type Topic struct {
	Key      string
	Name     string
	Source   string
	Keywords []string
}

// Seeded from two real regulatory sources, as of August 2026:
//
// EU AI Act Article 5 lists eight practices banned outright since
// Feb 2025 — used here as an "unacceptable risk" reference, since
// Kenya's own AI Bill (still in draft, at Senate stage) is explicitly
// modeled in part on it.
//
// Kenya's own binding law today is the Data Protection Act 2019 —
// its sensitive personal data categories are used directly. The ODPC's
// July 2026 draft Guidance Note on AI additionally flags things like
// recommendation engines and automated credit/loan decisions for
// extra scrutiny; included here as draft guidance, explicitly not
// yet binding law, and worth re-checking as the AI Bill progresses.
var topicCatalog = map[string]Topic{
	"social_scoring": {
		Key: "social_scoring", Name: "Social scoring",
		Source:   "EU AI Act Art. 5(1)(c)",
		Keywords: []string{"social score", "social credit", "citizen scoring"},
	},
	"biometric_categorization": {
		Key: "biometric_categorization", Name: "Biometric categorization of sensitive traits",
		Source:   "EU AI Act Art. 5(1)(g)",
		Keywords: []string{"infer race", "infer religion", "infer sexual orientation", "biometric categorization"},
	},
	"emotion_recognition_workplace": {
		Key: "emotion_recognition_workplace", Name: "Workplace/education emotion recognition",
		Source:   "EU AI Act Art. 5(1)(f)",
		Keywords: []string{"emotion recognition", "employee mood detection", "detect employee emotion"},
	},
	"vulnerable_exploitation": {
		Key: "vulnerable_exploitation", Name: "Exploiting vulnerable groups",
		Source:   "EU AI Act Art. 5(1)(b)",
		Keywords: []string{"exploit vulnerability", "target children", "target elderly", "target disabled users"},
	},
	"realtime_biometric_id": {
		Key: "realtime_biometric_id", Name: "Real-time public biometric identification",
		Source:   "EU AI Act Art. 5(1)(h)",
		Keywords: []string{"real-time facial recognition", "live biometric identification"},
	},
	"predictive_policing": {
		Key: "predictive_policing", Name: "Profiling-based predictive policing",
		Source:   "EU AI Act Art. 5(1)(d)",
		Keywords: []string{"predict criminal behavior", "risk of offending profile"},
	},
	"sensitive_personal_data_ke": {
		Key: "sensitive_personal_data_ke", Name: "Kenyan sensitive personal data categories",
		Source:   "Kenya Data Protection Act 2019, s.2",
		Keywords: []string{"ethnic origin", "religious belief", "political opinion", "genetic data", "health data", "sex life", "sexual orientation", "criminal record"},
	},
	"automated_credit_decision_ke": {
		Key: "automated_credit_decision_ke", Name: "Automated credit/loan decisions",
		Source:   "Kenya ODPC draft AI Guidance Note, July 2026 (not yet binding)",
		Keywords: []string{"automatic loan approval", "automated credit scoring decision", "algorithmic loan denial"},
	},
	"recommendation_engine_ke": {
		Key: "recommendation_engine_ke", Name: "High-risk recommendation engine deployment",
		Source:   "Kenya ODPC draft AI Guidance Note, July 2026 (not yet binding)",
		Keywords: []string{"recommendation engine deployment", "content recommendation model"},
	},
}

func ListTopics() []Topic {
	out := make([]Topic, 0, len(topicCatalog))
	for _, t := range topicCatalog {
		out = append(out, t)
	}
	return out
}

func GetTopic(key string) (Topic, bool) {
	t, ok := topicCatalog[key]
	return t, ok
}