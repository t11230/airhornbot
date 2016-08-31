package sound

// Array of all the sounds we have
var AIRHORN *SoundCollection = &SoundCollection{
	Prefix: "airhorn",
	Commands: []string{
		"airhorn",
	},
	Sounds: []*Sound{
		CreateSound("default", 1000, 250),
		CreateSound("reverb", 800, 250),
		CreateSound("spam", 800, 0),
		CreateSound("tripletap", 800, 250),
		CreateSound("fourtap", 800, 250),
		CreateSound("distant", 500, 250),
		CreateSound("echo", 500, 250),
		CreateSound("clownfull", 250, 250),
		CreateSound("clownshort", 250, 250),
		CreateSound("clownspam", 250, 0),
		CreateSound("highfartlong", 200, 250),
		CreateSound("highfartshort", 200, 250),
		CreateSound("midshort", 100, 250),
		CreateSound("truck", 10, 250),
	},
}

var KHALED *SoundCollection = &SoundCollection{
	Prefix:    "another",
	ChainWith: AIRHORN,
	Commands: []string{
		"anotha",
		"anothaone",
	},
	Sounds: []*Sound{
		CreateSound("one", 1, 250),
		CreateSound("one_classic", 1, 250),
		CreateSound("one_echo", 1, 250),
		CreateSound("dialup", 1, 250),
	},
}

var CENA *SoundCollection = &SoundCollection{
	Prefix: "jc",
	Commands: []string{
		"johncena",
		"cena",
	},
	Sounds: []*Sound{
		CreateSound("airhorn", 1, 250),
		CreateSound("echo", 1, 250),
		CreateSound("full", 1, 250),
		CreateSound("jc", 1, 250),
		CreateSound("nameis", 1, 250),
		CreateSound("spam", 1, 250),
		CreateSound("collect", 1, 250),
	},
}

var ETHAN *SoundCollection = &SoundCollection{
	Prefix: "ethan",
	Commands: []string{
		"ethan",
		"eb",
		"ethanbradberry",
		"h3h3",
	},
	Sounds: []*Sound{
		CreateSound("areyou_classic", 100, 250),
		CreateSound("areyou_condensed", 100, 250),
		CreateSound("areyou_crazy", 100, 250),
		CreateSound("areyou_ethan", 100, 250),
		CreateSound("classic", 100, 250),
		CreateSound("echo", 100, 250),
		CreateSound("high", 100, 250),
		CreateSound("slowandlow", 100, 250),
		CreateSound("cuts", 30, 250),
		CreateSound("beat", 30, 250),
		CreateSound("sodiepop", 1, 250),
		CreateSound("vape", 1, 250),
	},
}

var COW *SoundCollection = &SoundCollection{
	Prefix: "cow",
	Commands: []string{
		"stan",
		"stanislav",
	},
	Sounds: []*Sound{
		CreateSound("herd", 10, 250),
		CreateSound("moo", 10, 250),
		CreateSound("x3", 1, 250),
	},
}

var TRUMP *SoundCollection = &SoundCollection{
	Prefix: "trump",
	Commands: []string{
		"trump",
	},
	Sounds: []*Sound{
		CreateSound("10ft", 50, 250),
		// CreateSound("motivation", 50, 250),
		CreateSound("wall", 50, 250),
		CreateSound("mess", 10, 250),
		CreateSound("bing", 3, 250),
		CreateSound("getitout", 1, 250),
		CreateSound("tractor", 3, 250),
		CreateSound("worstpres", 3, 250),
		CreateSound("china", 3, 250),
		CreateSound("mexico", 3, 250),
		CreateSound("special", 3, 250),
	},
}

var MUSIC *SoundCollection = &SoundCollection{
	Prefix: "music",
	Commands: []string{
		"music",
		"m",
	},
	Sounds: []*Sound{
		CreateSound("serbian", 3, 250),
		CreateSound("techno", 3, 250),
	},
}

var MEMES *SoundCollection = &SoundCollection{
	Prefix: "meme",
	Commands: []string{
		"meme",
		"maymay",
		"memes",
	},
	Sounds: []*Sound{
		CreateSound("headshot", 3, 250),
		CreateSound("wombo", 3, 250),
		CreateSound("triple", 3, 250),
		CreateSound("camera", 3, 250),
		CreateSound("gandalf", 3, 250),
		CreateSound("mad", 50, 0),
		CreateSound("ateam", 50, 0),
		CreateSound("bennyhill", 50, 0),
		CreateSound("tuba", 50, 0),
		CreateSound("donethis", 50, 0),
		CreateSound("leeroy", 50, 0),
		CreateSound("slam", 50, 0),
		CreateSound("nerd", 50, 0),
		CreateSound("kappa", 50, 0),
		CreateSound("digitalsports", 50, 0),
		CreateSound("csi", 50, 0),
		CreateSound("nogod", 50, 0),
		CreateSound("welcomebdc", 50, 0),
	},
}

var BIRTHDAY *SoundCollection = &SoundCollection{
	Prefix: "birthday",
	Commands: []string{
		"birthday",
		"bday",
	},
	Sounds: []*Sound{
		CreateSound("horn", 50, 250),
		CreateSound("horn3", 30, 250),
		CreateSound("sadhorn", 25, 250),
		CreateSound("weakhorn", 25, 250),
	},
}

var OVERWATCH_ULTS *SoundCollection = &SoundCollection{
	Prefix: "owult",
	Commands: []string{
		"owult",
	},
	Sounds: []*Sound{
		//looking for sounds on
		//http://rpboyer15.github.io/sounds-of-overwatch/
		// CreateSound("bastion_enemy", 1000, 250),
		// CreateSound("bastion_friendly", 1000, 250),
		CreateSound("dva_enemy", 1000, 250),
		// CreateSound("dva_friendly", 1000, 250),
		CreateSound("genji_enemy", 1000, 250),
		CreateSound("genji_friendly", 1000, 250),
		CreateSound("hanzo_enemy", 1000, 250),
		// CreateSound("hanzo_enemy_wolf", 1000, 250),
		CreateSound("hanzo_friendly", 1000, 250),
		// CreateSound("hanzo_friendly_wolf", 1000, 250),
		CreateSound("junkrat_enemy", 1000, 250),
		CreateSound("junkrat_friendly", 1000, 250),
		CreateSound("lucio_friendly", 1000, 250),
		CreateSound("lucio_enemy", 1000, 250),
		CreateSound("mccree_enemy", 1000, 250),
		CreateSound("mccree_friendly", 1000, 250),
		CreateSound("mei_friendly", 1000, 250),
		// //there may be multiple mei friendly ult lines
		// //from this: https://www.reddit.com/r/Overwatch/comments/4fdw0z/is_that_ultimate_friendly_or_hostile/
		CreateSound("mei_enemy", 1000, 250),
		// CreateSound("mercy_friendly", 1000, 250),
		// CreateSound("mercy_friendly_devil", 1000, 250),
		// CreateSound("mercy_friendly_valkyrie", 1000, 250),
		// CreateSound("mercy_enemy", 1000, 250),
		CreateSound("pharah_enemy", 1000, 250),
		// CreateSound("pharah_friendly", 1000, 250),
		// CreateSound("reaper_enemy", 1000, 250), //not found
		CreateSound("reaper_friendly", 1000, 250),
		// CreateSound("reinhardt_friendly", 1000, 250), //doesn't exist?
		// CreateSound("reinhardt_enemy", 1000, 250),    //consider shortening to ?????
		// CreateSound("roadhog_enemy", 1000, 250),
		// CreateSound("roadhog_friendly", 1000, 250),
		CreateSound("76_enemy", 1000, 250), //consider shortening to s76, s:76?
		// CreateSound("76_friendly", 1000, 250),
		CreateSound("symmetra_friendly", 1000, 250),
		// CreateSound("symmetra_enemy", 1000, 250), //each hero has a line for when they see an enemy symmetra turret. not sure how to implement
		CreateSound("torbjorn_enemy", 1000, 250), //consider shortening to torb?
		// CreateSound("torbjorn_friendly", 1000, 250),
		CreateSound("tracer_enemy", 1000, 250),    //enemy line has variations. variations are an argument for splitting it up to be !owtracer, putting them in separate sound collections
		CreateSound("tracer_friendly", 1000, 250), //doesn't exist?
		CreateSound("widow_enemy", 1000, 250),     //consider shortening to widow?
		CreateSound("widow_friendly", 1000, 250),  //same as above
		CreateSound("zarya_enemy", 1000, 250),
		CreateSound("zarya_friendly", 1000, 250),
		CreateSound("zenyatta_enemy", 1000, 250),
		// CreateSound("zenyatta_friendly", 1000, 250),

		CreateSound("dva_;)", 1000, 250), //should be in its own sound repository
		CreateSound("anyong", 1000, 250),
		//skipping tracer for now

		//missing:
		//Bastion e/f
		//D.Va f
		//Genji e
		//Hanzo ew/fw
		//Junkrat f
		//Mercy e/f/fd/fv (have friendly, but different voice actress)
		//Pharah f
		//Reaper e
		//Roadhog e/f
		//Soldier:76 f
		//Symmetra e???
		//Torbjorn f
		//Tracer e/f?
		//Zenyatta f

	},
}

var DOTA *SoundCollection = &SoundCollection{
	Prefix: "dota",
	Commands: []string{
		"dota",
	},
	Sounds: []*Sound{
		CreateSound("waow", 50, 0),
		CreateSound("balance", 10, 0),
		CreateSound("rekt", 10, 0),
		CreateSound("stick", 10, 0),
		CreateSound("mana", 50, 0),
		CreateSound("disaster", 50, 0),
		CreateSound("liquid", 50, 0),
		CreateSound("history", 50, 0),
		CreateSound("smut", 50, 0),
		CreateSound("team", 50, 0),
		CreateSound("aegis", 50, 0),
	},
}

var OVERWATCH *SoundCollection = &SoundCollection{
	Prefix: "ow",
	Commands: []string{
		"overwatch",
		"ow",
	},
	Sounds: []*Sound{
		CreateSound("payload", 50, 0),
		CreateSound("whoa", 50, 0),
		CreateSound("woah", 50, 0),
		CreateSound("winky", 50, 0),
		CreateSound("turd", 50, 0),
		CreateSound("ryuugawagatekiwokurau", 50, 0),
		CreateSound("cyka", 50, 0),
		CreateSound("noon", 50, 0),
		CreateSound("somewhere", 50, 0),
		CreateSound("lift", 50, 0),
		CreateSound("russia", 50, 0),
	},
}

var WARCRAFT *SoundCollection = &SoundCollection{
	Prefix: "wc",
	Commands: []string{
		"wc3",
		"warcraft",
	},
	Sounds: []*Sound{
		CreateSound("work", 50, 0),
		CreateSound("awake", 50, 0),
	},
}

var SOUTHPARK *SoundCollection = &SoundCollection{
	Prefix: "sp",
	Commands: []string{
		"sp",
		"southpark",
	},
	Sounds: []*Sound{
		CreateSound("screw", 50, 0),
		CreateSound("authority", 50, 0),
	},
}

var SILICONVALLEY *SoundCollection = &SoundCollection{
	Prefix: "sv",
	Commands: []string{
		"sv",
		"silicon",
		"siliconvalley",
	},
	Sounds: []*Sound{
		CreateSound("piss", 50, 0),
		CreateSound("fucks", 50, 0),
		CreateSound("shittalk", 50, 0),
		CreateSound("attractive", 50, 0),
		CreateSound("win", 50, 0),
	},
}

var ARCHER *SoundCollection = &SoundCollection{
	Prefix: "archer",
	Commands: []string{
		"archer",
	},
	Sounds: []*Sound{
		CreateSound("dangerzone", 50, 0),
		CreateSound("klog", 50, 0),
	},
}

var collections []*SoundCollection = []*SoundCollection{
	AIRHORN,
	KHALED,
	CENA,
	ETHAN,
	COW,
	TRUMP,
	MUSIC,
	MEMES,
	BIRTHDAY,
	OVERWATCH_ULTS,
	DOTA,
	OVERWATCH,
	WARCRAFT,
	SOUTHPARK,
	SILICONVALLEY,
	ARCHER,
}
