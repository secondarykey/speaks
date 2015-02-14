package main

import ()

type setting struct {
	Database database
	Web      web
	Session  session
}

type database struct {
	Path string
}

type web struct {
	Port   string
	Root   string
	Upload string
}

type session struct {
	Secret string
	Name   string
}
