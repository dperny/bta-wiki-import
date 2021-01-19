package export

import (
	"fmt"
	"strings"
)

type WikiTemplate struct {
	name string
	args []templateArg
}

type templateArg struct {
	name  string
	value string
}

func NewWikiTemplate(name string) *WikiTemplate {
	return &WikiTemplate{
		name: name,
	}
}

func (wt *WikiTemplate) AddArg(name string, value interface{}) {
	switch v := value.(type) {
	case int:
		wt.args = append(wt.args, templateArg{name: name, value: fmt.Sprint(v)})
	case float64:
		wt.args = append(wt.args, templateArg{name: name, value: fmt.Sprint(v)})
	case bool:
		wt.args = append(wt.args, templateArg{name: name, value: fmt.Sprint(v)})
	case string:
		wt.args = append(wt.args, templateArg{name: name, value: v})
	}
}

func (wt *WikiTemplate) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "{{%s\n", wt.name)
	for _, arg := range wt.args {
		// skip adding any args that are empty strings
		if arg.value != "" {
			fmt.Fprintf(&b, "|%s=%s\n", arg.name, arg.value)
		}
	}
	fmt.Fprintln(&b, "}}")

	return b.String()
}
