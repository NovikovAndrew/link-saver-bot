package telegram

const msgHelp = `I can save your pages. Also i can offer you them to read.

In order to save your page, just send me link!

In order to get a random link from your list, send command /random
After that, this page will be removed from your list! 
`
const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command!"
	msgNoSavedPages   = "You have no saved paged"
	msgSaved          = "Saved!"
	msgAlreadyExist   = "You already have this page in your list"
)
