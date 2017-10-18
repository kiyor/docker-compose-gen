/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : flag.go

* Purpose :

* Creation Date : 10-18-2017

* Last Modified : Wed 18 Oct 2017 12:30:12 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

var (
	flagAddHeader flagSliceString
)

type flagSliceString []string

func (i *flagSliceString) String() string {
	return ""
}

func (i *flagSliceString) Set(value string) error {
	*i = append(*i, value)
	return nil
}
