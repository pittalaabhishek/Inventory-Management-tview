package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"github.com/rivo/tview"
)

type Item struct {
	Name string `json:"name"`
	Stock int   `json:"stock"`
}

var (
	inventory =[]Item{}
	inventoryFile = "inventory.json"
)

func loadInventory() {
	if _,err := os.Stat(inventoryFile); err == nil {{
		data, err := os.ReadFile(inventoryFile)
		if err != nil {
			log.Fatal("Error reading inventory file! - ", err)
		}
		json.Unmarshal(data, &inventory)
	}}
}

func saveInventory() {
	data, err := json.MarshalIndent(inventory, "", " ")
	if err != nil {
		log.Fatal("Error saving inventory! - ", err)
	}
	os.WriteFile(inventoryFile, data, 0644)
}

func deleteItem(index int) {
	if index < 0 || index >= len(inventory) {
		fmt.Println("Invalid item index!")
		return
	}
	inventory = append(inventory[:index], inventory[index+1:]...)
	saveInventory()
}

func main() {
	// Create a new TUI application
	app := tview.NewApplication()

	// Load existing inventory from the JSON file
	loadInventory()

	// Create a TextView that will display the inventory items in the TUI
	inventoryList := tview.NewTextView().
		SetDynamicColors(true). // Enable dynamic coloring of text
		SetRegions(true).       // Allows regions for interaction (not used here)
		SetWordWrap(true)       // Enables word wrapping to fit the TextView size

	inventoryList.SetBorder(true).SetTitle("Inventory Items") // Set border and title

	// This function refreshes the inventory display whenever there are changes
	refreshInventory := func() {
		// Clear the current content of the TextView
		inventoryList.Clear()
		// If inventory is empty, display a message
		if len(inventory) == 0 {
			fmt.Fprintln(inventoryList, "No items in inventory.")
		} else {
			// Iterate through inventory and print each item to the TextView
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (Stock: %d)\n", i+1, item.Name, item.Stock)
			}
		}
	}

	// Create input fields for item name and stock quantity
	itemNameInput := tview.NewInputField().SetLabel("Item Name: ")
	itemStockInput := tview.NewInputField().SetLabel("Stock: ")

	// Create an input field for deleting an item by its index (ID)
	itemIDInput := tview.NewInputField().SetLabel("Item ID to delete: ")

	// Create a form that lets the user add or delete items
	form := tview.NewForm().
		AddFormItem(itemNameInput).     
		AddFormItem(itemStockInput).
		AddFormItem(itemIDInput).      
		AddButton("Add Item", func() { 
			// Get the text input for name and stock
			name := itemNameInput.GetText()
			stock := itemStockInput.GetText()
			// Check if both fields are filled
			if name != "" && stock != "" {
				// Convert the stock input to an integer
				quantity, err := strconv.Atoi(stock)
				if err != nil {
					fmt.Fprintln(inventoryList, "Invalid stock value.")
					return
				}
				// Add the new item to the inventory slice
				inventory = append(inventory, Item{Name: name, Stock: quantity})
				// Save the updated inventory
				saveInventory()
				// Refresh the inventory display
				refreshInventory()
				// Clear the input fields after adding the item
				itemNameInput.SetText("")
				itemStockInput.SetText("")
			}
		}).
		AddButton("Delete Item", func() { // Button to delete an item
			idStr := itemIDInput.GetText()
			// Ensure the ID field is not empty
			if idStr == "" {
				fmt.Fprintln(inventoryList, "Please enter an item ID to delete.")
				return
			}
			// Convert the ID to an integer and check if it's valid
			id, err := strconv.Atoi(idStr)
			if err != nil || id < 1 || id > len(inventory) {
				fmt.Fprintln(inventoryList, "Invalid item ID.")
				return
			}
			// Delete the item (adjust for zero-based index)
			deleteItem(id - 1)
			fmt.Fprintf(inventoryList, "Item [%d] deleted.\n", id)
			// Refresh the inventory display after deletion
			refreshInventory()
			itemIDInput.SetText("") // Clear the ID input field
		}).
		AddButton("Exit", func() { // Button to exit the application
			app.Stop()
		})

	// Set a border and title for the form
	form.SetBorder(true).SetTitle("Manage Inventory").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		AddItem(inventoryList, 0, 1, false).
		AddItem(form, 0, 1, true)            

	refreshInventory()

	// Start the TUI application
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}