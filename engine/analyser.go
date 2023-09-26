package engine

import (
	"context"
	"regexp"
	"strings"

	"phatic_dialogue/database"
	"phatic_dialogue/types"
)

type Analyser struct {
	templates *database.Templates
}

func NewAnalyser(templates *database.Templates) *Analyser {
	return &Analyser{templates: templates}
}

func (analyser *Analyser) AnalyseTopics(ctx context.Context, inStr string) []types.Topic {
	templates, err := analyser.templates.List(ctx)
	if err != nil {
		return []types.Topic{types.UnknownTopic}
	}

	topics := filterTopics(templates, normalizeSentence(inStr))
	if len(topics) == 0 {
		return []types.Topic{types.UnknownTopic}
	}

	return topics
}

func normalizeSentence(inStr string) string {
	inStr = strings.ToLower(inStr)
	var outStr string
	for i, symb := range inStr {
		if (symb == '.' || symb == ',' || symb == '!' || symb == '?') && i >= 1 && inStr[i-1] != ' ' {
			outStr += " " + string(symb)
			continue
		}
		outStr += string(symb)
	}

	return outStr
}

func filterTopics(templates []types.Template, normalisedSentence string) []types.Topic {
	topics := make([]types.Topic, 0)
	for _, template := range templates {
		template.Template = strings.Replace(template.Template, "_", "[а-я0-9]*", -1)  // _ - single insert / one word -> [а-я0-9]*.
		template.Template = strings.Replace(template.Template, "$", "[а-я0-9 ]*", -1) // $ - group insert / many words -> [а-я0-9 ]*.
		templateRegEx, err := regexp.Compile(template.Template)
		if err != nil {
			continue
		}

		if templateRegEx.MatchString(normalisedSentence) {
			if len(topics) >= 1 && topics[len(topics)-1] == template.Topic {
				continue
			}
			topics = append(topics, template.Topic)
		}
	}

	return topics
}
