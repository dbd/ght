package components

import (
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/charmbracelet/lipgloss"
)

func RenderBoxWithTitle(title, body string, width int) string {
	doc := strings.Builder{}
	styleTop := BoxTitleBorderStyle.Copy().Width(width)
	styleBody := BoxBodyBorderStyle.Copy().Width(width)
	doc.WriteString(styleTop.Render(title))
	doc.WriteString("\n")
	doc.WriteString(styleBody.Render(body))
	return doc.String()
}

func RenderHelpBox(body, background string, x, y, width int) string {
	rb := BoxBorderStyle.Render(CenterAll.Render(body))
	offset := lipgloss.Width(rb) / 2

	rbg := BackgroundStyle.Render(stripansi.Strip(background))
	dialog := PlaceOverlay(x-offset, y, rb, rbg)
	return dialog

}
