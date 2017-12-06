package resolvconf

import (
	"io"
	"text/template"
)

var templates = map[string]string{
	"domain":     "{{if .Domain.Name}}domain {{ .Domain.Name }}\n{{end}}",
	"nameserver": "{{if .Nameservers}}{{range $nameserver := .Nameservers}}nameserver {{$nameserver.IP}}\n{{end}}\n{{end}}",
	"options":    "{{if .Options}}options{{range $opt := .Options}} {{$opt}}{{end}}\n\n{{end}}",
	"sortlist":   "{{if .Sortlist.Pairs}}sortlist{{range $pair := .Sortlist.Pairs}} {{$pair}}{{end}}\n\n{{end}}",
	"search":     "{{if .Search.Domains}}search{{range $dom := .Search.Domains}} {{$dom.Name}}{{end}}\n\n{{end}}",
}

// Write configuration to an io.Writer
//
// return an error if unsuccessful
func (this *Conf) Write(w io.Writer) error {
	for _, key := range []string{"domain", "nameserver", "sortlist", "search", "options"} {
		tmpl, err := template.New(key).Parse(templates[key])
		if err != nil {
			return err
		}

		if err := tmpl.Execute(w, this); err != nil {
			return err
		}
	}

	return nil
}
