package util

var AllColors = []string { "clear", "grey", "bluegrey", "red", "orange", "yellow", "green", "blue", "purple" }

type Theme uint8

const (
	ThemeLight Theme = iota
	ThemeDark
	endTypes
)

var AllThemes = []Theme{ ThemeLight, ThemeDark }

func (t Theme) String() string {
	if t < endTypes {
		return themeNames[t]
	}
	return themeNames[ThemeLight]
}

func ThemeFromString(s string) Theme {
	for _, t := range AllThemes {
		if t.String() == s {
			return t
		}
	}
	return ThemeLight
}

func (t Theme) CssClass() string {
	if t < endTypes {
		return themeClasses[t]
	}
	return themeClasses[ThemeLight]
}

func (t Theme) Valid() bool {
	return t > ThemeLight && t < endTypes
}

var (
	themeNames = [...]string{
		ThemeLight: "Light",
		ThemeDark: "Dark",
	}

	themeClasses = [...]string{
		ThemeLight: "uk-dark",
		ThemeDark: "uk-light",
	}
)

type UserProfile struct {
	Name string
	Theme Theme
	NavColor string
	LinkColor string
}

func (p *UserProfile) LinkClass() string {
	return p.LinkColor + "-fg"
}

var SystemProfile = NewUserProfile()

func NewUserProfile() UserProfile {
	return UserProfile{
		Name:      "System",
		Theme:     ThemeLight,
		NavColor:  "bluegrey",
		LinkColor: "bluegrey",
	}
}
