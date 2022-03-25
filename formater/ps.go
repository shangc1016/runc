package formater

import (
	"os"

	"gitee.com/shangc1016/runc/status"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PsQuiet(ids []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Container-ID"})
	for k, id := range ids {
		t.AppendRow(table.Row{k, id})
		t.AppendSeparator()
	}
	t.Render()
}

func PsAll(info []status.ContainerInfo) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "id", "name", "command", "pid", "status", "volume:(outer, inner)"})
	for k, item := range info {
		t.AppendRow(table.Row{k, item.Id, item.Name, item.Command, item.Pid, item.Status, item.Volumes})
		t.AppendSeparator()
	}
	t.Render()
}
