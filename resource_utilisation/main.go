package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	app := tview.NewApplication()

	// Create TextViews for displaying stats
	cpuUsageView := tview.NewTextView()
	cpuUsageView.SetTextAlign(tview.AlignLeft).SetBorder(true).SetTitle("CPU Usage")

	memoryUsageView := tview.NewTextView()
	memoryUsageView.SetTextAlign(tview.AlignLeft).SetBorder(true).SetTitle("Memory Usage")

	diskUsageView := tview.NewTextView()
	diskUsageView.SetTextAlign(tview.AlignLeft).SetBorder(true).SetTitle("Disk Usage")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(cpuUsageView, 3, 1, false).
		AddItem(memoryUsageView, 3, 1, false).
		AddItem(diskUsageView, 3, 1, false)

	// CPU usage updater
	go func() {
		for {
			percentages, err := cpu.Percent(0, false)
			if err == nil && len(percentages) > 0 {
				app.QueueUpdateDraw(func() {
					cpuUsageView.SetText(fmt.Sprintf("CPU Usage: %.2f%%", percentages[0]))
				})
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Memory usage updater
	go func() {
		for {
			vmStat, err := mem.VirtualMemory()
			if err == nil {
				app.QueueUpdateDraw(func() {
					memoryUsageView.SetText(fmt.Sprintf("Memory Usage: %.2f%% (Used: %.2f GB / Total: %.2f GB)",
						vmStat.UsedPercent, float64(vmStat.Used)/1e9, float64(vmStat.Total)/1e9))
				})
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Disk usage updater
	go func() {
		for {
			diskStat, err := disk.Usage("/")
			if err == nil {
				app.QueueUpdateDraw(func() {
					diskUsageView.SetText(fmt.Sprintf("Disk Usage: %.2f%% (Used: %.2f GB / Total: %.2f GB)",
						diskStat.UsedPercent, float64(diskStat.Used)/1e9, float64(diskStat.Total)/1e9))
				})
			}
			time.Sleep(1 * time.Second)
		}
	}()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
