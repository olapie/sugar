package plurals

// RegularReplacement is a regexp Find replacement
type RegularReplacement struct {
	Find    string
	Replace string
}

// IrregularReplacement is a hard replacement,
// containing both Singular and Plural forms
type IrregularReplacement struct {
	Singular string
	Plural   string
}

// RegularSlice is a slice of RegularReplacement Replacements
type RegularSlice []RegularReplacement

// IrregularSlice is a slice of IrregularReplacement Replacements
type IrregularSlice []IrregularReplacement

var PluralReplacements = RegularSlice{
	{"([a-z])$", "${1}s"},
	{"s$", "s"},
	{"^(ax|test)is$", "${1}es"},
	{"(octop|vir)us$", "${1}i"},
	{"(octop|vir)i$", "${1}i"},
	{"(alias|status)$", "${1}es"},
	{"(bu)s$", "${1}ses"},
	{"(buffal|tomat)o$", "${1}oes"},
	{"([ti])um$", "${1}a"},
	{"([ti])a$", "${1}a"},
	{"sis$", "ses"},
	{"(?:([^f])fe|([lr])f)$", "${1}${2}ves"},
	{"(hive)$", "${1}s"},
	{"([^aeiouy]|qu)y$", "${1}ies"},
	{"(x|ch|ss|sh)$", "${1}es"},
	{"(matr|vert|ind)(?:ix|ex)$", "${1}ices"},
	{"^(m|l)ouse$", "${1}ice"},
	{"^(m|l)ice$", "${1}ice"},
	{"^(ox)$", "${1}en"},
	{"^(oxen)$", "${1}"},
	{"(quiz)$", "${1}zes"},
}

var SingularReplacements = RegularSlice{
	{"s$", ""},
	{"(ss)$", "${1}"},
	{"(n)ews$", "${1}ews"},
	{"([ti])a$", "${1}um"},
	{"((a)naly|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)(sis|ses)$", "${1}sis"},
	{"(^analy)(sis|ses)$", "${1}sis"},
	{"([^f])ves$", "${1}fe"},
	{"(hive)s$", "${1}"},
	{"(tive)s$", "${1}"},
	{"([lr])ves$", "${1}f"},
	{"([^aeiouy]|qu)ies$", "${1}y"},
	{"(s)eries$", "${1}eries"},
	{"(m)ovies$", "${1}ovie"},
	{"(c)ookies$", "${1}ookie"},
	{"(x|ch|ss|sh)es$", "${1}"},
	{"^(m|l)ice$", "${1}ouse"},
	{"(bus)(es)?$", "${1}"},
	{"(o)es$", "${1}"},
	{"(shoe)s$", "${1}"},
	{"(cris|test)(is|es)$", "${1}is"},
	{"^(a)x[ie]s$", "${1}xis"},
	{"(octop|vir)(us|i)$", "${1}us"},
	{"(alias|status)(es)?$", "${1}"},
	{"^(ox)en", "${1}"},
	{"(vert|ind)ices$", "${1}ex"},
	{"(matr)ices$", "${1}ix"},
	{"(quiz)zes$", "${1}"},
	{"(database)s$", "${1}"},
}

var IrregularReplacements = IrregularSlice{
	{"person", "people"},
	{"man", "men"},
	{"child", "children"},
	{"sex", "sexes"},
	{"move", "moves"},
	{"mombie", "mombies"},
}

var UncountableReplacements = []string{"equipment", "information", "rice", "money", "species", "series", "fish", "sheep", "jeans", "police"}
