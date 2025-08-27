package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"ui_example/ui"
	"ui_example/ui/providers"
	"ui_example/ui/styles"
)

// CustomContentProvider shows how to create your own content
type CustomContentProvider struct {
	appName string
}

func NewCustomContentProvider(appName string) *CustomContentProvider {
	return &CustomContentProvider{
		appName: appName,
	}
}

func (c *CustomContentProvider) RenderContent(width, height int) string {
	t := styles.CurrentTheme()
	
	availableWidth := width - 6
	availableHeight := height - 6
	
	if availableWidth <= 0 || availableHeight <= 0 {
		return ""
	}
	
	var sections []string
	
	// Custom welcome message
	welcome := styles.ApplyBoldForegroundGrad(c.appName, t.Primary, t.Secondary)
	sections = append(sections, welcome)
	sections = append(sections, t.S().Muted.Render("Custom Application Example"))
	sections = append(sections, "")
	
	// Custom content
	content := []string{
		"This demonstrates how to create a custom application",
		"using the reusable CRUSH UI framework:",
		"",
		"• Custom content provider",
		"• Custom sidebar sections", 
		"• Custom header data",
		"• Maintains all CRUSH styling",
		"",
		"The framework handles all the layout, styling,",
		"and responsive behavior automatically.",
	}
	
	for _, line := range content {
		if strings.HasPrefix(line, "•") {
			sections = append(sections, t.S().Success.Render(line))
		} else if line == "" {
			sections = append(sections, "")
		} else {
			sections = append(sections, t.S().Text.Render(line))
		}
	}
	
	// Center content
	contentStr := strings.Join(sections, "\n")
	contentLines := strings.Split(contentStr, "\n")
	
	if len(contentLines) < availableHeight {
		paddingTop := (availableHeight - len(contentLines)) / 2
		for i := 0; i < paddingTop; i++ {
			contentLines = append([]string{""}, contentLines...)
		}
	}
	
	var centeredLines []string
	for _, line := range contentLines {
		lineWidth := len(line) // Simple width calculation
		if lineWidth < availableWidth {
			padding := (availableWidth - lineWidth) / 2
			line = strings.Repeat(" ", padding) + line
		}
		centeredLines = append(centeredLines, line)
	}
	
	return strings.Join(centeredLines, "\n")
}

func (c *CustomContentProvider) HandleContentUpdate(msg tea.Msg) tea.Cmd {
	return nil
}

func (c *CustomContentProvider) InitContent() tea.Cmd {
	return nil
}

// CustomHeaderProvider shows how to customize header data
type CustomHeaderProvider struct {
	appName    string
	lastUpdate time.Time
}

func NewCustomHeaderProvider(appName string) *CustomHeaderProvider {
	return &CustomHeaderProvider{
		appName:    appName,
		lastUpdate: time.Now(),
	}
}

func (h *CustomHeaderProvider) GetBrandName() string {
	return "MyApp™"
}

func (h *CustomHeaderProvider) GetAppName() string {
	return h.appName
}

func (h *CustomHeaderProvider) GetStatusData() map[string]interface{} {
	return map[string]interface{}{
		"time":    h.lastUpdate.Format("15:04:05"),
		"status":  "custom",
		"users":   42,
		"version": "v2.1.0",
	}
}

func (h *CustomHeaderProvider) HandleHeaderUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tickMsg:
		h.lastUpdate = time.Time(msg)
		return tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return nil
}

func (h *CustomHeaderProvider) InitHeader() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

// CustomTasksSection shows how to create a custom sidebar section
type CustomTasksSection struct {
	tasks []Task
}

type Task struct {
	Name     string
	Status   string
	Priority string
}

func NewCustomTasksSection() *CustomTasksSection {
	return &CustomTasksSection{
		tasks: []Task{
			{Name: "Design Review", Status: "in-progress", Priority: "high"},
			{Name: "Code Review", Status: "pending", Priority: "medium"},
			{Name: "Testing", Status: "completed", Priority: "high"},
			{Name: "Documentation", Status: "pending", Priority: "low"},
			{Name: "Deployment", Status: "blocked", Priority: "high"},
		},
	}
}

func (t *CustomTasksSection) GetTitle() string {
	return "Tasks"
}

func (t *CustomTasksSection) RenderItems(maxItems, width int) []providers.SidebarItem {
	var items []providers.SidebarItem
	
	count := 0
	for _, task := range t.tasks {
		if count >= maxItems {
			break
		}
		
		icon := t.getTaskIcon(task.Status)
		status := t.getTaskStatus(task.Status)
		
		name := task.Name
		if len(name) > width-10 {
			name = name[:width-13] + "..."
		}
		
		items = append(items, providers.SidebarItem{
			Icon:   icon,
			Text:   name,
			Value:  task.Priority,
			Status: status,
		})
		count++
	}
	
	return items
}

func (t *CustomTasksSection) getTaskIcon(status string) string {
	switch status {
	case "completed":
		return "✓"
	case "in-progress":
		return "●"
	case "pending":
		return "○"
	case "blocked":
		return "×"
	default:
		return "○"
	}
}

func (t *CustomTasksSection) getTaskStatus(status string) string {
	switch status {
	case "completed":
		return "success"
	case "in-progress":
		return "info"
	case "pending":
		return "warning"
	case "blocked":
		return "error"
	default:
		return "muted"
	}
}

func (t *CustomTasksSection) HandleSectionUpdate(msg tea.Msg) tea.Cmd {
	return nil
}

func (t *CustomTasksSection) InitSection() tea.Cmd {
	return nil
}

func (t *CustomTasksSection) RefreshSection() tea.Cmd {
	// Simulate task updates
	now := time.Now()
	if now.Second()%10 == 0 && len(t.tasks) > 0 {
		// Toggle first task status
		if t.tasks[0].Status == "in-progress" {
			t.tasks[0].Status = "completed"
		} else {
			t.tasks[0].Status = "in-progress"
		}
	}
	return nil
}

func main() {
	// Create custom providers
	contentProvider := NewCustomContentProvider("CUSTOM APP")
	headerProvider := NewCustomHeaderProvider("Custom App")
	
	// Create custom sidebar sections
	customSections := []providers.SidebarSection{
		NewCustomTasksSection(),
		providers.NewServersSection(), // Mix custom with default sections
		providers.NewStatusSection(),
	}
	
	// Configure the app
	config := &providers.AppConfig{
		ContentProvider:             contentProvider,
		HeaderDataProvider:          headerProvider,
		SidebarSections:             customSections,
		ShowSidebarByDefault:        true,
		CompactModeWidthBreakpoint:  100, // Custom breakpoint
		CompactModeHeightBreakpoint: 25,
	}
	
	// Create the app with custom configuration
	app := ui.NewReusableApp(config)
	
	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		log.Fatal(err)
		os.Exit(1)
	}
}