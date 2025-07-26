package body_menu

import (
	"os"
	"reflect"
	"testing"

	keyvalue_menu_component "github.com/diogopereiradev/httpzen/internal/components/keyvalue_menu"
	select_menu_component "github.com/diogopereiradev/httpzen/internal/components/select_menu"
	textarea_component "github.com/diogopereiradev/httpzen/internal/components/textarea"
	"github.com/diogopereiradev/httpzen/internal/utils/http_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
)

var (
	origSelectMenu   = select_menu_component.New
	origTextarea     = textarea_component.New
	origKeyValueMenu = keyvalue_menu_component.New
)

func mockSelectMenu(selected int) func(select_menu_component.MenuImpl) {
	return func(m select_menu_component.MenuImpl) {
		m.Events.OnSelect(selected)
	}
}

func mockTextarea(value string) func(textarea_component.TextareaImpl) {
	return func(m textarea_component.TextareaImpl) {
		m.Events.OnSubmit(value)
	}
}

func mockKeyValueMenu(pairs []keyvalue_menu_component.KeyValue) func(keyvalue_menu_component.KeyValueMenuImpl) {
	return func(m keyvalue_menu_component.KeyValueMenuImpl) {
		m.OnSubmit(pairs)
	}
}

func Test_applicationJsonMenu(t *testing.T) {
	textarea_component.New = mockTextarea("{\"foo\":1}")
	defer func() { textarea_component.New = origTextarea }()
	
	got := applicationJsonMenu()
	want := []http_utility.HttpContentData{{ContentType: "application/json", Value: "{\"foo\":1}"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_applicationJsonMenu_empty(t *testing.T) {
	var exited bool
	Exit = func(code int) { exited = true }
	textarea_component.New = mockTextarea("")
	defer func() {
		textarea_component.New = origTextarea
		Exit = os.Exit
	}()

	_ = applicationJsonMenu()
	if !exited {
		t.Error("should exit on empty json")
	}
}

func Test_plainText(t *testing.T) {
	textarea_component.New = mockTextarea("abc")
	defer func() { textarea_component.New = origTextarea }()

	got := plainText()
	want := []http_utility.HttpContentData{{ContentType: "text/plain", Value: "abc"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_plainText_empty(t *testing.T) {
	var exited bool
	Exit = func(code int) { exited = true }
	textarea_component.New = mockTextarea("")
	defer func() {
		textarea_component.New = origTextarea
		Exit = os.Exit
	}()

	_ = plainText()
	if !exited {
		t.Error("should exit on empty plain text")
	}
}

func Test_xWWWFormUrlEncodedMenu(t *testing.T) {
	keyvalue_menu_component.New = mockKeyValueMenu([]keyvalue_menu_component.KeyValue{{Key: "a", Value: "b"}})
	defer func() { keyvalue_menu_component.New = origKeyValueMenu }()

	got := xWWWFormUrlEncodedMenu()
	want := []http_utility.HttpContentData{{ContentType: "application/x-www-form-urlencoded", Key: "a", Value: "b"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_multipartFormDataMenu(t *testing.T) {
	keyvalue_menu_component.New = mockKeyValueMenu([]keyvalue_menu_component.KeyValue{{Key: "x", Value: "y"}})
	defer func() { keyvalue_menu_component.New = origKeyValueMenu }()

	got := multipartFormDataMenu()
	want := []http_utility.HttpContentData{{ContentType: "multipart/form-data", Key: "x", Value: "y"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_New_all_types(t *testing.T) {
	select_menu_component.New = mockSelectMenu(0)
	textarea_component.New = mockTextarea("abc")
	keyvalue_menu_component.New = mockKeyValueMenu([]keyvalue_menu_component.KeyValue{{Key: "k", Value: "v"}})

	TermClear = func() {}
	ErrorLogger = func(string, int) {}
	Exit = func(code int) { panic(code) }
	defer func() {
		select_menu_component.New = origSelectMenu
		textarea_component.New = origTextarea
		keyvalue_menu_component.New = origKeyValueMenu
		TermClear = terminal_utility.Clear
		ErrorLogger = func(string, int) {}
		Exit = os.Exit
	}()

	types := []struct {
		idx  int
		want []http_utility.HttpContentData
	}{
		{0, []http_utility.HttpContentData{{ContentType: "application/json", Value: "abc"}}},
		{1, []http_utility.HttpContentData{{ContentType: "application/x-www-form-urlencoded", Key: "k", Value: "v"}}},
		{2, []http_utility.HttpContentData{{ContentType: "text/plain", Value: "txt"}}},
		{3, []http_utility.HttpContentData{{ContentType: "multipart/form-data", Key: "f", Value: "g"}}},
	}

	for _, typ := range types {
		select_menu_component.New = mockSelectMenu(typ.idx)
		textarea_component.New = mockTextarea("abc")
		keyvalue_menu_component.New = mockKeyValueMenu([]keyvalue_menu_component.KeyValue{{Key: "k", Value: "v"}})

		if typ.idx == 2 {
			textarea_component.New = mockTextarea("txt")
		}
		
		if typ.idx == 3 {
			keyvalue_menu_component.New = mockKeyValueMenu([]keyvalue_menu_component.KeyValue{{Key: "f", Value: "g"}})
		}

		var got []http_utility.HttpContentData
		New(nil, &got)
		if !reflect.DeepEqual(got, typ.want) {
			t.Errorf("type %d: got %v, want %v", typ.idx, got, typ.want)
		}
	}
}

func Test_New_invalid_type(t *testing.T) {
	TermClear = func() {}
	ErrorLogger = func(string, int) {}

	var exited bool
	Exit = func(code int) { exited = true }

	select_menu_component.New = mockSelectMenu(99)
	defer func() {
		select_menu_component.New = origSelectMenu
		Exit = os.Exit
		ErrorLogger = func(string, int) {}
		TermClear = terminal_utility.Clear
	}()

	var got []http_utility.HttpContentData
	New(nil, &got)
	if !exited {
		t.Error("should exit on invalid type")
	}
}
