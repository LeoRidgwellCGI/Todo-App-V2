package actor

import (
	"context"
	"todo-app/storage"
)

type CommandType int

const (
	CreateCmd CommandType = iota
	UpdateCmd
	DeleteCmd
	ListAllCmd
	ListCmd
)

type Command struct {
	Type        CommandType
	ID          int
	Description string
	Status      string
	ResultChan  chan interface{}
}

type Actor struct {
	cmdChan chan Command
}

func NewActor(ctx context.Context) *Actor {
	actor := &Actor{
		cmdChan: make(chan Command),
	}
	go actor.run(ctx)
	return actor
}

func (a *Actor) run(ctx context.Context) {
	for cmd := range a.cmdChan {
		switch cmd.Type {
		case CreateCmd:
			// reload storage to ensure we have the latest data
			reloadStorage(ctx)

			// create the item
			item, err := storage.CreateItem(ctx, cmd.Description, cmd.Status)

			// send back result
			if err != nil {
				cmd.ResultChan <- err
			} else {
				cmd.ResultChan <- item
			}

		case UpdateCmd:
			// reload storage to ensure we have the latest data
			reloadStorage(ctx)

			// update the item
			item := storage.Item{ID: cmd.ID, Description: cmd.Description, Status: cmd.Status}
			updated, err := storage.UpdateItem(ctx, item)

			// send back result
			if err != nil {
				cmd.ResultChan <- err
			} else {
				cmd.ResultChan <- updated
			}

		case DeleteCmd:
			// reload storage to ensure we have the latest data
			reloadStorage(ctx)

			// delete the item
			err := storage.DeleteItem(ctx, cmd.ID)
			// send back result
			cmd.ResultChan <- err
		case ListAllCmd:
			// reload storage to ensure we have the latest data
			reloadStorage(ctx)

			// get all items
			items, err := storage.GetAllItems()

			// send back result
			if err != nil {
				cmd.ResultChan <- err
			} else {
				cmd.ResultChan <- items
			}
		case ListCmd:
			// reload storage to ensure we have the latest data
			reloadStorage(ctx)

			// get the item by ID
			item, err := storage.GetItemByID(cmd.ID)

			// send back result
			if err != nil {
				cmd.ResultChan <- err
			} else {
				cmd.ResultChan <- item
			}
		}
	}
}

// Create creates a new item with the given description and status.
func (a *Actor) Create(ctx context.Context, description string, status string) (storage.Item, error) {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: CreateCmd, Description: description, Status: status, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return storage.Item{}, err
	}
	return result.(storage.Item), nil
}

// Update updates an existing item with the given ID, description, and status.
func (a *Actor) Update(ctx context.Context, id int, description string, status string) (storage.Item, error) {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: UpdateCmd, ID: id, Description: description, Status: status, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return storage.Item{}, err
	}
	return result.(storage.Item), nil
}

// Delete deletes the item with the given ID.
func (a *Actor) Delete(ctx context.Context, id int) error {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: DeleteCmd, ID: id, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

// ListAll returns all items.
func (a *Actor) ListAll(ctx context.Context) (storage.Items, error) {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: ListAllCmd, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return storage.Items{}, err
	}
	return result.(storage.Items), nil
}

// List returns the item with the given ID.
func (a *Actor) List(ctx context.Context, id int) (storage.Item, error) {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: ListCmd, ID: id, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return storage.Item{}, err
	}
	return result.(storage.Item), nil
}

// Helper to reload storage before every read
func reloadStorage(ctx context.Context) {
	if storageFile := storage.GetDataFile(); storageFile != "" {
		_ = storage.Open(ctx, storageFile)
	}
}
