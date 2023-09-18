package tx

import (
	"fmt"

	"github.com/alex-miller-0/safe-global-smartcard/internal/db"
	"github.com/alex-miller-0/safe-global-smartcard/internal/util"
)

func (p *DecodedParam) String() string {
	str := fmt.Sprintf(
		"%s- %s (%s): ",
		util.PrintTabs(p.Tabs),
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
			util.PrintTabs(p.Tabs+1),
			i+1,
			len(p.ValueDecoded),
			util.PrintTabs(p.Tabs+1),
			db.SwapAddress(value.To),
			util.PrintTabs(p.Tabs+1),
			v,
			util.PrintTabs(p.Tabs+1),
		)
		value.DataDecoded.Tabs = p.Tabs + 2
		str += value.DataDecoded.String()
	}
	return str
}
