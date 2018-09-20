package client

import (
	"context"
	"testing"
)

func TestSearch(t *testing.T) {
	books, _, err := Search(context.Background(), "1", 0, 5)
	if err != nil {
		t.Errorf("Search returned unexpected error: %v", err)
	}
	if books == nil {
		t.Error("Search returned nil result")
	}
}

func TestGetBook(t *testing.T) {
	book, _, err := GetBook(context.Background(), "5825762ba6f6ad45211e65ae")
	if err != nil {
		t.Errorf("GetBook returned unexpected error: %v", err)
	}
	if book == nil {
		t.Error("GetBook returned nil result")
	}
}

func TestListTocs(t *testing.T) {
	tocs, _, err := ListTocs(context.Background(), "5825762ba6f6ad45211e65ae")
	if err != nil {
		t.Errorf("ListTocs returned unexpected error: %v", err)
	}
	if tocs == nil {
		t.Error("ListTocs returned nil result")
	}
}

func TestListChapters(t *testing.T) {
	toc, _, err := ListChapters(context.Background(), "58306f869b4ca72609ada6aa")
	if err != nil {
		t.Errorf("TestListChapters returned unexpected error: %v", err)
	}
	if toc == nil {
		t.Error("TestListChapters returned nil result")
	}
}

func TestGetChapter(t *testing.T) {
	rst, _, err := GetChapter(context.Background(), "https://www.hunhun520.com/book/zhisiwuxian/27732194.html")
	if err != nil {
		t.Errorf("GetChapter returned unexpected error: %v", err)
	}
	if rst == nil {
		t.Error("GetChapter returned nil result")
	}
}
