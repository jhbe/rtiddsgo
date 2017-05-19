package parse

type ConstsDef []ConstDef

type ConstDef struct {
	Name, Type, Value string
}

// getConstsDef traverses the result of the parsing of the IDL file ("spec")
// and extracts all constants definitions.
func getConstsDef(spec *Node) ConstsDef {
	var c ConstsDef

	spec.Traverse(func(n *Node) bool {
		return n.Kind == KindConst
	}, func(n *Node) {
		c = append(c, ConstDef{
			Name:  n.Name,
			Type:  n.TypeName,
			Value: n.Value})
	})

	return c
}
