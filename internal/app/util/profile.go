package util


type Theme uint8

const (
	ThemeLight Theme = iota
	ThemeDark
	endTypes
)

func (t Theme) String() string {
	if t < endTypes {
		return themeNames[t]
	}
	return themeNames[ThemeLight]
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

var SystemProfile = NewUserProfile()

func NewUserProfile() UserProfile {
	return UserProfile{
		Name:      "System",
		Theme:     ThemeLight,
		NavColor:  "bluegrey",
		LinkColor: "bluegrey",
	}
}
