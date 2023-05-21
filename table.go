package a3spow

import (
	"bytes"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func RenderTable(addresses []OutAddress) []byte {
	buf := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"#", "address", "salt"})
	for idx, address := range addresses {
		table.Append([]string{
			strconv.Itoa(idx + 1), address.Address.Hex(), address.Salt,
		})
	}
	table.Render()

	return buf.Bytes()
}
