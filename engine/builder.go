package engine

import (
	"context"
	"math/rand"
	"strings"

	"phatic_dialogue/database"
	"phatic_dialogue/types"
)

type Builder struct {
	singleInserts *database.SingleInserts
	groupInserts  *database.GroupInserts
	answers       *database.Answers
}

func NewBuilder(singleInserts *database.SingleInserts, groupInserts *database.GroupInserts, answers *database.Answers) *Builder {
	return &Builder{
		singleInserts: singleInserts,
		groupInserts:  groupInserts,
		answers:       answers,
	}
}

func (builder *Builder) MakeAnswer(ctx context.Context, topics []types.Topic) string {
	if len(topics) == 0 {
		return "..."
	}

	answer := builder.generateAnswer(ctx, getRandomElement(topics))

	return normaliseAnswer(answer)
}

type possibleElements interface {
	types.Topic | types.SingleInsert | types.GroupInsert | types.Answer
}

func getRandomElement[T possibleElements](elems []T) T {
	return elems[rand.Intn(len(elems))]
}

func (builder *Builder) generateAnswer(ctx context.Context, topic types.Topic) string {
	answers, err := builder.answers.List(ctx, topic)
	if err != nil {
		return "..."
	}

	if len(answers) == 0 {
		answers, err = builder.answers.List(ctx, types.UnknownTopic)
		if err != nil {
			return "..."
		}
	}

	return builder.insertWords(ctx, getRandomElement(answers))
}

func (builder *Builder) insertWords(ctx context.Context, answer types.Answer) string {
	for strings.Contains(answer.Answer, "$") { // $ - group insert / many words -> [а-я0-9 ]*.
		groupInserts, err := builder.groupInserts.List(ctx, answer.Topic)
		if err != nil {
			return "..."
		}

		answer.Answer = strings.Replace(answer.Answer, "$", getRandomElement(groupInserts).Words, 1)
	}

	for strings.Contains(answer.Answer, "_") { // _ - single insert / one word -> [а-я0-9]*.
		singleInserts, err := builder.singleInserts.List(ctx, answer.Topic)
		if err != nil {
			return "..."
		}

		answer.Answer = strings.Replace(answer.Answer, "_", getRandomElement(singleInserts).Word, 1)
	}

	return answer.Answer
}

func normaliseAnswer(answer string) string {
	oldStrings := []string{" ,", " .", " !", " ?"}
	newStrings := []string{",", ".", "!", "?"}

	for i, oldString := range oldStrings {
		answer = strings.Replace(answer, oldString, newStrings[i], -1)
	}

	return answer
}
