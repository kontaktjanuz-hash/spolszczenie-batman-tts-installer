module spolszczenie

go 1.24

require (
	github.com/lxn/walk v0.0.0-20210112085537-c389da54e794
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e
	github.com/ulikunitz/xz v0.5.12
	golang.org/x/sys v0.30.0
)

require gopkg.in/Knetic/govaluate.v3 v3.0.0-00010101000000-000000000000 // indirect

replace golang.org/x/sys => github.com/golang/sys v0.30.0

replace gopkg.in/Knetic/govaluate.v3 => github.com/Knetic/govaluate v3.0.0+incompatible
