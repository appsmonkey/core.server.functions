package helper

import (
	"math/rand"
)

// CognitoIDZeroValue is the zero value for a Cognito ID
const CognitoIDZeroValue = "none"

// IsCognitoIDEmpty indicates if the cognito ID is empty.
// returns `true` if empty
func IsCognitoIDEmpty(cid string) bool {
	if len(cid) == 0 || cid == CognitoIDZeroValue {
		return true
	}

	return false
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq generated random string with spec. length
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetLang translations
func GetLang(lang string) interface{} {

	BA := map[string]string{
		"Mali_park":                     "Mali park",
		"Jevrejsko_groblje":             "Jevrejsko groblje",
		"Veliki_park":                   "Veliki park",
		"Budakovici":                    "Budakovići",
		"Bare":                          "Bare",
		"Dzidzikovac":                   "Džidžikovac",
		"Bascarsija":                    "Baščarsija",
		"Mihrivode":                     "Mihrivode",
		"Breka":                         "Breka",
		"Seri_Hill":                     "Šeri Hill",
		"Olimpijsko_naselje":            "Olimpijsko naselje",
		"Alifakovac":                    "Alifakovac",
		"Gorica":                        "Gorica",
		"Podhrastovi":                   "Podhrastovi",
		"Pogledi":                       "Pogledi",
		"Ciglane":                       "Ciglane",
		"Naselje_Logavina":              "Naselje Logavina",
		"Mejtas":                        "Mejtaš",
		"Vojnicko_polje":                "Vojnicko polje",
		"Vrbanjusa":                     "Vrbanjusa",
		"Sip-Betanija":                  "Sip-Betanija",
		"Sumbulusa":                     "Sumbulusa",
		"Medrese":                       "Medrese",
		"Panjina_Kula":                  "Panjina Kula",
		"Kovaci":                        "Kovači",
		"Bare_groblje":                  "Bare groblje",
		"Crni_vrh":                      "Crni vrh",
		"Mahmutovac":                    "Mahmutovac",
		"Hrid":                          "Hrid",
		"Gradska_deponija":              "Gradska deponija",
		"Dolac_Malta ":                  "Dolac Malta ",
		"Kozarevici":                    "Kozarevici",
		"Bjelave":                       "Bjelave",
		"Kosevsko_brdo":                 "Kosevsko brdo",
		"Vratnik":                       "Vratnik",
		"Socijalno":                     "Socijalno",
		"Hrasno":                        "Hrasno",
		"Jarcedoli":                     "Jarcedoli",
		"Skenderija":                    "Skenderija",
		"Zmajevac":                      "Zmajevac",
		"Vreoca":                        "Vreoca",
		"Centar":                        "Centar",
		"Bistrik":                       "Bistrik",
		"Stupsko_Brdo":                  "Stupsko Brdo",
		"Ugljesici":                     "Ugljesici",
		"Vraca":                         "Vraca",
		"Orahov_Brijeg":                 "Orahov Brijeg",
		"Svrakino_Selo":                 "Svrakino Selo",
		"Prljevo_Brdo":                  "Prljevo Brdo",
		"Nedarici":                      "Nedarici",
		"Vitkovac":                      "Vitkovac",
		"Kovacici":                      "Kovacici",
		"Luzani":                        "Luzani",
		"Sedrenik":                      "Sedrenik",
		"Pavlovac":                      "Pavlovac",
		"Hrasno_Brdo":                   "Hrasno Brdo",
		"Aneks":                         "Aneks",
		"Otoka":                         "Otoka",
		"Sirokaca":                      "Sirokaca",
		"Grbavica ":                     "Grbavica ",
		"Soukbunar":                     "Soukbunar",
		"Cengic_Vila":                   "Cengic Vila",
		"Reljevo":                       "Reljevo",
		"Marijin_Dvor":                  "Marijin Dvor",
		"Brijesce":                      "Brijesce",
		"Brijesce_brdo":                 "Brijesce brdo",
		"Sokolovic_Kolonija":            "Sokolovic Kolonija",
		"Alipasino_Polje":               "Alipasino Polje",
		"Otes":                          "Otes",
		"Zabrde":                        "Zabrde",
		"Blagovac":                      "Blagovac",
		"Boljakov_potok":                "Boljakov potok",
		"Medunarodni_aerodrom_Sarajevo": "Medunarodni aerodrom Sarajevo",
		"Sokolje":                       "Sokolje",
		"Donja_Josanica":                "Donja Josanica",
		"Mojmilo":                       "Mojmilo",
		"Mladicko_polje":                "Mladicko polje",
		"Azici":                         "Azici",
		"Velesici":                      "Velesici",
		"Aerodromsko_naselje":           "Aerodromsko naselje",
		"Doglodi":                       "Doglodi",
		"Kula":                          "Kula",
		"Dobrinja":                      "Dobrinja",
		"Toplik":                        "Toplik",
		"Ivanici":                       "Ivanici",
		"Kobilja_glava":                 "Kobilja glava",
		"Miljevici":                     "Miljevici",
		"Poljine":                       "Poljine",
		"Buca_Potok":                    "Buca Potok",
		"Pofalici":                      "Pofalici",
		"Ugorsko":                       "Ugorsko",
		"Stup":                          "Stup",
		"Halilovici":                    "Halilovici",
		"Hotonj":                        "Hotonj",
		"Kosevo":                        "Kosevo",
		"Ilidza":                        "Ilidza",
		"Vogosca":                       "Vogosca",
		"Tilava":                        "Tilava",
		"Kasindo":                       "Kasindo",
		"Rajlovac":                      "Rajlovac",
		"Lukavica":                      "Lukavica",
		"Butmir":                        "Butmir",
		"Trebevic":                      "Trebevic",
	}

	switch lang {
	case "BA":
		return BA
	default:
		return BA
	}
}
