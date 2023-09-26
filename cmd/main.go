package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"phatic_dialogue/cli"
	"phatic_dialogue/database"
	"phatic_dialogue/engine"
	"phatic_dialogue/types"
)

// commands.
var (
	rootCmd = &cobra.Command{
		Use:   "",
		Short: "cli for interacting with program",
	}
	runCmd = &cobra.Command{
		Use:         "run",
		Short:       "runs the program",
		RunE:        cmdRun,
		Annotations: map[string]string{"type": "run"},
	}
	runSeed = &cobra.Command{
		Use:         "seed",
		Short:       "fills database",
		RunE:        cmdSeed,
		Annotations: map[string]string{"type": "seed"},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(runSeed)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func cmdRun(cmd *cobra.Command, args []string) error {
	dbURL := "postgres://postgres:123456@localhost:7766/phatic_dialogue?sslmode=disable"

	ctx, cancel := context.WithCancel(context.Background())
	onSigInt(func() {
		// starting graceful exit on context cancellation.
		cancel()
	})

	db, err := database.New(dbURL)
	if err != nil {
		return err
	}

	analyser := engine.NewAnalyser(db.Templates())
	builder := engine.NewBuilder(db.SingleInserts(), db.GroupInserts(), db.Answers())

	cli := cli.NewCLI(analyser, builder)

	return cli.Run(ctx)
}

func cmdSeed(cmd *cobra.Command, args []string) error {
	dbURL := "postgres://postgres:123456@localhost:7766/phatic_dialogue?sslmode=disable"

	ctx, cancel := context.WithCancel(context.Background())
	onSigInt(func() {
		// starting graceful exit on context cancellation.
		cancel()
	})

	db, err := database.New(dbURL)
	if err != nil {
		return err
	}

	err = db.CreateSchema(ctx)
	if err != nil {
		return err
	}

	data := []struct {
		topic         types.Topic
		templates     []string
		answers       []string
		singleInserts []string
		groupInserts  []string
	}{
		{
			topic:         types.UnknownTopic,
			templates:     []string{},
			answers:       []string{"перепрошую , я не зовсім вас розумію .", "дуже _ !", "хммм . . .", "я $"},
			singleInserts: []string{"захопливо", "цікаво", "дивно", "раціонально"},
			groupInserts:  []string{"здивований більеш вас !", "б і не міг собі уявити !"},
		},
		{
			topic:         "привітання",
			templates:     []string{"привіт", "вітаю", "доброго ранку", "добрий _"},
			answers:       []string{"привіт", "вітаю", "добрий _"},
			singleInserts: []string{"ранок", "день", "вечір"},
			groupInserts:  []string{},
		},
		{
			topic:         "погода1",
			templates:     []string{"яка сьогодні _ погода"},
			answers:       []string{"так , сьгодні дуже _ погода !"},
			singleInserts: []string{"дивна", "тепла", "холодна"},
			groupInserts:  []string{},
		},
	}

	for _, datum := range data {
		// create topic.
		err = db.Topics().Create(ctx, datum.topic)
		if err != nil {
			return err
		}

		// create templates.
		for _, template := range datum.templates {
			err = db.Templates().Create(ctx, types.Template{Template: template, Topic: datum.topic})
			if err != nil {
				return err
			}
		}

		// create answers.
		for _, answer := range datum.answers {
			err = db.Answers().Create(ctx, types.Answer{Answer: answer, Topic: datum.topic})
			if err != nil {
				return err
			}
		}

		// create singleInserts.
		for _, singleInsert := range datum.singleInserts {
			err = db.SingleInserts().Create(ctx, types.SingleInsert{Word: singleInsert, Topic: datum.topic})
			if err != nil {
				return err
			}
		}

		// create groupInserts.
		for _, groupInsert := range datum.groupInserts {
			err = db.GroupInserts().Create(ctx, types.GroupInsert{Words: groupInsert, Topic: datum.topic})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// onSigInt fires in SIGINT or SIGTERM event (usually CTRL+C).
func onSigInt(onSigInt func()) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		onSigInt()
	}()
}
