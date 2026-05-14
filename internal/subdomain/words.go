package subdomain

// adjectives and nouns yield ~4 800 slug combinations.
var adjectives = []string{
	"amber", "azure", "bare", "bold", "brave", "bright", "brisk", "calm",
	"clean", "clear", "cold", "cool", "crisp", "dark", "deep", "deft",
	"dry", "dusk", "eager", "fair", "fast", "firm", "fleet", "fresh",
	"full", "grand", "great", "green", "grey", "hardy", "high", "keen",
	"kind", "large", "lean", "light", "long", "loud", "mild", "neat",
	"new", "nice", "noble", "odd", "open", "plain", "proud", "pure",
	"quick", "quiet", "rare", "rich", "safe", "sharp", "short", "silent",
	"slim", "slow", "small", "smart", "smooth", "soft", "still", "strong",
	"swift", "tall", "thick", "thin", "true", "warm", "wide", "wild",
}

var nouns = []string{
	"bay", "bird", "bluff", "brook", "cape", "cave", "cliff", "cloud",
	"coast", "cove", "crest", "creek", "dale", "dawn", "delta", "drift",
	"dune", "dusk", "dust", "echo", "falls", "fen", "field", "floe",
	"fog", "ford", "frost", "gale", "glen", "grove", "gust", "haze",
	"hill", "ice", "isle", "lake", "leaf", "marsh", "mesa", "mist",
	"moon", "moor", "moss", "oak", "peak", "pine", "pond", "pool",
	"rain", "reef", "ridge", "rift", "rise", "river", "rock", "root",
	"shore", "sky", "slate", "snow", "star", "stone", "storm", "stream",
	"surge", "tide", "trail", "vale", "vault", "wave", "wind", "wood",
}

// words yields single-word subdomains (must all satisfy IsValid).
var words = []string{
	"apex", "arch", "ark", "atlas", "aura", "axis",
	"beam", "blaze", "bloom", "bolt", "bond", "bridge",
	"calm", "crest", "dawn", "delta", "drift", "dune",
	"echo", "edge", "ember", "falcon", "fern", "flame",
	"flare", "flint", "flow", "foam", "forge", "frost",
	"gale", "gem", "glade", "glow", "grove", "gust",
	"haven", "haze", "hive", "iris", "isle", "jade",
	"lance", "lark", "leap", "lens", "link", "lynx",
	"mesa", "mist", "moon", "moss", "nova", "opal",
	"orb", "peak", "pine", "pivot", "plume", "pool",
	"prism", "pulse", "quest", "reef", "ridge", "rift",
	"ring", "rise", "root", "rune", "rush", "sage",
	"scale", "scout", "shade", "shore", "slate", "span",
	"spark", "spire", "star", "stem", "storm", "surge",
	"swift", "tide", "torch", "trail", "vale", "vault",
	"vine", "wake", "ward", "wave", "wisp", "wolf",
	"zeal", "zone",
}
