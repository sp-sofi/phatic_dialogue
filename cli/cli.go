package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"phatic_dialogue/engine"
)

type CLI struct {
	analyser *engine.Analyser
	builder  *engine.Builder
}

func NewCLI(analyser *engine.Analyser, builder *engine.Builder) *CLI {
	return &CLI{
		analyser: analyser,
		builder:  builder,
	}
}

func (cli *CLI) Run(ctx context.Context) error {
	fmt.Println("WELCOME TO PHATIC-DIALOGUE PROGRAM")

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fmt.Print("you>> ")
		in := bufio.NewReader(os.Stdin)
		sentence, err := in.ReadString('\n')
		if err != nil {
			return err
		}

		sentence = sentence[:len(sentence)-1]

		if sentence == "/q" || sentence == "\\q" || sentence == "quit" {
			fmt.Println("BYE-BYE")
			return nil
		}

		fmt.Println("ms.X>> ", cli.builder.MakeAnswer(ctx, cli.analyser.AnalyseTopics(ctx, sentence)))
	}
}
