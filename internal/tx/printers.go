package tx

import (
	"fmt"

	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
)

func (s *SafeTransaction) String() string {
	str := fmt.Sprintf(
		"To: %s\nNonce: %d\nValue: %s\n",
		db.SwapAddress(s.To),
		s.Nonce,
		s.Value,
	)
	if s.DataDecoded.Method != "" {
		str += "Data:\n"
		str += s.DataDecoded.String()
	}
	return str
}

func (d *DecodedData) String() string {
	str := ""
	if d.Method == "" {
		return ""
	}
	if d.Tabs == 0 {
		d.Tabs = 1
	}
	str += fmt.Sprintf("%s{%s}\n", printTabs(d.Tabs), d.Method)
	for _, p := range d.Params {
		p.Tabs = d.Tabs + 1
		str += p.String()
	}
	return str
}

func (p *DecodedParam) String() string {
	str := fmt.Sprintf(
		"%s- %s (%s): ",
		printTabs(p.Tabs),
		p.Name,
		p.Type,
	)
	if len(p.ValueDecoded) == 0 {
		if p.Type == "address" {
			p.Value = db.SwapAddress(fmt.Sprintf("%v", p.Value))
		}
		str += fmt.Sprintf("%v\n", p.Value)
		return str
	}
	str += "\n"
	for i, value := range p.ValueDecoded {
		v := value.Value
		if p.Type == "address" {
			v = db.SwapAddress(v)
		}
		str += fmt.Sprintf(
			"%s--- [%d/%d] ---\n%sTo: %s\n%sValue: %s\n%sData:\n",
			printTabs(p.Tabs+1),
			i+1,
			len(p.ValueDecoded),
			printTabs(p.Tabs+1),
			db.SwapAddress(value.To),
			printTabs(p.Tabs+1),
			v,
			printTabs(p.Tabs+1),
		)
		value.DataDecoded.Tabs = p.Tabs + 2
		str += value.DataDecoded.String()
	}
	return str
}

func printTabs(tabs int) string {
	str := ""
	for i := 0; i < tabs; i++ {
		str += "  "
	}
	return str
}
