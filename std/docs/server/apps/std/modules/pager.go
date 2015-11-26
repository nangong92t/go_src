package modules

import (

)

type Pager struct {
    Total   int64           `json:"total,omitempty"`
    List    interface{}     `json:"list"`
}

func NewPager(data interface{}, total int64) *Pager {
    return &Pager{
        Total: total,
        List: data,
    }
}
