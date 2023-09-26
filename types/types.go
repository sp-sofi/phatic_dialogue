package types

type (
	Topic string

	SingleInsert struct {
		Word  string
		Topic Topic
	}

	GroupInsert struct {
		Words string
		Topic Topic
	}

	Template struct {
		Template string
		Topic    Topic
	}

	Answer struct {
		Answer string
		Topic  Topic
	}
)

const UnknownTopic Topic = "unknown_topic"
