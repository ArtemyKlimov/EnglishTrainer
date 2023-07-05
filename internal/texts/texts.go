package texts

const (
	StartText string = "Hi, I'm your English trainer. You can use me to repeat vocabulary\n" +
		"Press /help to check what I can do"
	HelpText string = "You can control me by sending these commands: \n" +
		"/add - create new pair of phrases\n" +
		"/learn - learn some new words and phrases\n" +
		"/train - start practice your vocabulary\n" +
		"/exit - exit from current operation to the main menu"
	AddOperationText string = "Введите фразы в формате значение,синоним:value,synonym \n" +
		"Синонимов может быть несколько. В таком случае они разделяются запятыми"
	LearnOperationText string = "Let's learn some new phrases: "
	UnknownCommandText string = "Unknown command. Press /help for information"
	UnknownInputText   string = "Unknown input. Press /help for information"
	TrainInputText     string = "Let's practice some phrases:"
)
