package resolvconf

import (
	"io"
	"text/template"
)

var templates = map[string]string{
	"domain":     "{{if .GetDomain.Name}}domain {{ .GetDomain.Name }}\n{{end}}",
	"Nameserver": "{{if .GetNameservers}}{{range $nameserver := .GetNameservers}}nameserver {{$nameserver.IP}}\n{{end}}\n{{end}}",
	"options":    "{{if .GetOptions}}options{{range $opt := .GetOptions}} {{$opt}}{{end}}\n\n{{end}}",
	"sortlist":   "{{if .GetSortItems}}sortlist{{range $pair := .GetSortItems}} {{$pair}}{{end}}\n\n{{end}}",
	"search":     "{{if .GetSearchDomains}}search{{range $dom := .GetSearchDomains}} {{$dom.Name}}{{end}}\n\n{{end}}",
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
