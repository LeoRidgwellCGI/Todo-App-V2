package actor

import (
	"todo-app/storage"
)

func Create(description string, status string) storage.Item {
	// create new item
	item := storage.Item{}
	return item
}

func Update(id int, description string, status string) storage.Item {
	// update item by id
	item := storage.Item{ID: id, Description: description, Status: status}
	return item
}

func Delete(id int) {
	// delete item by id
}

func ListAll() storage.Items {
	// list all items
	items := storage.Items{}
	return items
}

func List(id int) storage.Item {
	// list item by id
	item := storage.Item{ID: id}
	return item
}
