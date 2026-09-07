package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Task struct {
	Name string
	Done bool
}

func main() {
	a := app.New()

	w := a.NewWindow("Fyne Dashboard")

	tasks := []Task{}
	taskLabe := widget.NewLabel("Tasks: 0")
	doneLabel := widget.NewLabel("Done: 0")
	pendingLabel := widget.NewLabel("Pending: 0")

	taskEntry := widget.NewEntry()
	taskEntry.SetPlaceHolder("Enter task name...")

	taskList := container.NewVBox()

	var refreshTasks func()

	refreshTasks = func() {
		taskList.Objects = nil

		done := 0

		for i := range tasks {
			index := i

			check := widget.NewCheck(tasks[index].Name, func(checked bool) {
				tasks[index].Done = checked
				refreshTasks()
			})

			check.SetChecked(tasks[index].Done)

			taskList.Add(check)

			if tasks[index].Done {
				done++
			}
		}

		pending := len(tasks) - done

		taskLabe.SetText(
			fmt.Sprintf("Tasks: %d", len(tasks)),
		)

		doneLabel.SetText(
			fmt.Sprintf("Done: %d", done),
		)

		pendingLabel.SetText(
			fmt.Sprintf("Pending: %d", pending),
		)

		taskList.Refresh()
	}

	addButton := widget.NewButton("Add Task", func() {
		if taskEntry.Text == "" {
			return
		}

		task := Task{
			Name: taskEntry.Text,
			Done: false,
		}

		tasks = append(tasks, task)
		taskEntry.SetText("")
		refreshTasks()
	})

	// header
	header := container.NewHBox(
		widget.NewLabel("🤖 Fyne Dashboard"),
		widget.NewButton("Settings", func() {}),
	)

	// sidebar
	sidebar := container.NewVBox(
		widget.NewButton("+ New Task", func() {}),
		widget.NewButton("Tasks", func() {}),
		widget.NewButton("Statistics", func() {}),
		widget.NewButton("Settings", func() {}),
	)

	// stat
	stats := container.NewGridWithColumns(
		3,
		taskLabe,
		doneLabel,
		pendingLabel,
	)

	// loading
	loading := container.NewVBox(
		widget.NewLabel("Loading..."),
		widget.NewProgressBarInfinite(),
	)

	loading.Hide()

	// dashboard stack
	dashboardStack := container.NewStack(
		stats,
		loading,
	)

	// refresh button
	var refreshButton *widget.Button
	refreshButton = widget.NewButton("Refresh", func() {
		loading.Show()
		refreshButton.Disable()

		go func() {
			time.Sleep(3 * time.Second)
			refreshTasks()

			loading.Hide()
			refreshButton.Enable()
		}()
	})

	// main content
	mainContent := container.NewCenter(

		container.NewVBox(
			widget.NewLabel("Dashboard"),
			dashboardStack,
			refreshButton,
			taskEntry,
			addButton,
			taskList,
		),
	)

	// footer
	footer := container.NewHBox(
		widget.NewLabel("Status: Ready"),
	)

	content := container.NewBorder(
		header,
		footer,
		sidebar,
		nil,
		mainContent,
	)

	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)
	w.SetContent(content)
	w.ShowAndRun()
}
