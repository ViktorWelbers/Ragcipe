package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorecipe/db"
	"gorecipe/links"
	"gorecipe/recipes"
	"log"
	"os"
	"sync"

	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"golang.org/x/sync/semaphore"
)

func LinkEntryPoint(wg *sync.WaitGroup) {
	sem := semaphore.NewWeighted(10)
	sum := 1
	for sum < 100 {
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go func() {
			fetchLinks := func(s string) {
				links.FetchAlLRecipeLinks(s, wg)
			}
			scrapeUrl(fmt.Sprintf("https://www.rewe.de/rezepte/?pageNumber=%d", sum), fetchLinks)
			sem.Release(1)
		}()
		sum++
	}
}

func RecipeEntryPoint(wg *sync.WaitGroup, db *sql.DB) {
	sem := semaphore.NewWeighted(10)
	file, err := os.Open("links.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Fatal(err)
		}
		link := scanner.Text()
		wg.Add(1)
		recipeFunc := func(s string) {
			recipes.FetchRecipe(s, wg, db)
		}
		go scrapeUrl(link, recipeFunc)
		break
	}
}

func scrapeUrl(url string, dataExtractorFunc func(string)) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			data := string(r.Body)
			dataExtractorFunc(data)
		},
	}).Start()
}

func main() {
	var wg sync.WaitGroup
	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat("links.txt"); err == nil {
		RecipeEntryPoint(&wg, db)
	} else if errors.Is(err, os.ErrNotExist) {
		LinkEntryPoint(&wg)
	} else {
		log.Fatal(err)
	}
	wg.Wait()
}
