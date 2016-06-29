package main

// Array of all the sounds we have
var AIRHORN *SoundCollection = &SoundCollection{
	Prefix: "airhorn",
	Commands: []string{
		"airhorn",
	},
	Sounds: []*Sound{
		sndCreateSound("default", 1000, 250),
		sndCreateSound("reverb", 800, 250),
		sndCreateSound("spam", 800, 0),
		sndCreateSound("tripletap", 800, 250),
		sndCreateSound("fourtap", 800, 250),
		sndCreateSound("distant", 500, 250),
		sndCreateSound("echo", 500, 250),
		sndCreateSound("clownfull", 250, 250),
		sndCreateSound("clownshort", 250, 250),
		sndCreateSound("clownspam", 250, 0),
		sndCreateSound("highfartlong", 200, 250),
		sndCreateSound("highfartshort", 200, 250),
		sndCreateSound("midshort", 100, 250),
		sndCreateSound("truck", 10, 250),
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
		sndCreateSound("one", 1, 250),
		sndCreateSound("one_classic", 1, 250),
		sndCreateSound("one_echo", 1, 250),
		sndCreateSound("game", 1, 250),
		sndCreateSound("dial", 1, 250),
	},
}

var CENA *SoundCollection = &SoundCollection{
	Prefix: "jc",
	Commands: []string{
		"johncena",
		"cena",
	},
	Sounds: []*Sound{
		sndCreateSound("airhorn", 1, 250),
		sndCreateSound("echo", 1, 250),
		sndCreateSound("full", 1, 250),
		sndCreateSound("jc", 1, 250),
		sndCreateSound("nameis", 1, 250),
		sndCreateSound("spam", 1, 250),
		sndCreateSound("collect", 1, 250),
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
		sndCreateSound("areyou_classic", 100, 250),
		sndCreateSound("areyou_condensed", 100, 250),
		sndCreateSound("areyou_crazy", 100, 250),
		sndCreateSound("areyou_ethan", 100, 250),
		sndCreateSound("classic", 100, 250),
		sndCreateSound("echo", 100, 250),
		sndCreateSound("high", 100, 250),
		sndCreateSound("slowandlow", 100, 250),
		sndCreateSound("cuts", 30, 250),
		sndCreateSound("beat", 30, 250),
		sndCreateSound("sodiepop", 1, 250),
	},
}

var COW *SoundCollection = &SoundCollection{
	Prefix: "cow",
	Commands: []string{
		"stan",
		"stanislav",
	},
	Sounds: []*Sound{
		sndCreateSound("herd", 10, 250),
		sndCreateSound("moo", 10, 250),
		sndCreateSound("x3", 1, 250),
	},
}

var TRUMP *SoundCollection = &SoundCollection{
	Prefix: "trump",
	Commands: []string{
		"trump",
	},
	Sounds: []*Sound{
		sndCreateSound("mess", 10, 250),
		sndCreateSound("bing", 3, 250),
		sndCreateSound("getitout", 1, 250),
		sndCreateSound("tractor", 3, 250),
		sndCreateSound("worstpres", 3, 250),
		sndCreateSound("china", 3, 250),
		sndCreateSound("mexico", 3, 250),
		sndCreateSound("special", 3, 250),
	},
}

var MUSIC *SoundCollection = &SoundCollection{
	Prefix: "music",
	Commands: []string{
		"music",
		"m",
	},
	Sounds: []*Sound{
		sndCreateSound("serbian", 3, 250),
		sndCreateSound("techno", 3, 250),
	},
}

var MEMES *SoundCollection = &SoundCollection{
	Prefix: "memes",
	Commands: []string{
		"meme",
		"maymay",
		"memes",
	},
	Sounds: []*Sound{
		sndCreateSound("headshot", 3, 250),
		sndCreateSound("wombo", 3, 250),
		sndCreateSound("triple", 3, 250),
		sndCreateSound("camera", 3, 250),
		sndCreateSound("gandalf", 3, 250),
	},
}

var BIRTHDAY *SoundCollection = &SoundCollection{
	Prefix: "birthday",
	Commands: []string{
		"birthday",
		"bday",
	},
	Sounds: []*Sound{
		sndCreateSound("horn", 50, 250),
		sndCreateSound("horn3", 30, 250),
		sndCreateSound("sadhorn", 25, 250),
		sndCreateSound("weakhorn", 25, 250),
	},
}

var OVERWATCH *SoundCollection = &SoundCollection{
	Prefix: "owult",
	Commands: []string{
		"overwatch",
	},
	Sounds: []*Sound{
		//looking for sounds on
		//http://rpboyer15.github.io/sounds-of-overwatch/
		// sndCreateSound("bastion_enemy", 1000, 250),
		// sndCreateSound("bastion_friendly", 1000, 250),
		sndCreateSound("dva_enemy", 1000, 250),
		// sndCreateSound("dva_friendly", 1000, 250),
		// sndCreateSound("genji_enemy", 1000, 250),
		sndCreateSound("genji_friendly", 1000, 250),
		sndCreateSound("hanzo_enemy", 1000, 250),
		// sndCreateSound("hanzo_enemy_wolf", 1000, 250),
		sndCreateSound("hanzo_friendly", 1000, 250),
		// sndCreateSound("hanzo_friendly_wolf", 1000, 250),
		sndCreateSound("junkrat_enemy", 1000, 250),
		// sndCreateSound("junkrat_friendly", 1000, 250),
		sndCreateSound("lucio_friendly", 1000, 250),
		sndCreateSound("lucio_enemy", 1000, 250),
		sndCreateSound("mccree_enemy", 1000, 250),
		sndCreateSound("mccree_friendly", 1000, 250),
		sndCreateSound("mei_friendly", 1000, 250),
		// //there may be multiple mei friendly ult lines
		// //from this: https://www.reddit.com/r/Overwatch/comments/4fdw0z/is_that_ultimate_friendly_or_hostile/
		sndCreateSound("mei_enemy", 1000, 250),
		// sndCreateSound("mercy_friendly", 1000, 250),
		// sndCreateSound("mercy_friendly_devil", 1000, 250),
		// sndCreateSound("mercy_friendly_valkyrie", 1000, 250),
		// sndCreateSound("mercy_enemy", 1000, 250),
		sndCreateSound("pharah_enemy", 1000, 250),
		// sndCreateSound("pharah_friendly", 1000, 250),
		// sndCreateSound("reaper_enemy", 1000, 250), //not found
		sndCreateSound("reaper_friendly", 1000, 250),
		// sndCreateSound("reinhardt_friendly", 1000, 250), //doesn't exist?
		// sndCreateSound("reinhardt_enemy", 1000, 250),    //consider shortening to ?????
		// sndCreateSound("roadhog_enemy", 1000, 250),
		// sndCreateSound("roadhog_friendly", 1000, 250),
		sndCreateSound("76_enemy", 1000, 250), //consider shortening to s76, s:76?
		// sndCreateSound("76_friendly", 1000, 250),
		sndCreateSound("symmetra_friendly", 1000, 250),
		// sndCreateSound("symmetra_enemy", 1000, 250), //each hero has a line for when they see an enemy symmetra turret. not sure how to implement
		sndCreateSound("torbjorn_enemy", 1000, 250), //consider shortening to torb?
		// sndCreateSound("torbjorn_friendly", 1000, 250),
		// sndCreateSound("tracer_enemy", 1000, 250),    //enemy line has variations. variations are an argument for splitting it up to be !owtracer, putting them in separate sound collections
		// sndCreateSound("tracer_friendly", 1000, 250), //doesn't exist?
		sndCreateSound("widow_enemy", 1000, 250),     //consider shortening to widow?
		sndCreateSound("widow_friendly", 1000, 250),  //same as above
		// sndCreateSound("zarya_enemy", 1000, 250),
		// sndCreateSound("zarya_friendly", 1000, 250),
		sndCreateSound("zenyatta_enemy", 1000, 250),
		// sndCreateSound("zenyatta_friendly", 1000, 250),

		sndCreateSound("dva_;)", 1000, 250), //should be in its own sound repository
		sndCreateSound("anyong", 1000, 250),
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

var COLLECTIONS []*SoundCollection = []*SoundCollection{
	AIRHORN,
	KHALED,
	CENA,
	ETHAN,
	COW,
	TRUMP,
	MUSIC,
	MEMES,
	BIRTHDAY,
	OVERWATCH,
}