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
			singleInserts: []string{"захопливо", "дивно", "раціонально"},
			groupInserts:  []string{"здивований більше вас !", "б і не міг собі уявити !"},
		},
		{
			topic:         "привітання",
			templates:     []string{"привіт", "вітаю"},
			answers:       []string{"привіт с ", "вітаю $"},
			singleInserts: []string{},
			groupInserts:  []string{" , чим я можу вам допомогти ?", " , чи є у вас якісь запитання ?"},
		},
		{
			topic:         "привітання ранок",
			templates:     []string{"доброго ранку", "добрий ранок"},
			answers:       []string{"доброго ранку і вам , сьогодні _ день $"},
			singleInserts: []string{"вдалий", "гарний", "класний", "прекрасний", "чудовий"},
			groupInserts:  []string{", чим я можу бути корисний ?", ", чи є у вас якісь питання ?", ", що бажаєте дізнатись ?"},
		},
		{
			topic:         "привітання день",
			templates:     []string{"добрий день", "доброго дня"},
			answers:       []string{"доброго дня , сподіваюсь ваш день проходить _ $", "добрий день , сподіваюсь ваш день проходить  _ $"},
			singleInserts: []string{"вдало", "класно", "прекрасно", "чудово", "цікаво"},
			groupInserts:  []string{", чим я можу бути корисний ?", ", чи є у вас якісь питання ?", ", що бажаєте дізнатись ?"},
		},
		{
			topic:         "привітання вечір",
			templates:     []string{"добрий вечір", "доброго вечора"},
			answers:       []string{"доброго вечора , сподіваюсь ваш день пройшов _ $", "добрий вечір , сподіваюсь ваш день пройшов  _ $"},
			singleInserts: []string{"вдало", "класно", "прекрасно", "чудово", "цікаво"},
			groupInserts:  []string{", чим я можу бути корисний ?", ", чи є у вас якісь питання ?", ", що бажаєте дізнатись ?"},
		},
		{
			topic:         "погода твердження",
			templates:     []string{"яка сьогодні _ погода", "сьогодні на вулиці так _", "завтра пронозують _ погоду"},
			answers:       []string{"не можу з вами не погодитись", "так , прогноз _ говорить про те саме", "якщо вірити прогнозу _"},
			singleInserts: []string{"погоди"},
			groupInserts:  []string{},
		},
		{
			topic:         "погода питання",
			templates:     []string{"яка сьогодні _ погода ?", "який прогноз погоди на _ ?", "яка погода буде _", "яка погода буде у _", "чи буде _ дощ?", "чи буде у _ дощ?", "чи треба мені _ брати парасольку ?"},
			answers:       []string{"я би радив вам переглянути прогноз погоди на _", "на _ ви можете це дізнатись", "ви можете дізнатись про це на _"},
			singleInserts: []string{"https://ua.sinoptik.ua/", "https://meteofor.com.ua/", "https://www.meteo.gov.ua/"},
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
