package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
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
	sem := semaphore.NewWeighted(5)
	sum := 1
	for sum < 100 {
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go func() {
			scrapeUrl(fmt.Sprintf("https://www.rewe.de/rezepte/?pageNumber=%d", sum), links.FetchAlLRecipeLinks, wg)
			sem.Release(1)
		}()
		sum++
	}
}

func RecipeEntryPoint(wg *sync.WaitGroup) {
	file, err := os.Open("links.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		link := scanner.Text()
		wg.Add(1)
		go scrapeUrl(link, recipes.FetchRecipe, wg)
		break
	}
}

func main() {
	var wg sync.WaitGroup
	if _, err := os.Stat("links.txt"); err == nil {
		RecipeEntryPoint(&wg)
	} else if errors.Is(err, os.ErrNotExist) {
		LinkEntryPoint(&wg)
	} else {
		log.Fatal(err)
	}
	wg.Wait()
}

func scrapeUrl(url string, dataextractorFunc func(string, *sync.WaitGroup), wg *sync.WaitGroup) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			data := string(r.Body)
			dataextractorFunc(data, wg)
		},
	}).Start()
}
