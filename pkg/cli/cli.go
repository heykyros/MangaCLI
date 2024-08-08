package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title, author string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.author }
func (i item) FilterValue() string { return i.title }

type model struct {
	list            list.Model
	textInput       textinput.Model
	quitting        bool
	choice          string
	enteringChapter bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.enteringChapter {
				m.quitting = true
				return m, tea.Quit
			} else {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = i.title
					m.enteringChapter = true
					return m, nil
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	if m.enteringChapter {
		m.textInput, cmd = m.textInput.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.quitting {

		Start(m.textInput.Value(), m.choice)
	}
	if m.enteringChapter {
		return docStyle.Render(fmt.Sprintf("Enter chapter for %s: %s", m.choice, m.textInput.View()))
	}
	return docStyle.Render(m.list.View())
}

func main() {
	items := []list.Item{
		item{title: "Attack on Titan", author: "Hajime Isayama"},
		item{title: "Black Clover", author: "Yūki Tabata"},
		item{title: "Bleach", author: "Tite Kubo"},
		item{title: "Chainsaw Man", author: "Tatsuki Fujimoto"},
		item{title: "Demon Slayer: Kimetsu no Yaiba", author: "Koyoharu Gotouge"},
		item{title: "Hunter X Hunter", author: "Yoshihiro Togashi"},
		item{title: "Jujutsu Kaisen", author: "Gege Akutami"},
		item{title: "My Hero Academia", author: "Kōhei Horikoshi"},
		item{title: "One Piece", author: "Eiichiro Oda"},
		item{title: "One-Punch Man", author: "ONE"},
		item{title: "Spy X Family", author: "Tatsuya Endo"},
	}

	ti := textinput.New()
	ti.Placeholder = "Chapter Number"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	m := model{
		list:      list.New(items, list.NewDefaultDelegate(), 0, 0),
		textInput: ti,
	}

	m.list.Title = "Manga Selector"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
