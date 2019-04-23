package main

import (
	"bytes"
	"github.com/ansel1/merry"
	"github.com/boltdb/bolt"
	"github.com/fpawel/elco/pkg/winapp"
	"github.com/fpawel/vimec/internal/data"
	"log"
	"path/filepath"
	"sort"
	"time"
)

type DB = *bolt.DB

func root(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("root"))
}

func dataFolderPath() (string, error) {
	dataFolderPath, err := winapp.AppDataFolderPath()
	if err != nil {
		return "", merry.Wrap(err)
	}
	dataFolderPath = filepath.Join(dataFolderPath, "ogmvik")
	err = winapp.EnsuredDirectory(dataFolderPath)
	if err != nil {
		return "", merry.Wrap(err)
	}
	return dataFolderPath, nil
}

func Import() {

	folderPath, err := dataFolderPath()
	if err != nil {
		panic(err)
	}
	fileName := filepath.Join(folderPath, "data.db")

	log.Println("file to import:", fileName)

	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	var entries []Entity

	err = db.View(func(tx *bolt.Tx) error {
		c := root(tx).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			a := readEntity(k, bytes.NewReader(v))
			entries = append(entries, a)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	sort.Slice(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		ta, tb := a.Time.Truncate(24*time.Hour), b.Time.Truncate(24*time.Hour)
		if ta == tb {
			return a.Order < b.Order
		}
		return a.Time.Before(b.Time)
	})

	log.Println(len(entries), "records to import")

	for _, entry := range entries {
		act := &data.Act{
			Year:          entry.Time.Year(),
			Month:         entry.Time.Month(),
			Day:           entry.Time.Day(),
			ActNumber:     int(entry.Order),
			DocCode:       int(entry.DocCode),
			ProductsCount: int(entry.ProductsCount),
			RouteSheet:    entry.RouteSheet,
		}
		if act.ProductsCount == 0 {
			act.ProductsCount = 60
		}
		if err := data.DB.Save(act); err != nil {
			log.Printf("ERROR: %v:\n\t%+v\n", err, entry)
		}
	}

	log.Println("ok")

	return
}
