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
			topic:     types.UnknownTopic,
			templates: []string{},
			answers: []string{"ваше питання таке цікаве, навіть не знаю, як на нього відповісти",
				"я це знаю , але сьогдні забув", "приходьте завтра з серйозними питаннями, сьогодні тільки кіно, іжа і музика",
				"мої розробники були в грайливому муді коли мене писали, тому я розмовляю лише на не серйозні теми",
				"сподіваюсь такі серйозні питання - це жарт", "в гуглі точно знають, мене туди поки не взяли",
				"я не в ресурсі сьогодні, приходьте завтра", "давай без отєтого от усього, будь ласка",
				"я сьогодні грайливий, з такими складними питаннями до chat gpt", "пикол зайшов занадто далеко"},
			//answers:       []string{"перефраазуйте $", "чи могли б ви переформулювати $"},
			singleInserts: []string{},
			groupInserts:  []string{"питання, будь ласка, може так ми зможемо знайти спільну мову"},
		},
		{
			topic:         "привітання",
			templates:     []string{"привіт", "вітаю"},
			answers:       []string{"привіт $як ", "вітаю $"},
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
			topic:         "смолток",
			templates:     []string{"як справи ?", "як день проходить ?", "як настрій ?"},
			answers:       []string{"_ $"},
			singleInserts: []string{"вдало", "класно", "прекрасно", "чудово", "цікаво"},
			groupInserts:  []string{", чим я можу бути корисний ?", ", чи є у вас якісь питання ?", ", що бажаєте дізнатись ?"},
		},
		{
			topic:         "вдячність",
			templates:     []string{"дякую", "дякую _", "дякую !", "дякую _ !", "дякую $", "дякую $ !"},
			answers:       []string{"будь ласка $", "$", "завжди радий допомогти $", "мені з вами теж було приємно працювати !", "мені подобається допомагати $"},
			singleInserts: []string{},
			groupInserts:  []string{", звертайтесь ще !"},
		},
		{
			topic:         "так",
			templates:     []string{"так", "погоджуюсь", "не  можу не погодитись", "є момент"},
			answers:       []string{"_ $"},
			singleInserts: []string{"файно", "прекрасно", "чудово", "приємно чути", "приємно знати"},
			groupInserts:  []string{", що ми знайшли з вами спільну мову", ", що ми це погодили", ", що ми це затвердили"},
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
			templates:     []string{"яка сьогодні _ погода ?", "який прогноз погоди $ ?", "яка погода буде _", "яка погода буде $", "чи буде $ дощ?", "чи буде _ дощ?", "чи треба мені _ брати парасольку ?"},
			answers:       []string{"я би радив вам переглянути прогноз погоди на _", "на _ ви можете це дізнатись", "ви можете дізнатись про це на _"},
			singleInserts: []string{"https://ua.sinoptik.ua/", "https://meteofor.com.ua/", "https://www.meteo.gov.ua/"},
			groupInserts:  []string{},
		},
		{
			topic:         "фільми",
			templates:     []string{"порадь фільми", "порадь фільми $", "напиши _ фільми $ ", "напиши _ фільми", "що мені подивитись $ ?", "що _ подивитись _ ?", "що _ подивитись $ ?", "що подивитись ?", "які  цікаві фільми $ ?", "що ти порадиш подивитись $ ?", "що ти порадиш подиивтись _ ?", "$ фільми $"},
			answers:       []string{"я би радив вам переглянути пропозиції на _", "на _ ви зможете собі щось підібрати", "ознайомтесь з підбіркою на _", "мені особисто подобаються : $", "я би вам порадив : $", "в трендах зараз : $", "я чув зараз модно диивтись : $"},
			singleInserts: []string{"https://megogo.net/ua/films", "https://sweet.tv/movie", "https://uakino.club/", "https://kinovezha.com/films/"},
			groupInserts:  []string{" 'Люксембург , Люксембург' , 'Довбуш' , 'Аватар' , 'Астероїд - Сіті' , 'Вавілон'", "'Месники' , 'Вартові галактики' , 'Чорна пантера' , 'Тор' , 'Людина Павук'", "'Три тисячі років нудьги' , 'Барбі', 'Першому гравцю приготуватися' , 'БлекБеррі'"},
		},
		{
			topic:         "книги1",
			templates:     []string{"що почитати $", "$ книжки $", "_ книжки $", "книжки $", "_ книгу $", "$ книгу $"},
			answers:       []string{"я би радив вам переглянути пропозиції на _", "на _ ви зможете собі щось підібрати", "ознайомтесь з підбіркою на _"},
			singleInserts: []string{"https://www.yakaboo.ua/", "https://book-ye.com.ua/", "https://vivat-book.com.ua/", "https://laboratoria.pro/"},
			groupInserts:  []string{},
		},
		{
			topic:         "книги2",
			templates:     []string{"порадь книгу", "що почитати ?", "що почитати _ ?", "які книжки зараз _ ?", "які книжки зараз  ?", "що зараз читають ?"},
			answers:       []string{"мені особисто подобаються : $", "я би вам порадив : $", "в трендах зараз : $", "я чув зараз модно читати : $"},
			singleInserts: []string{},
			groupInserts:  []string{" 'За перекопом є земля' , 'Наше спільне' , 'Дзвінка' , 'Ворошиловград' , 'Тигролови'", "'Кафе на краю світу' , 'Квіти для Елджерона' , 'Лбдина в пошуках справжнього сенсу' , 'Пляжне чтиво' , 'Драбина'"},
		},
		{
			topic:         "квитки",
			templates:     []string{"куди сходити $ ?", "куди сходити _ ?", "як провести вихідні ?", "що _ буде $ ?", "як провести вільний час ?", "що буде на $ ?", "що буде на _ ?", "що буде у $ ?", "що буде у _ ?"},
			answers:       []string{"ви можете ознайомитись з подіями на _", "переглянте пропозиції на _ ", "є кілька варіантів на _", "на _ ви зможете собі щось підібрати"},
			singleInserts: []string{"https://kontramarka.ua/uk/standUp", "https://molodyytheatre.com/", "http://ft.org.ua/ua/program", "http://newtheatre.kiev.ua/"},
			groupInserts:  []string{},
		},
		{
			topic:         "програмування",
			templates:     []string{"яку мову _ вивчити ?", "модна мова _", "на чому _ програмують ?", "яку мову _ обрати ?"},
			answers:       []string{"моїм розробникам подобається _ , $", "краще ніж _ ще нічого не придумали , $", "мені наспівала пташечка, що зараз модна _ , $"},
			singleInserts: []string{"goLang"},
			groupInserts:  []string{"ви можете дізнатись більше на https://go.dev/"},
		},
		{
			topic:         "рецепти",
			templates:     []string{"як приготувати _ ?", "як приготувати $ ?", "як готується _ ?", "як готується $ ?", "рецепт _", "рецепт $"},
			answers:       []string{"спробуйте відвідати _", "Клопотенко звичайно підозрілий тип, але спробуйте його рецепти https://klopotenko.com/reczepti/", "може спробуйте $ ", "особисто я спробував би $"},
			singleInserts: []string{"https://jisty.com.ua/category/howtocookthat/", "https://fayni-recepty.com.ua/"},
			groupInserts:  []string{"зварити ля пельмені"},
		},
		{
			topic:         "іжа",
			templates:     []string{"що приготувати $ ?", "чим здивувати $ ?", "чим здивувати _ ?"},
			answers:       []string{" спробуйте приготувтаи щось від Клопотенка $ ", "може спробуйте знайти щось на _ ", "можливо щось цікаве попадеться вам на _", "мені порадили подивтись на _", "приготуйте щось незвичайне $", "поексперементуйте на кухні $"},
			singleInserts: []string{"https://jisty.com.ua/category/howtocookthat/", "https://fayni-recepty.com.ua/"},
			groupInserts:  []string{"тут ви зможете дізнатись більше https://klopotenko.com/reczepti/"},
		},
		{
			topic:         "вільний час",
			templates:     []string{"що робити $ ?", "як провести вільний _ ?", "чим зайнятись у вільний _ ?", "чим зайнятись _ ?"},
			answers:       []string{"є пропозиція сходити $", "як варіант сходити $", "зараз в тренді сходити $", "пропоную вам піти $"},
			singleInserts: []string{""},
			groupInserts:  []string{"на тілесний перформанс", "на медитацію", "в спортзал", "в клуб", "в бар", "в бібліотеку", "в торгівельний центр", "за покупками", "прогулятись містом", "на виставку"},
		},
		{
			topic:         "музика",
			templates:     []string{"що мені послухати ?", "порекомендуй музику", "яка музика $ ?"},
			answers:       []string{"слухайте українське!"},
			singleInserts: []string{""},
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
