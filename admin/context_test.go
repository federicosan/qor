package admin_test

import (
	"fmt"
	"testing"

	"github.com/qor/qor"
	"github.com/qor/qor/admin"

	_ "github.com/mattn/go-sqlite3"
)

// Template helpers test

func TestLinkTo(t *testing.T) {
	context := &admin.Context{Admin: Admin}

	link := context.LinkTo("test link", "/link")

	if link != "<a href=\"/link\">test link</a>" {
		t.Error("link not generated by LinkTo")
	}
}

func TestUrlForAdmin(t *testing.T) {
	context := &admin.Context{Admin: Admin}

	rootLink := context.UrlFor(Admin)

	if rootLink != "/admin" {
		t.Error("Admin link not generated by UrlFor")
	}
}

func TestUrlForResource(t *testing.T) {
	context := &admin.Context{Admin: Admin}
	user := Admin.GetResource("user")

	userLink := context.UrlFor(user)

	if userLink != "/admin/user" {
		t.Error("resource link not generated by UrlFor")
	}
}

func TestUrlForResourceName(t *testing.T) {
	user := &User{Name: "test"}
	db.Create(&user)

	context := &admin.Context{Admin: Admin, Context: &qor.Context{}}
	context.SetDB(&db)

	userLink := context.UrlFor(user)

	if userLink != "/admin/user/"+fmt.Sprintf("%v", user.Id) {
		t.Error("resource link not generated by UrlFor")
	}
}

func TestPagination(t *testing.T) {
	context := &admin.Context{Admin: Admin}
	context.Searcher = &admin.Searcher{Context: context}

	// Test current page 1
	context.Searcher.Pagination.Pages = 10
	context.Searcher.Pagination.CurrentPage = 1

	pages := context.Pagination()

	if !pages[0].Current {
		t.Error("first page not set as current page")
	}

	// +1 for "Next page" link which is a "Page" too
	if len(pages) != admin.VISIBLE_PAGE_COUNT+1 {
		t.Error("visible pages in current context beyond the bound of VISIBLE_PAGE_COUNT")
	}

	// Test current page 8 => the length between start and end less than MAX_VISIBLE_PAGES
	context.Searcher.Pagination.CurrentPage = 8
	pages = context.Pagination()

	if !pages[6].Current {
		t.Error("visible previous pages count incorrect")
	}

	// 1 for "Prev"
	if len(pages) != admin.VISIBLE_PAGE_COUNT+1 {
		t.Error("visible pages in current context beyond the bound of VISIBLE_PAGE_COUNT")
	}

	// Test current page at last
	context.Searcher.Pagination.CurrentPage = 10
	pages = context.Pagination()

	if !pages[len(pages)-1].Current {
		t.Error("last page is not the current page")
	}

	if len(pages) != admin.VISIBLE_PAGE_COUNT+1 {
		t.Error("visible pages count is incorrect")
	}

	// Test current page at last but total page count less than VISIBLE_PAGE_COUNT
	context.Searcher.Pagination.Pages = 5
	context.Searcher.Pagination.CurrentPage = 5
	pages = context.Pagination()

	if len(pages) != 5 {
		t.Error("incorrect pages count")
	}
}