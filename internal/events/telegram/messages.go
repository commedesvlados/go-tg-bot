package telegram

const msgHelp = `I can save and keep you pages. Also I can offer you then to read.

In order to save the page, just send me all link to it.

In order to get a random page from your list, send me command /rnd. Caution! After that, this page will be removed from your list`

const msgHello = "Hi there! ;)) \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command -_-"
	msgNoSavedPages   = "You have no saved pages *_*"
	msgSaved          = "Saved! >_<"
	msgAlreadyExists  = "You have already this page in your list (-=-)"
)
