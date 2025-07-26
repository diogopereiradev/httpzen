package body_menu

import (
	"os"

	keyvalue_menu_component "github.com/diogopereiradev/httpzen/internal/components/keyvalue_menu"
	select_menu_component "github.com/diogopereiradev/httpzen/internal/components/select_menu"
	textarea_component "github.com/diogopereiradev/httpzen/internal/components/textarea"
	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	"github.com/diogopereiradev/httpzen/internal/utils/http_utility"
	"github.com/diogopereiradev/httpzen/internal/utils/terminal_utility"
)

var Exit = os.Exit
var ErrorLogger = logger_module.Error
var TermClear = terminal_utility.Clear

func New(req *request_module.RequestOptions, bodyPointer *[]http_utility.HttpContentData) {
	TermClear()

	choices := []string{"application/json", "application/x-www-form-urlencoded", "text/plain", "multipart/form-data"}
	var contentType string

	select_menu_component.New(select_menu_component.MenuImpl{
		Choices: choices,
		Messages: select_menu_component.MenuMessages{
			Title:        "Select Content-Type",
			EmptyOptions: "No Content-Type options available.",
		},
		PerPage: 5,
		Events: select_menu_component.MenuEvents{
			OnSelect: func(choice int) {
				if choice >= 0 && choice < len(choices) {
					contentType = choices[choice]
				} else {
					ErrorLogger("Invalid Content-Type selected.", 50)
					Exit(1)
				}
			},
		},
	})

	var bodyResult []http_utility.HttpContentData

	switch contentType {
	case "application/json":
		bodyResult = applicationJsonMenu()
	case "application/x-www-form-urlencoded":
		bodyResult = xWWWFormUrlEncodedMenu()
	case "text/plain":
		bodyResult = plainText()
	case "multipart/form-data":
		bodyResult = multipartFormDataMenu()
	default:
		ErrorLogger("Invalid Content-Type selected.", 50)
		Exit(1)
	}

	*bodyPointer = bodyResult
}

func applicationJsonMenu() []http_utility.HttpContentData {
	var bodyResult []http_utility.HttpContentData
	textarea_component.New(textarea_component.TextareaImpl{
		Title:     "Enter JSON body",
		MaxLength: 50000,
		Events: textarea_component.TextareaEvents{
			OnSubmit: func(value string) {
				if value == "" {
					ErrorLogger("JSON body cannot be empty.", 50)
					Exit(1)
					return
				}
				bodyResult = append(bodyResult, http_utility.HttpContentData{
					ContentType: "application/json",
					Value:       value,
				})
			},
		},
	})
	return bodyResult
}

func plainText() []http_utility.HttpContentData {
	var bodyResult []http_utility.HttpContentData
	textarea_component.New(textarea_component.TextareaImpl{
		Title:     "Enter plain text body",
		MaxLength: 50000,
		Events: textarea_component.TextareaEvents{
			OnSubmit: func(value string) {
				if value == "" {
					ErrorLogger("Plain text body cannot be empty.", 50)
					Exit(1)
					return
				}
				bodyResult = append(bodyResult, http_utility.HttpContentData{
					ContentType: "text/plain",
					Value:       value,
				})
			},
		},
	})
	return bodyResult
}

func xWWWFormUrlEncodedMenu() []http_utility.HttpContentData {
	bodyResult := []http_utility.HttpContentData{}
	var resultPairs []keyvalue_menu_component.KeyValue

	keyvalue_menu_component.New(keyvalue_menu_component.KeyValueMenuImpl{
		Title: "Add fields to application/x-www-form-urlencoded",
		OnSubmit: func(pairs []keyvalue_menu_component.KeyValue) {
			resultPairs = pairs
		},
	})

	for _, kv := range resultPairs {
		bodyResult = append(bodyResult, http_utility.HttpContentData{
			ContentType: "application/x-www-form-urlencoded",
			Key:         kv.Key,
			Value:       kv.Value,
		})
	}
	return bodyResult
}

func multipartFormDataMenu() []http_utility.HttpContentData {
	bodyResult := []http_utility.HttpContentData{}
	var resultPairs []keyvalue_menu_component.KeyValue

	keyvalue_menu_component.New(keyvalue_menu_component.KeyValueMenuImpl{
		Title: "Add fields to multipart/form-data",
		OnSubmit: func(pairs []keyvalue_menu_component.KeyValue) {
			resultPairs = pairs
		},
	})

	for _, kv := range resultPairs {
		bodyResult = append(bodyResult, http_utility.HttpContentData{
			ContentType: "multipart/form-data",
			Key:         kv.Key,
			Value:       kv.Value,
		})
	}
	return bodyResult
}
