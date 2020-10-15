package reporting

import (
	"bytes"
	"html/template"
	"os"

	"github.com/activecm/rita/pkg/useragent"
	"github.com/activecm/rita/reporting/templates"
	"github.com/activecm/rita/resources"
)

func printUserAgents(db string, res *resources.Resources) error {
	f, err := os.Create("useragents.html")
	if err != nil {
		return err
	}
	defer f.Close()
	out, err := template.New("useragents.html").Parse(templates.UserAgentsTempl)
	if err != nil {
		return err
	}

	data, err := useragent.Results(res, 1, 1000, false)
	if err != nil {
		return err
	}

	w, err := getUserAgentsWriter(data)
	if err != nil {
		return err
	}
	return out.Execute(f, &templates.ReportingInfo{DB: db, Writer: template.HTML(w)})
}

func getUserAgentsWriter(agents []useragent.Result) (string, error) {
	tmpl := "<tr><td>{{.UserAgent}}</td><td>{{.TimesUsed}}</td></tr>\n"
	out, err := template.New("Agents").Parse(tmpl)
	if err != nil {
		return "", err
	}
	w := new(bytes.Buffer)
	for _, agent := range agents {
		err := out.Execute(w, agent)
		if err != nil {
			return "", err
		}
	}
	return w.String(), nil
}
