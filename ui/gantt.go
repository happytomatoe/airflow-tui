package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/airflow-tui/airflow-tui/airflow"
	"github.com/charmbracelet/lipgloss"
)

type GanttChart struct {
	width       int
	tasks       []airflow.TaskInstance
	timeScaleMs int64
}

func NewGanttChart() GanttChart {
	return GanttChart{
		width:       80,
		timeScaleMs: 1000,
	}
}

func (g *GanttChart) SetTasks(tasks []airflow.TaskInstance) {
	g.tasks = tasks
}

func (g *GanttChart) SetWidth(width int) {
	g.width = width
}

func (g *GanttChart) View() string {
	if len(g.tasks) == 0 {
		return mutedStyle.Render("No tasks to display")
	}

	minStart, maxEnd := g.findTimeRange()
	if minStart == 0 && maxEnd == 0 {
		maxEnd = 60000
	}
	duration := maxEnd - minStart
	if duration == 0 {
		duration = 60000
	}

	availableWidth := g.width - 25
	if availableWidth < 10 {
		availableWidth = 10
	}

	var lines []string
	lines = append(lines, ganttHeaderStyle.Render("Gantt Chart"))

	for _, task := range g.tasks {
		taskID := derefString(task.TaskId)
		if len(taskID) > 20 {
			taskID = taskID[:20]
		}

		startMs := int64(0)
		if task.StartDate != nil {
			if t, err := time.Parse("2006-01-02T15:04:05Z", *task.StartDate); err == nil {
				startMs = t.UnixMilli()
			}
		}
		endMs := int64(0)
		if task.EndDate != nil {
			if t, err := time.Parse("2006-01-02T15:04:05Z", *task.EndDate); err == nil {
				endMs = t.UnixMilli()
			}
		} else if task.Duration != nil {
			endMs = startMs + int64(*task.Duration*1000)
		}

		relativeStart := startMs - minStart
		relativeEnd := endMs - minStart

		if relativeStart < 0 {
			relativeStart = 0
		}
		if relativeEnd > duration {
			relativeEnd = duration
		}

		startPos := int(float64(relativeStart) / float64(duration) * float64(availableWidth))
		endPos := int(float64(relativeEnd) / float64(duration) * float64(availableWidth))

		if endPos <= startPos {
			endPos = startPos + 1
		}

		state := "-"
		if task.State != nil {
			state = string(*task.State)
		}

		bar := g.renderBar(startPos, endPos, availableWidth, state)
		taskLine := fmt.Sprintf("%-20s %s %s", taskID, stateColor(state).Render(bar), g.formatDuration(task.Duration))
		lines = append(lines, taskLine)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (g *GanttChart) findTimeRange() (int64, int64) {
	var minStart, maxEnd int64 = 0, 0

	parseTime := func(s *string) (int64, bool) {
		if s == nil {
			return 0, false
		}
		t, err := time.Parse("2006-01-02T15:04:05Z", *s)
		if err != nil {
			return 0, false
		}
		return t.UnixMilli(), true
	}

	for _, task := range g.tasks {
		if startMs, ok := parseTime(task.StartDate); ok {
			if minStart == 0 || startMs < minStart {
				minStart = startMs
			}
		}
		if endMs, ok := parseTime(task.EndDate); ok {
			if endMs > maxEnd {
				maxEnd = endMs
			}
		} else if task.Duration != nil {
			if startMs, ok := parseTime(task.StartDate); ok {
				endMs := startMs + int64(*task.Duration*1000)
				if endMs > maxEnd {
					maxEnd = endMs
				}
			}
		}
	}
	return minStart, maxEnd
}

func (g *GanttChart) renderBar(start, end, width int, state string) string {
	bar := make([]string, width)
	for i := range bar {
		bar[i] = " "
	}
	if start < width {
		for i := start; i < end && i < width; i++ {
			bar[i] = "▃"
		}
	}
	return strings.Join(bar, "")
}

func (g *GanttChart) formatDuration(d *float32) string {
	if d == nil {
		return "-"
	}
	if *d < 60 {
		return fmt.Sprintf("%.1fs", *d)
	}
	if *d < 3600 {
		return fmt.Sprintf("%.1fm", *d/60)
	}
	return fmt.Sprintf("%.1fh", *d/3600)
}

func stateColor(state string) lipgloss.Style {
	switch state {
	case "success":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
	case "failed":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	case "running":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	case "queued":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	}
}

var ganttHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("212")).
	Bold(true).
	Padding(0, 0, 1, 0)
