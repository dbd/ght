package components

import (
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/charmbracelet/lipgloss"
)

func RenderBoxWithTitle(title, body string, width int) string {
	return RenderBoxWithTitleCorner(title, body, width, false, false)
}

func RenderBox(body string, width int) string {
	doc := strings.Builder{}
	bs := BoxBorderStyle.Copy().Width(width)
	doc.WriteString(bs.Render(body))
	return doc.String()
}

func RenderBoxWithTitleCorner(title, body string, width int, topLeft, bottomLeft bool) string {
	doc := strings.Builder{}
	var styleTop lipgloss.Style
	var styleBody lipgloss.Style
	styleTop = BoxTitleBorderStyle.Copy().Width(width)
	styleBody = BoxBodyBorderStyle.Copy().Width(width)
	if topLeft {
		border := boxTitleBorder
		border.TopLeft = "├"
		styleTop = styleTop.Copy().BorderStyle(border)
	}
	if bottomLeft {
		border := boxBodyBorder
		border.BottomLeft = "├"
		styleBody = styleBody.Copy().BorderStyle(border)
	}
	doc.WriteString(styleTop.Render(title))
	doc.WriteString("\n")
	doc.WriteString(styleBody.Render(body))
	return doc.String()
}

func RenderHelpBox(body, background string, x, y, width int) string {
	rb := BoxBorderStyle.Width(lipgloss.Width(body)).Render(CenterAll.Render(body))
	offset := lipgloss.Width(rb) / 2

	rbg := BackgroundStyle.Render(stripansi.Strip(background))
	dialog := PlaceOverlay(x-offset, y, rb, rbg)
	return dialog

}

func RenderFilter(body, background string, x, y, width int) string {
	rb := BoxBorderStyle.Width(width).Render(CenterAll.Render(body))
	offset := lipgloss.Width(rb) / 2

	rbg := BackgroundStyle.Render(stripansi.Strip(background))
	dialog := PlaceOverlay(x-offset, y, rb, rbg)
	return dialog

}

func RenderColoredText(text, cs string) string {
	cm := map[string]lipgloss.Color{
		"black":    Black,
		"red":      Red,
		"green":    Green,
		"blue":     Blue,
		"purple":   Purple,
		"cyan":     Cyan,
		"yellow":   Yellow,
		"grey":     Grey,
		"darkGrey": DarkGrey,
	}
	var c lipgloss.Color
	cs = strings.ToLower(cs)
	c, ok := cm[cs]
	if !ok {
		c = lipgloss.Color("#" + cs)
	}
	s := lipgloss.NewStyle().Foreground(c).Render(text)
	return s
}
