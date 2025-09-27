package contents

import "lytemp/pkg/schemaparser"

func Router(h IHandler, r schemaparser.IEchoRouter) {
	g := r.Group("")
	g.GET("", h.Index)
	g.GET("/:id", h.Show)
}
