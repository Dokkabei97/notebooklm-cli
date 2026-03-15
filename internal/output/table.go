package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	HeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	CellStyle   = lipgloss.NewStyle()
	IDStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func PrintTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		fmt.Fprintln(os.Stderr, "No items found.")
		return
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8"))).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderStyle
			}
			return CellStyle
		})

	for _, row := range rows {
		t.Row(row...)
	}

	fmt.Fprintln(os.Stdout, t)
}

func PrintKeyValue(pairs [][2]string) {
	maxKeyLen := 0
	for _, p := range pairs {
		if len(p[0]) > maxKeyLen {
			maxKeyLen = len(p[0])
		}
	}
	for _, p := range pairs {
		key := lipgloss.NewStyle().Bold(true).Render(p[0])
		padding := strings.Repeat(" ", maxKeyLen-len(p[0])+2)
		fmt.Fprintf(os.Stdout, "%s%s%s\n", key, padding, p[1])
	}
}

func PrintSuccess(msg string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	fmt.Fprintln(os.Stdout, style.Render("✓ "+msg))
}

func PrintError(msg string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	fmt.Fprintln(os.Stderr, style.Render("✗ "+msg))
}

func PrintInfo(msg string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	fmt.Fprintln(os.Stdout, style.Render("ℹ "+msg))
}
