package util

var AllColors = []string{"clear", "grey", "bluegrey", "red", "orange", "yellow", "green", "blue", "purple"}

type Theme uint8

const (
	ThemeLight Theme = iota
	ThemeDark
	endTypes
)

var AllThemes = []Theme{ThemeLight, ThemeDark}

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

func (t Theme) BodyClass() string {
	if t < endTypes {
		return themeBodyClasses[t]
	}
	return themeBodyClasses[ThemeLight]
}

func (t Theme) CardClass() string {
	if t < endTypes {
		return themeCardClasses[t]
	}
	return themeCardClasses[ThemeLight]
}

func (t Theme) Valid() bool {
	return t > ThemeLight && t < endTypes
}

var (
	themeNames = [...]string{
		ThemeLight: "Light",
		ThemeDark:  "Dark",
	}

	themeBodyClasses = [...]string{
		ThemeLight: "uk-dark",
		ThemeDark:  "uk-light",
	}

	themeCardClasses = [...]string{
		ThemeLight: "uk-card-default",
		ThemeDark:  "uk-card-secondary",
	}
)

type UserProfile struct {
	Name      string
	Theme     Theme
	NavColor  string
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
