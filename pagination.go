package gitcode

import (
	"context"
	"iter"
)

// PageResult represents a single page of results.
type PageResult[T any] struct {
	Items []T
	Page  int
}

// Paginate creates an iterator that automatically fetches all pages of results.
// Usage:
//
//	for items, err := range client.Paginate(ctx, func(opts ListOptions) ([]*Issue, error) {
//	    return client.ListIssues(ctx, owner, repo, ListIssuesOptions{ListOptions: opts})
//	}) {
//	    if err != nil { return err }
//	    for _, issue := range items {
//	        // process issue
//	    }
//	}
func Paginate[T any](ctx context.Context, fetch func(ListOptions) ([]T, error)) iter.Seq2[[]T, error] {
	return func(yield func([]T, error) bool) {
		page := 1
		perPage := 100 // max per page for efficiency

		for {
			items, err := fetch(ListOptions{Page: page, PerPage: perPage})
			if err != nil {
				yield(nil, err)
				return
			}

			if len(items) == 0 {
				return
			}

			if !yield(items, nil) {
				return
			}

			// If we got fewer items than requested, we've reached the end
			if len(items) < perPage {
				return
			}

			page++

			// Check context cancellation
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}
}

// CollectAll collects all pages into a single slice.
func CollectAll[T any](ctx context.Context, fetch func(ListOptions) ([]T, error)) ([]T, error) {
	var all []T
	for items, err := range Paginate(ctx, fetch) {
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
	}
	return all, nil
}
