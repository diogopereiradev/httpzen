package body_menu

import (
	"fmt"
	"os"

	logger_module "github.com/diogopereiradev/httpzen/internal/logger"
	request_module "github.com/diogopereiradev/httpzen/internal/request"
	keyvalue_menu "github.com/diogopereiradev/httpzen/internal/utils/keyvalue_menu"
	"github.com/diogopereiradev/httpzen/internal/utils/select_menu"
	"github.com/diogopereiradev/httpzen/internal/utils/term_clear"
	"github.com/diogopereiradev/httpzen/internal/utils/textarea"
)

var Exit = os.Exit
var ErrorLogger = logger_module.Error
var TermClear = term_clear.Clear

func New(req *request_module.RequestOptions, bodyPointer *[]request_module.RequestBody) {
	TermClear()

	choices := []string{"application/json", "application/x-www-form-urlencoded", "text/plain", "multipart/form-data"}
	var contentType string

	select_menu.New(select_menu.MenuImpl{
		Choices: choices,
		Messages: select_menu.MenuMessages{
			Title:        "Select Content-Type",
			EmptyOptions: "No Content-Type options available.",
		},
		PerPage: 5,
		Events: select_menu.MenuEvents{
			OnSelect: func(choice int) {
				contentType = choices[choice]
			},
		},
	})

	var bodyResult []request_module.RequestBody

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
		ErrorLogger("Invalid Content-Type selected.")
		Exit(1)
	}

	fmt.Println(bodyResult)
}

func applicationJsonMenu() []request_module.RequestBody {
	var bodyResult []request_module.RequestBody
	textarea.New(textarea.TextareaImpl{
		Title:     "Enter JSON body",
		MaxLength: 50000,
		Events: textarea.TextareaEvents{
			OnSubmit: func(value string) {
				if value == "" {
					ErrorLogger("JSON body cannot be empty.")
					Exit(1)
					return
				}
				bodyResult = append(bodyResult, request_module.RequestBody{
					ContentType: "application/json",
					Value:       value,
				})
			},
		},
	})
	return bodyResult
}

func plainText() []request_module.RequestBody {
	var bodyResult []request_module.RequestBody
	textarea.New(textarea.TextareaImpl{
		Title:     "Enter plain text body",
		MaxLength: 50000,
		Events: textarea.TextareaEvents{
			OnSubmit: func(value string) {
				if value == "" {
					ErrorLogger("Plain text body cannot be empty.")
					Exit(1)
					return
				}
				bodyResult = append(bodyResult, request_module.RequestBody{
					ContentType: "text/plain",
					Value:       value,
				})
			},
		},
	})
	return bodyResult
}

func xWWWFormUrlEncodedMenu() []request_module.RequestBody {
	bodyResult := []request_module.RequestBody{}
	var resultPairs []keyvalue_menu.KeyValue

	keyvalue_menu.New(keyvalue_menu.KeyValueMenuImpl{
		Title:    "Add fields to application/x-www-form-urlencoded",
		OnSubmit: func(pairs []keyvalue_menu.KeyValue) {
			resultPairs = pairs
		},
	})

	for _, kv := range resultPairs {
		bodyResult = append(bodyResult, request_module.RequestBody{
			ContentType: "application/x-www-form-urlencoded",
			Key:         kv.Key,
			Value:       kv.Value,
		})
	}
	return bodyResult
}


func multipartFormDataMenu() []request_module.RequestBody {
	bodyResult := []request_module.RequestBody{}
	var resultPairs []keyvalue_menu.KeyValue

	keyvalue_menu.New(keyvalue_menu.KeyValueMenuImpl{
		Title:    "Add fields to multipart/form-data",
		OnSubmit: func(pairs []keyvalue_menu.KeyValue) {
			resultPairs = pairs
		},
	})

	for _, kv := range resultPairs {
		bodyResult = append(bodyResult, request_module.RequestBody{
			ContentType: "multipart/form-data",
			Key:         kv.Key,
			Value:       kv.Value,
		})
	}
	return bodyResult
}
