package mtemplate

import "embed"

//go:embed mail/*
var EmailTemplateFS embed.FS
