package naming

import (
	"regexp"
	"strings"

	"code.olapie.com/sugar/v2/stringutil"

	"code.olapie.com/sugar/v2/naming/internal/plurals"
)

type replacement struct {
	regexp  *regexp.Regexp
	replace string
}

// RegularReplacement is a regexp find replacement
type RegularReplacement = plurals.RegularReplacement

// IrregularReplacement is a hard replacement,
// containing both singular and plural forms
type IrregularReplacement = plurals.IrregularReplacement

// RegularSlice is a slice of RegularReplacement Replacements
type RegularSlice = plurals.RegularSlice

// IrregularSlice is a slice of IrregularReplacement Replacements
type IrregularSlice = plurals.IrregularSlice

var compiledPluralMaps []replacement
var compiledSingularMaps []replacement

func compile() {
	compiledPluralMaps = []replacement{}
	compiledSingularMaps = []replacement{}
	for _, uncountable := range plurals.UncountableReplacements {
		inf := replacement{
			regexp:  regexp.MustCompile("^(?i)(" + uncountable + ")$"),
			replace: "${1}",
		}
		compiledPluralMaps = append(compiledPluralMaps, inf)
		compiledSingularMaps = append(compiledSingularMaps, inf)
	}

	for _, value := range plurals.IrregularReplacements {
		replacements := []replacement{
			{regexp: regexp.MustCompile(strings.ToUpper(value.Singular) + "$"), replace: strings.ToUpper(value.Plural)},
			{regexp: regexp.MustCompile(stringutil.Title(value.Singular) + "$"), replace: stringutil.Title(value.Plural)},
			{regexp: regexp.MustCompile(value.Singular + "$"), replace: value.Plural},
		}
		compiledPluralMaps = append(compiledPluralMaps, replacements...)
	}

	for _, value := range plurals.IrregularReplacements {
		replacements := []replacement{
			{regexp: regexp.MustCompile(strings.ToUpper(value.Plural) + "$"), replace: strings.ToUpper(value.Singular)},
			{regexp: regexp.MustCompile(stringutil.Title(value.Plural) + "$"), replace: stringutil.Title(value.Singular)},
			{regexp: regexp.MustCompile(value.Plural + "$"), replace: value.Singular},
		}
		compiledSingularMaps = append(compiledSingularMaps, replacements...)
	}

	for i := len(plurals.PluralReplacements) - 1; i >= 0; i-- {
		value := plurals.PluralReplacements[i]
		replacements := []replacement{
			{regexp: regexp.MustCompile(strings.ToUpper(value.Find)), replace: strings.ToUpper(value.Replace)},
			{regexp: regexp.MustCompile(value.Find), replace: value.Replace},
			{regexp: regexp.MustCompile("(?i)" + value.Find), replace: value.Replace},
		}
		compiledPluralMaps = append(compiledPluralMaps, replacements...)
	}

	for i := len(plurals.SingularReplacements) - 1; i >= 0; i-- {
		value := plurals.SingularReplacements[i]
		replacements := []replacement{
			{regexp: regexp.MustCompile(strings.ToUpper(value.Find)), replace: strings.ToUpper(value.Replace)},
			{regexp: regexp.MustCompile(value.Find), replace: value.Replace},
			{regexp: regexp.MustCompile("(?i)" + value.Find), replace: value.Replace},
		}
		compiledSingularMaps = append(compiledSingularMaps, replacements...)
	}
}

func init() {
	compile()
}

// AddPlural adds a plural replacement
func AddPlural(find, replace string) {
	plurals.PluralReplacements = append(plurals.PluralReplacements, RegularReplacement{Find: find, Replace: replace})
	compile()
}

// AddSingular adds a singular replacement
func AddSingular(find, replace string) {
	plurals.SingularReplacements = append(plurals.SingularReplacements, RegularReplacement{Find: find, Replace: replace})
	compile()
}

// AddIrregular adds an irregular replacement
func AddIrregular(singular, plural string) {
	plurals.IrregularReplacements = append(plurals.IrregularReplacements, IrregularReplacement{Singular: singular, Plural: plural})
	compile()
}

// AddUncountable adds an uncountable replacement
func AddUncountable(values ...string) {
	plurals.UncountableReplacements = append(plurals.UncountableReplacements, values...)
	compile()
}

// GetPluralReplacements retrieves the plural replacements
func GetPluralReplacements() RegularSlice {
	pluralSlice := make(RegularSlice, len(plurals.PluralReplacements))
	copy(pluralSlice, plurals.PluralReplacements)
	return pluralSlice
}

// GetSingularReplacements retrieves the singular replacements
func GetSingularReplacements() RegularSlice {
	singulars := make(RegularSlice, len(plurals.SingularReplacements))
	copy(singulars, plurals.SingularReplacements)
	return singulars
}

// GetIrregularReplacements retrieves the irregular replacements
func GetIrregularReplacements() IrregularSlice {
	irregular := make(IrregularSlice, len(plurals.IrregularReplacements))
	copy(irregular, plurals.IrregularReplacements)
	return irregular
}

// GetUncountableReplacements retrieves the uncountable replacements
func GetUncountableReplacements() []string {
	uncountables := make([]string, len(plurals.UncountableReplacements))
	copy(uncountables, plurals.UncountableReplacements)
	return uncountables
}

// SetPluralReplacements sets the plural Replacements slice
func SetPluralReplacements(Replacements RegularSlice) {
	plurals.PluralReplacements = Replacements
	compile()
}

// SetSingularReplacements sets the singular Replacements slice
func SetSingularReplacements(Replacements RegularSlice) {
	plurals.SingularReplacements = Replacements
	compile()
}

// SetIrregularReplacements sets the irregular Replacements slice
func SetIrregularReplacements(Replacements IrregularSlice) {
	plurals.IrregularReplacements = Replacements
	compile()
}

// SetUncountableReplacements sets the uncountable Replacements slice
func SetUncountableReplacements(Replacements []string) {
	plurals.UncountableReplacements = Replacements
	compile()
}

// Plural converts a word to its plural form
func Plural(str string) string {
	for _, replacement := range compiledPluralMaps {
		if replacement.regexp.MatchString(str) {
			return replacement.regexp.ReplaceAllString(str, replacement.replace)
		}
	}
	return str
}

// Singular converts a word to its singular form
func Singular(str string) string {
	for _, replacement := range compiledSingularMaps {
		if replacement.regexp.MatchString(str) {
			return replacement.regexp.ReplaceAllString(str, replacement.replace)
		}
	}
	return str
}
