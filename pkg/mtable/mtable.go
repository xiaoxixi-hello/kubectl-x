package mtable

import (
    "fmt"

    "github.com/liushuochen/gotable"
    "github.com/liushuochen/gotable/table"
)

var title []string = []string{
    "NAME",
    "STATUS",
    "IP",
    "SVC",
    "INGRESS",
}

func GenTable(mapList []map[string]string) *table.Table {
    t, err := gotable.Create(title...)
    if err != nil {
        fmt.Printf("create table error: %s", err.Error())
        return nil
    }
    t.AddRows(mapList)
    return t
}
