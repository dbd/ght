package components

import "strings"

func RenderBoxWithTitle(title, body string, width int) string {
	doc := strings.Builder{}
	styleTop := BoxTitleBorderStyle.Copy().Width(width)
	styleBody := BoxBodyBorderStyle.Copy().Width(width)
	doc.WriteString(styleTop.Render(title))
	doc.WriteString("\n")
	doc.WriteString(styleBody.Render(body))
	return doc.String()
}
