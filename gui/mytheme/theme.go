package mytheme

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	//go:embed fonts/ysHarmonyOS_Sans_SC_Medium.ttf
	ysfontsMedium []byte
	//go:embed fonts/ysHarmonyOS_Sans_SC_Bold.ttf
	ysfontsBlack []byte
)

type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold {
		return &fyne.StaticResource{
			StaticName:    "ysfontsBlack.ttf",
			StaticContent: ysfontsBlack,
		}
	}
	return &fyne.StaticResource{
		StaticName:    "ysfontsMedium.ttf",
		StaticContent: ysfontsMedium,
	}
}

func (m *MyTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(c, v)
}

func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
