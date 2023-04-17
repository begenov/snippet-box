package controller

import "errors"

const homeError = "Home controller error"
const createSnippetError = "Create Snippet controler error"
const chowSnippetError = "Show Snippet controller error"

var ErrNoRecord = errors.New("models: подходящей записи не найдено")
