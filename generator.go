package resolvconf

import (
	"io"
	"text/template"
)

var templates = map[string]string{
	"domain":     "{{if .Domain.Name}}domain {{ .Domain.Name }}\n{{end}}",
	"Nameserver": "{{if .Nameservers}}{{range $nameserver := .Nameservers}}nameserver {{$nameserver.IP}}\n{{end}}\n{{end}}",
	"options":    "{{if .Options}}options{{range $opt := .Options}} {{$opt}}{{end}}\n\n{{end}}",
	"sortlist":   "{{if .Sortlist}}sortlist{{range $pair := .Sortlist}} {{$pair}}{{end}}\n\n{{end}}",
	"search":     "{{if .Search}}search{{range $dom := .Search}} {{$dom.Name}}{{end}}\n\n{{end}}",
}

// Write configuration to an io.Writer
//
// return an error if unsuccessful
func (conf *Conf) Write(w io.Writer) error {
	for _, key := range []string{"domain", "Nameserver", "sortlist", "search", "options"} {
		tmpl, err := template.New(key).Parse(templates[key])
		if err != nil {
			return err
		}

		if err := tmpl.Execute(w, conf); err != nil {
			return err
		}
	}

	return nil
}
