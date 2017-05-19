package parse

// EnumsDef represents a collection of enums.
type EnumsDef []EnumDef

// EnumDef represents a single enum.
type EnumDef struct {
	Name   string   // Fully qualified name of the enum.
	Values []string // Enum members. Values start at zero (same as array index).
}

// getEnumsDef traverses the result of the parsing of the IDL file ("spec")
// and extracts all enum definitions. All enum names will be fully qualified.
// Enum members, i.e. the actual values, are left unchanged. They must be
// unique already or the C DDS implementation won't work.
func getEnumsDef(spec *Node) EnumsDef {
	var e EnumsDef

	spec.Traverse(func(n *Node) bool {
		return n.Kind == KindEnum
	}, func(n *Node) {
		var members []string
		for _, child := range n.Children() {
			members = append(members, child.Name)
		}
		e = append(e, EnumDef{
			Name:   n.Name,//n.FullPathName("_"),
			Values: members})
	})

	return e
}
