package components

import "github.com/charmbracelet/lipgloss"

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	// Palette
	Black    = lipgloss.Color("#282C34")
	Red      = lipgloss.Color("#E06C75")
	Green    = lipgloss.Color("#98C379")
	Blue     = lipgloss.Color("#61AFEF")
	Purple   = lipgloss.Color("#C678DD")
	Cyan     = lipgloss.Color("#56B6C2")
	Yellow   = lipgloss.Color("#E5C07B")
	Grey     = lipgloss.Color("#ABB2BF")
	DarkGrey = lipgloss.Color("#6B717D")
	subtle   = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	PrTitleStyle        = lipgloss.NewStyle().Bold(true)
	AdditionsStyle      = lipgloss.NewStyle().Foreground(Green)
	DeletionsStyle      = lipgloss.NewStyle().Foreground(Red)
	BackgroundStyle     = lipgloss.NewStyle().Foreground(DarkGrey)
	DocStyle            = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	DiffLineNumberStyle = lipgloss.NewStyle().Foreground(DarkGrey)
	LineColor           = Blue
	InactiveTabStyle    = lipgloss.NewStyle().Border(InactiveTabBorder, true).BorderForeground(Blue).Padding(0, 1)
	ActiveTabStyle      = InactiveTabStyle.Copy().Border(ActiveTabBorder, true).Bold(true)
	ActiveTabBlurStyle  = ActiveTabStyle.Copy().Bold(false)
	WindowStyle         = lipgloss.NewStyle().BorderForeground(Blue).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	tab                 = lipgloss.NewStyle().
				Border(InactiveTabBorder, true).
				BorderForeground(Blue).
				Padding(0, 1)
	// Set a rounded, yellow-on-purple border to the top and left
	//BoxBorderStyle = InactiveTabStyle.Copy().BorderStyle(boxBorder)
	BoxBorderStyle = InactiveTabStyle.Copy().Padding(0, 0).BorderStyle(boxBorder)
	boxBorder      = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}
	BoxTitleBorderStyle = InactiveTabStyle.Copy().BorderStyle(boxTitleBorder)
	boxTitleBorder      = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "├",
		BottomRight: "┤",
	}
	BoxBodyBorderStyle = InactiveTabStyle.Copy().BorderStyle(boxBodyBorder)
	boxBodyBorder      = lipgloss.Border{
		Top:         " ",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "│",
		TopRight:    "│",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}
	ActiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}
	InactiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
	TabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	statusNugget = lipgloss.NewStyle().
			Foreground(Grey).
			Padding(0, 1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Grey).
			Background(Black)

	StatusSectionStyle = lipgloss.NewStyle().
				Inherit(StatusBarStyle).
				Foreground(Black).
				Background(Green).
				Padding(0, 1).
				MarginRight(1)

	StatusStyle = statusNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	StatusText = lipgloss.NewStyle().Inherit(StatusBarStyle)
	CenterAll  = lipgloss.NewStyle().Align(lipgloss.Center)

	StatusHelpStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))
	HelpBox         = CenterAll.Copy()
)
