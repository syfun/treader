package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"

	"github.com/syfun/treader/pkg/logging"
)

var logger = logging.GetLogger("client")

// Book represents book info.
type Book struct {
	ID          string `json:"_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Category    string `json:"cat,omitempty"`
	Author      string `json:"author,omitempty"`
	ShortIntro  string `json:"shortIntro,omitempty"`
	LastChapter string `json:"lastChapter,omitempty"`
	WordCount   int64  `json:"wordCount,omitempty"`
}

// SearchResult represents search result.
type SearchResult struct {
	Books   []*Book `json:"books,omitempty"`
	Total   int64   `json:"total,omitempty"`
	Success bool    `json:"ok,omitempty"`
}

// Search fuzzy search books with keyword.
func Search(ctx context.Context, keyword string, start, limit int) (*SearchResult, *Response, error) {
	query := make(url.Values)
	query.Add("query", keyword)
	query.Add("start", strconv.Itoa(start))
	query.Add("limit", strconv.Itoa(limit))
	req, err := DefaultClient.NewAPIRequest("/book/fuzzy-search", query)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot search")
	}
	rst := new(SearchResult)
	resp, err := DefaultClient.Do(ctx, req, rst)
	if err != nil {
		return nil, resp, errors.Wrap(err, "cannot search")
	}
	return rst, resp, nil
}

// GetBook get book info by book id.
func GetBook(ctx context.Context, id string) (*Book, *Response, error) {
	req, err := DefaultClient.NewAPIRequest("/book/"+id, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get book")
	}
	book := new(Book)
	resp, err := DefaultClient.Do(ctx, req, book)
	if err != nil {
		return nil, resp, errors.Wrap(err, "cannot get book")
	}
	return book, resp, nil
}

// Toc represents toc info.
type Toc struct {
	ID            string     `json:"_id,omitempty"`
	Source        string     `json:"source,omitempty"`
	Name          string     `json:"name,omitempty"`
	Link          string     `json:"link,omitempty"`
	LastChapter   string     `json:"lastChapter,omitempty"`
	IsCharge      bool       `json:"isCharge,omitempty"`
	ChaptersCount int        `json:"chaptersCount,omitempty"`
	Updated       string     `json:"updated,omitempty"`
	Host          string     `json:"host,omitempty"`
	Chapters      []*Chapter `json:"chapters,omitempty"`
}

// ListTocs get tocs of book.
func ListTocs(ctx context.Context, bookID string) ([]*Toc, *Response, error) {
	query := make(url.Values)
	query.Add("view", "summary")
	query.Add("book", bookID)
	req, err := DefaultClient.NewAPIRequest("/toc", query)
	logger.Info(req.URL.String())
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot list tocs")
	}
	tocs := make([]*Toc, 0, 100)
	resp, err := DefaultClient.Do(ctx, req, &tocs)
	if err != nil {
		return nil, resp, errors.Wrap(err, "cannot list tocs")
	}
	return tocs, resp, nil
}

// Chapter represents book chapter info.
type Chapter struct {
	Title string `json:"title,omitempty"`
	Link  string `json:"link,omitempty"`
	Body  string `json:"body,omitempty"`
}

// ListChapters list toc chapters.
func ListChapters(ctx context.Context, tocID string) (*Toc, *Response, error) {
	url := fmt.Sprintf("/toc/%s?view=chapters", tocID)
	req, err := DefaultClient.NewAPIRequest(url, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot list chapters")
	}
	rst := new(Toc)
	resp, err := DefaultClient.Do(ctx, req, rst)
	if err != nil {
		return nil, resp, errors.Wrap(err, "cannot list chapters")
	}
	return rst, resp, nil
}

// ChapterResult represents get chapter result.
type ChapterResult struct {
	Success bool     `json:"ok,omitempty"`
	Chapter *Chapter `json:"chapter,omitempty"`
}

// GetChapter get chapter.
func GetChapter(ctx context.Context, link string) (*ChapterResult, *Response, error) {
	query := make(url.Values)
	query.Add("k", "2124b73d7e2e1945")
	query.Add("t", "1468223717")
	url := fmt.Sprintf("/chapter/%s", url.QueryEscape(link))
	req, err := DefaultClient.NewChatperRequest(url, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get chapter content")
	}
	rst := new(ChapterResult)
	resp, err := DefaultClient.Do(ctx, req, rst)
	if err != nil {
		return nil, resp, errors.Wrap(err, "cannot get chapter content")
	}
	return rst, resp, nil
}
