package admin

import "html/template"

func ParseHtmlFiles(filenames ...string) (*template.Template, error) {
	if len(WebRoot) == 0 {
		return template.ParseFiles(filenames...)
	}
	tmp := make([]string, len(filenames))
	for i, fn := range filenames {
		tmp[i] = WebRoot + fn
	}
	return template.ParseFiles(tmp...)
}
