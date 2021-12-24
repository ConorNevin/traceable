package example

//go:generate ../bin/traceable -type=Stringer -output=traced/stringer.go

type Stringer interface {
	String() string
}
