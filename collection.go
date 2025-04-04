package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Db struct {
	tasks []Task
}

func (db *Db) Add(desc string) Task {
	task := NewTask(len(db.tasks)+1, desc)
	db.tasks = append(db.tasks, task)
	return task
}

func (db *Db) AddTask(t Task) {
	db.tasks = append(db.tasks, t)
}

func (db *Db) Delete(id int) {
	newTasks := []Task{}
	for _, v := range db.tasks {
		if v.Id == id {
			continue
		}
		newTasks = append(newTasks, v)
	}
	db.tasks = newTasks
}

func (db *Db) Update(id int, desc string) error {
	t, err := db.Find(id)
	if err != nil {
		return err
	}
	t.Description = desc
	db.Delete(id)
	db.AddTask(t)
	return nil
}

func (db *Db) Mark(id int, status Status) error {
	t, err := db.Find(id)
	if err != nil {
		return err
	}
	t.Status = status
	db.Delete(id)
	db.AddTask(t)
	return nil
}

func (db *Db) Find(id int) (Task, error) {
	for _, task := range db.tasks {
		if task.Id == id {
			return task, nil
		}
	}
	return Task{}, fmt.Errorf("no task found with id %d", id)
}

type Filter struct {
	status Status
}

func (db *Db) List(filter *Filter) []Task {
	if filter != nil {
		filtered := []Task{}
		for _, t := range db.tasks {
			if t.Status == filter.status {
				filtered = append(filtered, t)
			}
		}
		return filtered
	}
	return db.tasks
}

func (db *Db) Render(filter *Filter) string {
	res := ""
	tasks := db.List(filter)
	for _, task := range tasks {
		res = fmt.Sprintf("%s\n%s", res, task.Render())
	}
	return res
}

func (db *Db) Save() {
	file, _ := os.Create(getPath())
	defer file.Close()
	data, err := json.Marshal(db.tasks)
	if err != nil {
		fmt.Println("Unable to save", err)
		os.Exit(1)
	}
	_, werr := file.Write(data)
	if werr != nil {
		fmt.Println("Unable to write", werr)
		os.Exit(1)
	}
}

func readFromJson(r io.Reader) Db {
	decoder := json.NewDecoder(r)
	decoder.Token()

	data := []Task{}
	var task Task
	for decoder.More() {
		decoder.Decode(&task)
		data = append(data, task)
	}

	return Db{data}
}

func prepareDump(filename string) (string, error) {
	file, err := os.Open(filename)
	if err == nil {
		return filename, nil
	}
	defer file.Close()
	file, err = os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return filename, nil
}
