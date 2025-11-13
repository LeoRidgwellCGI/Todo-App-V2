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
			item, err := storage.CreateItem(ctx, cmd.Description, cmd.Status)
			if err != nil {
				cmd.ResultChan <- err
			} else {
				cmd.ResultChan <- item
			}
		case UpdateCmd:
			item := storage.Item{ID: cmd.ID, Description: cmd.Description, Status: cmd.Status}
			updated, err := storage.UpdateItem(ctx, item)
			if err != nil {
				cmd.ResultChan <- err
			} else {
				cmd.ResultChan <- updated
			}
		case DeleteCmd:
			err := storage.DeleteItem(ctx, cmd.ID)
			cmd.ResultChan <- err
		case ListAllCmd:
			// Return all items
			items := make(storage.Items)
			for k, v := range storage.Items(storage.Items{}) {
				items[k] = v
			}
			cmd.ResultChan <- items
		case ListCmd:
			item, err := storage.GetItemByID(cmd.ID)
			if err != nil {
				cmd.ResultChan <- err
			} else {
				cmd.ResultChan <- item
			}
		}
	}
}

// API functions
func (a *Actor) Create(ctx context.Context, description string, status string) (storage.Item, error) {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: CreateCmd, Description: description, Status: status, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return storage.Item{}, err
	}
	return result.(storage.Item), nil
}

func (a *Actor) Update(ctx context.Context, id int, description string, status string) (storage.Item, error) {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: UpdateCmd, ID: id, Description: description, Status: status, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return storage.Item{}, err
	}
	return result.(storage.Item), nil
}

func (a *Actor) Delete(ctx context.Context, id int) error {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: DeleteCmd, ID: id, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (a *Actor) ListAll(ctx context.Context) (storage.Items, error) {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: ListAllCmd, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return storage.Items{}, err
	}
	return result.(storage.Items), nil
}

func (a *Actor) List(ctx context.Context, id int) (storage.Item, error) {
	resultChan := make(chan interface{})
	a.cmdChan <- Command{Type: ListCmd, ID: id, ResultChan: resultChan}
	result := <-resultChan
	if err, ok := result.(error); ok {
		return storage.Item{}, err
	}
	return result.(storage.Item), nil
}
