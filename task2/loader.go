package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Film struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Year int `json:"year"`
	Director string `json:"director"`
}

type Result struct {
	Film Film
	Err error
	FilmID int
}

func worker(client *http.Client, jobs <-chan int, results chan<- Result, ctx context.Context) {
	for filmID := range jobs {
		url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", filmID)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			results <- Result{FilmID: filmID, Err: fmt.Errorf("request: %w", err)}
			continue
		}
		
		resp, err := client.Do(req)
		if err != nil {
			results <- Result{FilmID: filmID, Err: fmt.Errorf("request: %w", err)}
			continue
		}

		var film Film
		err = json.NewDecoder(resp.Body).Decode(&film)
		resp.Body.Close()

		if err != nil {
			results <- Result{FilmID: filmID, Err: fmt.Errorf("decode: %w", err)}
			continue
		}

		results <- Result{Film: film}
	}
}

func load (from int, to int, workers int, timeout int, ctx context.Context) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	numJobs := to - from + 1
	jobs := make(chan int, numJobs)
	results := make(chan Result, numJobs)

	for w := 1; w <= workers; w++ {
		go worker(client, jobs, results, ctx)
	}

	for id := from; id <= to; id++ {
		jobs <- id
	}
	close(jobs)

	for i := 0; i < numJobs; i++ {
		r := <-results
		if r.Err != nil {
			fmt.Printf("id=%d — ERROR: %v\n", r.FilmID, r.Err)
			continue
		}
		fmt.Printf("%d — %s — %d — %s\n",
			r.Film.Id, r.Film.Title, r.Film.Year, r.Film.Director)
	}
}
