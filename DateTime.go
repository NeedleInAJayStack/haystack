package haystack

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// DateTime models a timestamp with a specific timezone.
type DateTime struct {
	time time.Time
}

// NewDateTime creates a new DateTime object. The values are not validated for correctness.
func NewDateTime(date Date, htime Time, tzOffset int, tz string) DateTime {
	loc, _ := time.LoadLocation(tz)
	goTime := time.Date(
		date.year,
		time.Month(date.month),
		date.day,
		htime.hour,
		htime.min,
		htime.sec,
		htime.ms*1000,
		loc,
	)
	return DateTime{
		time: goTime,
	}
}

// NewDateTimeRaw creates a new DateTime object. The values are not validated for correctness.
func NewDateTimeRaw(year int, month int, day int, hour int, min int, sec int, ms int, tz string) (DateTime, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return DateTime{}, err
	}
	goTime := time.Date(
		year,
		time.Month(month),
		day,
		hour,
		min,
		sec,
		ms*1000,
		loc,
	)
	return DateTime{time: goTime}, nil
}

// NewDateTimeFromString creates a DateTime object from a string in the format: "YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"
func NewDateTimeFromString(str string) (DateTime, error) {
	split := strings.Split(str, " ")
	goTime, err := time.Parse(time.RFC3339Nano, split[0])
	if len(split) > 1 {
		shortName := split[1]
		longName := tzShortNameMap[shortName]
		location, _ := time.LoadLocation(longName)
		if location != nil {
			goTime = time.Date(
				goTime.Year(),
				goTime.Month(),
				goTime.Day(),
				goTime.Hour(),
				goTime.Minute(),
				goTime.Second(),
				goTime.Nanosecond(),
				location,
			)
		}
	}
	return DateTime{time: goTime}, err
}

func NewDateTimeFromGo(goTime time.Time) DateTime {
	return DateTime{
		time: goTime,
	}
}

func dateTimeDef() DateTime {
	return NewDateTime(Date{}, Time{}, 0, "UTC")
}

// Date returns the date of the object.
func (dateTime DateTime) Date() Date {
	return NewDate(
		dateTime.time.Year(),
		int(dateTime.time.Month()),
		dateTime.time.Day(),
	)
}

// Time returns the date of the object.
func (dateTime DateTime) Time() Time {
	return NewTime(
		dateTime.time.Hour(),
		dateTime.time.Minute(),
		dateTime.time.Second(),
		dateTime.time.Nanosecond()/1000,
	)
}

// Tz returns the timezone of the object.
func (dateTime DateTime) Tz() string {
	longName := dateTime.time.Location().String()
	split := strings.Split(longName, "/")
	return split[len(split)-1]
}

// ToTz returns a new DateTime adjusted to the requested timezone. The timezone must be in shortened name format
func (dateTime DateTime) ToTz(shortName string) (DateTime, error) {
	longName := tzShortNameMap[shortName]
	newLocation, err := time.LoadLocation(longName)
	if err != nil {
		return DateTime{}, err
	}
	return DateTime{time: dateTime.time.In(newLocation)}, nil
}

// ToZinc represents the object as: "YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"
func (dateTime DateTime) ToZinc() string {
	buf := strings.Builder{}
	dateTime.encodeTo(&buf, true)
	return buf.String()
}

// ToAxon represents the object as: "dateTime(YYYY-MM-DD, hh:mm:ss.FFF)"
func (dateTime DateTime) ToAxon() string {
	return "dateTime(" + dateTime.Date().ToZinc() + "," + dateTime.Time().ToZinc() + ",\"" + dateTime.Tz() + "\")"
}

// ToGo creates a Go time.Time representation of the object
func (dateTime DateTime) ToGo() time.Time {
	return dateTime.time
}

// MarshalJSON representes the object as: "t:YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"
func (dateTime DateTime) MarshalJSON() ([]byte, error) {
	buf := strings.Builder{}
	buf.WriteString("t:")
	dateTime.encodeTo(&buf, true)
	return json.Marshal(buf.String())
}

// UnmarshalJSON interprets the json value: "t:YYYY-MM-DD'T'hh:mm:ss.FFFz zzzz"
func (dateTime *DateTime) UnmarshalJSON(buf []byte) error {
	var jsonStr string
	err := json.Unmarshal(buf, &jsonStr)
	if err != nil {
		return err
	}

	newDateTime, newErr := dateTimeFromJSON(jsonStr)
	*dateTime = newDateTime
	return newErr
}

func dateTimeFromJSON(jsonStr string) (DateTime, error) {
	if !strings.HasPrefix(jsonStr, "t:") {
		return dateTimeDef(), errors.New("value does not begin with 't:'")
	}
	dateTimeStr := jsonStr[2:]

	parseDateTime, parseErr := NewDateTimeFromString(dateTimeStr)
	if parseErr != nil {
		return dateTimeDef(), parseErr
	}
	return parseDateTime, nil
}

// MarshalHayson representes the object as: "{\"_kind\":\"dateTime\",\"val\":\"YYYY-MM-DD'T'hh:mm:ss.FFFz\",\"tz\":\"zzzz\"}"
func (dateTime DateTime) MarshalHayson() ([]byte, error) {
	buf := strings.Builder{}
	buf.WriteString("{\"_kind\":\"dateTime\",\"val\":\"")
	dateTime.encodeTo(&buf, false)
	buf.WriteString("\",\"tz\":\"")
	buf.WriteString(dateTime.Tz())
	buf.WriteString("\"}")
	return []byte(buf.String()), nil
}

func (dateTime DateTime) encodeTo(buf *strings.Builder, includeTz bool) {
	buf.WriteString(dateTime.time.Format(time.RFC3339Nano))
	if includeTz {
		buf.WriteRune(' ')
		buf.WriteString(dateTime.Tz())
	}
}

// Maps from a short timezone name to the long version. The reverse mapping is to take the last element of a '/' split.
// This is important because Fantom uses short timezone names, whereas go requires full names
var tzShortNameMap = map[string]string{
	"Andorra":        "Europe/Andorra",
	"Dubai":          "Asia/Dubai",
	"Kabul":          "Asia/Kabul",
	"Antigua":        "America/Antigua",
	"Anguilla":       "America/Anguilla",
	"Tirane":         "Europe/Tirane",
	"Yerevan":        "Asia/Yerevan",
	"Luanda":         "Africa/Luanda",
	"McMurdo":        "Antarctica/McMurdo",
	"Casey":          "Antarctica/Casey",
	"Davis":          "Antarctica/Davis",
	"DumontDUrville": "Antarctica/DumontDUrville",
	"Mawson":         "Antarctica/Mawson",
	"Palmer":         "Antarctica/Palmer",
	"Rothera":        "Antarctica/Rothera",
	"Syowa":          "Antarctica/Syowa",
	"Troll":          "Antarctica/Troll",
	"Vostok":         "Antarctica/Vostok",
	"Buenos_Aires":   "America/Argentina/Buenos_Aires",
	"Cordoba":        "America/Argentina/Cordoba",
	"Salta":          "America/Argentina/Salta",
	"Jujuy":          "America/Argentina/Jujuy",
	"Tucuman":        "America/Argentina/Tucuman",
	"Catamarca":      "America/Argentina/Catamarca",
	"La_Rioja":       "America/Argentina/La_Rioja",
	"San_Juan":       "America/Argentina/San_Juan",
	"Mendoza":        "America/Argentina/Mendoza",
	"San_Luis":       "America/Argentina/San_Luis",
	"Rio_Gallegos":   "America/Argentina/Rio_Gallegos",
	"Ushuaia":        "America/Argentina/Ushuaia",
	"Pago_Pago":      "Pacific/Pago_Pago",
	"Vienna":         "Europe/Vienna",
	"Lord_Howe":      "Australia/Lord_Howe",
	"Macquarie":      "Antarctica/Macquarie",
	"Hobart":         "Australia/Hobart",
	"Melbourne":      "Australia/Melbourne",
	"Sydney":         "Australia/Sydney",
	"Broken_Hill":    "Australia/Broken_Hill",
	"Brisbane":       "Australia/Brisbane",
	"Lindeman":       "Australia/Lindeman",
	"Adelaide":       "Australia/Adelaide",
	"Darwin":         "Australia/Darwin",
	"Perth":          "Australia/Perth",
	"Eucla":          "Australia/Eucla",
	"Aruba":          "America/Aruba",
	"Mariehamn":      "Europe/Mariehamn",
	"Baku":           "Asia/Baku",
	"Sarajevo":       "Europe/Sarajevo",
	"Barbados":       "America/Barbados",
	"Dhaka":          "Asia/Dhaka",
	"Brussels":       "Europe/Brussels",
	"Ouagadougou":    "Africa/Ouagadougou",
	"Sofia":          "Europe/Sofia",
	"Bahrain":        "Asia/Bahrain",
	"Bujumbura":      "Africa/Bujumbura",
	"Porto-Novo":     "Africa/Porto-Novo",
	"St_Barthelemy":  "America/St_Barthelemy",
	"Bermuda":        "Atlantic/Bermuda",
	"Brunei":         "Asia/Brunei",
	"La_Paz":         "America/La_Paz",
	"Kralendijk":     "America/Kralendijk",
	"Noronha":        "America/Noronha",
	"Belem":          "America/Belem",
	"Fortaleza":      "America/Fortaleza",
	"Recife":         "America/Recife",
	"Araguaina":      "America/Araguaina",
	"Maceio":         "America/Maceio",
	"Bahia":          "America/Bahia",
	"Sao_Paulo":      "America/Sao_Paulo",
	"Campo_Grande":   "America/Campo_Grande",
	"Cuiaba":         "America/Cuiaba",
	"Santarem":       "America/Santarem",
	"Porto_Velho":    "America/Porto_Velho",
	"Boa_Vista":      "America/Boa_Vista",
	"Manaus":         "America/Manaus",
	"Eirunepe":       "America/Eirunepe",
	"Rio_Branco":     "America/Rio_Branco",
	"Nassau":         "America/Nassau",
	"Thimphu":        "Asia/Thimphu",
	"Gaborone":       "Africa/Gaborone",
	"Minsk":          "Europe/Minsk",
	"Belize":         "America/Belize",
	"St_Johns":       "America/St_Johns",
	"Halifax":        "America/Halifax",
	"Glace_Bay":      "America/Glace_Bay",
	"Moncton":        "America/Moncton",
	"Goose_Bay":      "America/Goose_Bay",
	"Blanc-Sablon":   "America/Blanc-Sablon",
	"Toronto":        "America/Toronto",
	"Iqaluit":        "America/Iqaluit",
	"Atikokan":       "America/Atikokan",
	"Winnipeg":       "America/Winnipeg",
	"Resolute":       "America/Resolute",
	"Rankin_Inlet":   "America/Rankin_Inlet",
	"Regina":         "America/Regina",
	"Swift_Current":  "America/Swift_Current",
	"Edmonton":       "America/Edmonton",
	"Cambridge_Bay":  "America/Cambridge_Bay",
	"Inuvik":         "America/Inuvik",
	"Creston":        "America/Creston",
	"Dawson_Creek":   "America/Dawson_Creek",
	"Fort_Nelson":    "America/Fort_Nelson",
	"Whitehorse":     "America/Whitehorse",
	"Dawson":         "America/Dawson",
	"Vancouver":      "America/Vancouver",
	"Cocos":          "Indian/Cocos",
	"Kinshasa":       "Africa/Kinshasa",
	"Lubumbashi":     "Africa/Lubumbashi",
	"Bangui":         "Africa/Bangui",
	"Brazzaville":    "Africa/Brazzaville",
	"Zurich":         "Europe/Zurich",
	"Abidjan":        "Africa/Abidjan",
	"Rarotonga":      "Pacific/Rarotonga",
	"Santiago":       "America/Santiago",
	"Punta_Arenas":   "America/Punta_Arenas",
	"Easter":         "Pacific/Easter",
	"Douala":         "Africa/Douala",
	"Shanghai":       "Asia/Shanghai",
	"Urumqi":         "Asia/Urumqi",
	"Bogota":         "America/Bogota",
	"Costa_Rica":     "America/Costa_Rica",
	"Havana":         "America/Havana",
	"Cape_Verde":     "Atlantic/Cape_Verde",
	"Curacao":        "America/Curacao",
	"Christmas":      "Indian/Christmas",
	"Nicosia":        "Asia/Nicosia",
	"Famagusta":      "Asia/Famagusta",
	"Prague":         "Europe/Prague",
	"Berlin":         "Europe/Berlin",
	"Busingen":       "Europe/Busingen",
	"Djibouti":       "Africa/Djibouti",
	"Copenhagen":     "Europe/Copenhagen",
	"Dominica":       "America/Dominica",
	"Santo_Domingo":  "America/Santo_Domingo",
	"Algiers":        "Africa/Algiers",
	"Guayaquil":      "America/Guayaquil",
	"Galapagos":      "Pacific/Galapagos",
	"Tallinn":        "Europe/Tallinn",
	"Cairo":          "Africa/Cairo",
	"El_Aaiun":       "Africa/El_Aaiun",
	"Asmara":         "Africa/Asmara",
	"Madrid":         "Europe/Madrid",
	"Ceuta":          "Africa/Ceuta",
	"Canary":         "Atlantic/Canary",
	"Addis_Ababa":    "Africa/Addis_Ababa",
	"Helsinki":       "Europe/Helsinki",
	"Fiji":           "Pacific/Fiji",
	"Stanley":        "Atlantic/Stanley",
	"Chuuk":          "Pacific/Chuuk",
	"Pohnpei":        "Pacific/Pohnpei",
	"Kosrae":         "Pacific/Kosrae",
	"Faroe":          "Atlantic/Faroe",
	"Paris":          "Europe/Paris",
	"Libreville":     "Africa/Libreville",
	"London":         "Europe/London",
	"Grenada":        "America/Grenada",
	"Tbilisi":        "Asia/Tbilisi",
	"Cayenne":        "America/Cayenne",
	"Guernsey":       "Europe/Guernsey",
	"Accra":          "Africa/Accra",
	"Gibraltar":      "Europe/Gibraltar",
	"Nuuk":           "America/Nuuk",
	"Danmarkshavn":   "America/Danmarkshavn",
	"Scoresbysund":   "America/Scoresbysund",
	"Thule":          "America/Thule",
	"Banjul":         "Africa/Banjul",
	"Conakry":        "Africa/Conakry",
	"Guadeloupe":     "America/Guadeloupe",
	"Malabo":         "Africa/Malabo",
	"Athens":         "Europe/Athens",
	"South_Georgia":  "Atlantic/South_Georgia",
	"Guatemala":      "America/Guatemala",
	"Guam":           "Pacific/Guam",
	"Bissau":         "Africa/Bissau",
	"Guyana":         "America/Guyana",
	"Hong_Kong":      "Asia/Hong_Kong",
	"Tegucigalpa":    "America/Tegucigalpa",
	"Zagreb":         "Europe/Zagreb",
	"Port-au-Prince": "America/Port-au-Prince",
	"Budapest":       "Europe/Budapest",
	"Jakarta":        "Asia/Jakarta",
	"Pontianak":      "Asia/Pontianak",
	"Makassar":       "Asia/Makassar",
	"Jayapura":       "Asia/Jayapura",
	"Dublin":         "Europe/Dublin",
	"Jerusalem":      "Asia/Jerusalem",
	"Isle_of_Man":    "Europe/Isle_of_Man",
	"Kolkata":        "Asia/Kolkata",
	"Chagos":         "Indian/Chagos",
	"Baghdad":        "Asia/Baghdad",
	"Tehran":         "Asia/Tehran",
	"Reykjavik":      "Atlantic/Reykjavik",
	"Rome":           "Europe/Rome",
	"Jersey":         "Europe/Jersey",
	"Jamaica":        "America/Jamaica",
	"Amman":          "Asia/Amman",
	"Tokyo":          "Asia/Tokyo",
	"Nairobi":        "Africa/Nairobi",
	"Bishkek":        "Asia/Bishkek",
	"Phnom_Penh":     "Asia/Phnom_Penh",
	"Tarawa":         "Pacific/Tarawa",
	"Kanton":         "Pacific/Kanton",
	"Kiritimati":     "Pacific/Kiritimati",
	"Comoro":         "Indian/Comoro",
	"St_Kitts":       "America/St_Kitts",
	"Pyongyang":      "Asia/Pyongyang",
	"Seoul":          "Asia/Seoul",
	"Kuwait":         "Asia/Kuwait",
	"Cayman":         "America/Cayman",
	"Almaty":         "Asia/Almaty",
	"Qyzylorda":      "Asia/Qyzylorda",
	"Qostanay":       "Asia/Qostanay",
	"Aqtobe":         "Asia/Aqtobe",
	"Aqtau":          "Asia/Aqtau",
	"Atyrau":         "Asia/Atyrau",
	"Oral":           "Asia/Oral",
	"Vientiane":      "Asia/Vientiane",
	"Beirut":         "Asia/Beirut",
	"St_Lucia":       "America/St_Lucia",
	"Vaduz":          "Europe/Vaduz",
	"Colombo":        "Asia/Colombo",
	"Monrovia":       "Africa/Monrovia",
	"Maseru":         "Africa/Maseru",
	"Vilnius":        "Europe/Vilnius",
	"Luxembourg":     "Europe/Luxembourg",
	"Riga":           "Europe/Riga",
	"Tripoli":        "Africa/Tripoli",
	"Casablanca":     "Africa/Casablanca",
	"Monaco":         "Europe/Monaco",
	"Chisinau":       "Europe/Chisinau",
	"Podgorica":      "Europe/Podgorica",
	"Marigot":        "America/Marigot",
	"Antananarivo":   "Indian/Antananarivo",
	"Majuro":         "Pacific/Majuro",
	"Kwajalein":      "Pacific/Kwajalein",
	"Skopje":         "Europe/Skopje",
	"Bamako":         "Africa/Bamako",
	"Yangon":         "Asia/Yangon",
	"Ulaanbaatar":    "Asia/Ulaanbaatar",
	"Hovd":           "Asia/Hovd",
	"Choibalsan":     "Asia/Choibalsan",
	"Macau":          "Asia/Macau",
	"Saipan":         "Pacific/Saipan",
	"Martinique":     "America/Martinique",
	"Nouakchott":     "Africa/Nouakchott",
	"Montserrat":     "America/Montserrat",
	"Malta":          "Europe/Malta",
	"Mauritius":      "Indian/Mauritius",
	"Maldives":       "Indian/Maldives",
	"Blantyre":       "Africa/Blantyre",
	"Mexico_City":    "America/Mexico_City",
	"Cancun":         "America/Cancun",
	"Merida":         "America/Merida",
	"Monterrey":      "America/Monterrey",
	"Matamoros":      "America/Matamoros",
	"Chihuahua":      "America/Chihuahua",
	"Ciudad_Juarez":  "America/Ciudad_Juarez",
	"Ojinaga":        "America/Ojinaga",
	"Mazatlan":       "America/Mazatlan",
	"Bahia_Banderas": "America/Bahia_Banderas",
	"Hermosillo":     "America/Hermosillo",
	"Tijuana":        "America/Tijuana",
	"Kuala_Lumpur":   "Asia/Kuala_Lumpur",
	"Kuching":        "Asia/Kuching",
	"Maputo":         "Africa/Maputo",
	"Windhoek":       "Africa/Windhoek",
	"Noumea":         "Pacific/Noumea",
	"Niamey":         "Africa/Niamey",
	"Norfolk":        "Pacific/Norfolk",
	"Lagos":          "Africa/Lagos",
	"Managua":        "America/Managua",
	"Amsterdam":      "Europe/Amsterdam",
	"Oslo":           "Europe/Oslo",
	"Kathmandu":      "Asia/Kathmandu",
	"Nauru":          "Pacific/Nauru",
	"Niue":           "Pacific/Niue",
	"Auckland":       "Pacific/Auckland",
	"Chatham":        "Pacific/Chatham",
	"Muscat":         "Asia/Muscat",
	"Panama":         "America/Panama",
	"Lima":           "America/Lima",
	"Tahiti":         "Pacific/Tahiti",
	"Marquesas":      "Pacific/Marquesas",
	"Gambier":        "Pacific/Gambier",
	"Port_Moresby":   "Pacific/Port_Moresby",
	"Bougainville":   "Pacific/Bougainville",
	"Manila":         "Asia/Manila",
	"Karachi":        "Asia/Karachi",
	"Warsaw":         "Europe/Warsaw",
	"Miquelon":       "America/Miquelon",
	"Pitcairn":       "Pacific/Pitcairn",
	"Puerto_Rico":    "America/Puerto_Rico",
	"Gaza":           "Asia/Gaza",
	"Hebron":         "Asia/Hebron",
	"Lisbon":         "Europe/Lisbon",
	"Madeira":        "Atlantic/Madeira",
	"Azores":         "Atlantic/Azores",
	"Palau":          "Pacific/Palau",
	"Asuncion":       "America/Asuncion",
	"Qatar":          "Asia/Qatar",
	"Reunion":        "Indian/Reunion",
	"Bucharest":      "Europe/Bucharest",
	"Belgrade":       "Europe/Belgrade",
	"Kaliningrad":    "Europe/Kaliningrad",
	"Moscow":         "Europe/Moscow",
	"Simferopol":     "Europe/Simferopol",
	"Kirov":          "Europe/Kirov",
	"Volgograd":      "Europe/Volgograd",
	"Astrakhan":      "Europe/Astrakhan",
	"Saratov":        "Europe/Saratov",
	"Ulyanovsk":      "Europe/Ulyanovsk",
	"Samara":         "Europe/Samara",
	"Yekaterinburg":  "Asia/Yekaterinburg",
	"Omsk":           "Asia/Omsk",
	"Novosibirsk":    "Asia/Novosibirsk",
	"Barnaul":        "Asia/Barnaul",
	"Tomsk":          "Asia/Tomsk",
	"Novokuznetsk":   "Asia/Novokuznetsk",
	"Krasnoyarsk":    "Asia/Krasnoyarsk",
	"Irkutsk":        "Asia/Irkutsk",
	"Chita":          "Asia/Chita",
	"Yakutsk":        "Asia/Yakutsk",
	"Khandyga":       "Asia/Khandyga",
	"Vladivostok":    "Asia/Vladivostok",
	"Ust-Nera":       "Asia/Ust-Nera",
	"Magadan":        "Asia/Magadan",
	"Sakhalin":       "Asia/Sakhalin",
	"Srednekolymsk":  "Asia/Srednekolymsk",
	"Kamchatka":      "Asia/Kamchatka",
	"Anadyr":         "Asia/Anadyr",
	"Kigali":         "Africa/Kigali",
	"Riyadh":         "Asia/Riyadh",
	"Guadalcanal":    "Pacific/Guadalcanal",
	"Mahe":           "Indian/Mahe",
	"Khartoum":       "Africa/Khartoum",
	"Stockholm":      "Europe/Stockholm",
	"Singapore":      "Asia/Singapore",
	"St_Helena":      "Atlantic/St_Helena",
	"Ljubljana":      "Europe/Ljubljana",
	"Longyearbyen":   "Arctic/Longyearbyen",
	"Bratislava":     "Europe/Bratislava",
	"Freetown":       "Africa/Freetown",
	"San_Marino":     "Europe/San_Marino",
	"Dakar":          "Africa/Dakar",
	"Mogadishu":      "Africa/Mogadishu",
	"Paramaribo":     "America/Paramaribo",
	"Juba":           "Africa/Juba",
	"Sao_Tome":       "Africa/Sao_Tome",
	"El_Salvador":    "America/El_Salvador",
	"Lower_Princes":  "America/Lower_Princes",
	"Damascus":       "Asia/Damascus",
	"Mbabane":        "Africa/Mbabane",
	"Grand_Turk":     "America/Grand_Turk",
	"Ndjamena":       "Africa/Ndjamena",
	"Kerguelen":      "Indian/Kerguelen",
	"Lome":           "Africa/Lome",
	"Bangkok":        "Asia/Bangkok",
	"Dushanbe":       "Asia/Dushanbe",
	"Fakaofo":        "Pacific/Fakaofo",
	"Dili":           "Asia/Dili",
	"Ashgabat":       "Asia/Ashgabat",
	"Tunis":          "Africa/Tunis",
	"Tongatapu":      "Pacific/Tongatapu",
	"Istanbul":       "Europe/Istanbul",
	"Port_of_Spain":  "America/Port_of_Spain",
	"Funafuti":       "Pacific/Funafuti",
	"Taipei":         "Asia/Taipei",
	"Dar_es_Salaam":  "Africa/Dar_es_Salaam",
	"Kyiv":           "Europe/Kyiv",
	"Kampala":        "Africa/Kampala",
	"Midway":         "Pacific/Midway",
	"Wake":           "Pacific/Wake",
	"New_York":       "America/New_York",
	"Detroit":        "America/Detroit",
	"Louisville":     "America/Kentucky/Louisville",
	"Monticello":     "America/Kentucky/Monticello",
	"Indianapolis":   "America/Indiana/Indianapolis",
	"Vincennes":      "America/Indiana/Vincennes",
	"Winamac":        "America/Indiana/Winamac",
	"Marengo":        "America/Indiana/Marengo",
	"Petersburg":     "America/Indiana/Petersburg",
	"Vevay":          "America/Indiana/Vevay",
	"Chicago":        "America/Chicago",
	"Tell_City":      "America/Indiana/Tell_City",
	"Knox":           "America/Indiana/Knox",
	"Menominee":      "America/Menominee",
	"Center":         "America/North_Dakota/Center",
	"New_Salem":      "America/North_Dakota/New_Salem",
	"Beulah":         "America/North_Dakota/Beulah",
	"Denver":         "America/Denver",
	"Boise":          "America/Boise",
	"Phoenix":        "America/Phoenix",
	"Los_Angeles":    "America/Los_Angeles",
	"Anchorage":      "America/Anchorage",
	"Juneau":         "America/Juneau",
	"Sitka":          "America/Sitka",
	"Metlakatla":     "America/Metlakatla",
	"Yakutat":        "America/Yakutat",
	"Nome":           "America/Nome",
	"Adak":           "America/Adak",
	"Honolulu":       "Pacific/Honolulu",
	"Montevideo":     "America/Montevideo",
	"Samarkand":      "Asia/Samarkand",
	"Tashkent":       "Asia/Tashkent",
	"Vatican":        "Europe/Vatican",
	"St_Vincent":     "America/St_Vincent",
	"Caracas":        "America/Caracas",
	"Tortola":        "America/Tortola",
	"St_Thomas":      "America/St_Thomas",
	"Ho_Chi_Minh":    "Asia/Ho_Chi_Minh",
	"Efate":          "Pacific/Efate",
	"Wallis":         "Pacific/Wallis",
	"Apia":           "Pacific/Apia",
	"Aden":           "Asia/Aden",
	"Mayotte":        "Indian/Mayotte",
	"Johannesburg":   "Africa/Johannesburg",
	"Lusaka":         "Africa/Lusaka",
	"Harare":         "Africa/Harare",
}
