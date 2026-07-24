package nutsandbolts

import (
	"fmt"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//-----------------------------------------------------------------------------
// Print localized message with up to 2 parameters to Stderr
//-----------------------------------------------------------------------------

func PrintLocalizedMessage(localizer *i18n.Localizer, MessageID string, value1 string, value2 string) {
	localizedMessage := localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: MessageID,
			TemplateData: map[string]string{
				"value1": value1,
				"value2": value2,
			},
		})

	fmt.Fprint(os.Stderr, localizedMessage)
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func PrintLocalizedListHeader(localizer *i18n.Localizer, MessageID string, listLength int) {
	localizedMessage := localizer.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: MessageID,
			TemplateData: map[string]int{
				"count": listLength,
			},
			PluralCount: listLength,
		})

	fmt.Fprint(os.Stderr, localizedMessage)
}
