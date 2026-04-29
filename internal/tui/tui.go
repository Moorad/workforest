package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const listHeight = 14

type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

func newStyles() styles {
	var s styles
	s.title = lipgloss.NewStyle().MarginLeft(2)
	s.item = lipgloss.NewStyle().PaddingLeft(4)
	s.item.Align(lipgloss.Center)
	s.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	s.pagination = list.DefaultStyles(true).PaginationStyle.PaddingLeft(4)
	s.help = list.DefaultStyles(true).HelpStyle.PaddingLeft(4).PaddingBottom(1)
	s.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	return s
}

type Item struct {
	Label string
	Value string
}

func NewItem(label string, value string) Item {
	return Item{
		Label: label,
		Value: value,
	}
}

func (i Item) FilterValue() string { return string(i.Label) }

type itemDelegate struct {
	styles *styles
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Label)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list       list.Model
	choice     Item
	styles     styles
	quitting   bool
	termWidth  int
	termHeight int
}

func initialModel(prompt string, choices []list.Item) model {
	const defaultWidth = 20

	l := list.New(choices, itemDelegate{}, defaultWidth, listHeight)
	l.Title = prompt
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	m := model{list: l}
	m.updateStyles() // default to dark styles.
	return m
}

func (m *model) updateStyles() {
	m.styles = newStyles()
	m.list.Styles.Title = m.styles.title
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.Styles.HelpStyle = m.styles.help
	m.list.SetDelegate(itemDelegate{styles: &m.styles})
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.choice = i
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	centered := lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, m.list.View())
	v := tea.NewView(centered)
	v.AltScreen = true
	return v
}

func PromptList(prompt string, options []Item) Item {
	items := []list.Item{}
	for _, opt := range options {
		items = append(items, Item(opt))
	}

	m := initialModel(prompt, items)

	result, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return result.(model).choice
}
