// This file contains utility functions for finding the fully qualified name of things.

package parse

import (
	"strings"
)

// getFullName returns a string with the fully qualified name of "name"
// starting at "n". Will leave golang builtin types unchanged.
// Will return "name" if there are no modules above "n".
// Returns the name unchanged if the name is already fully qualified.
func getFullName(name string, n *Node) string {
	if n == nil {
		return ""
	}
	if isBaseType(name) {
		return name
	}
	// The names is fully qualified if it starts with the topmost module
	// name and an underscore.
	if strings.HasPrefix(name, getTopmostModuleName(n.GetTop()) + "_") {
		return name
	}
	if modulePath := modulePath(n); modulePath != "" {
		return modulePath + "_" + name
	}
	return name
}

// modulePath returns a string containing the full path of the module closest
// to "n" or an empty string if there is no module above "n" (in a tree parent
// sense) or "n" if is nil.
func modulePath(n *Node) string {
	if n == nil {
		return ""
	}

	// Move up the tree until we hit a module. Once we're at a module,
	// the fully qualified name is the full module name plus the type name we
	// started with.
	var f func(n *Node) string
	f = func(n *Node) string {
		if n == nil {
			return ""
		}
		if n.Kind == KindModule {
			return n.FullPathName("_")
		}
		return f(n.Parent)
	}

	return f(n.Parent)
}

// getTopmostModuleName finds and returns the name of the module closest to
// the top of the tree or an empty string is no module was found.
func getTopmostModuleName(spec *Node) string {
	var firstModuleName string
	spec.Traverse(func(n *Node) bool {
		return n.Kind == KindModule && firstModuleName == ""
	}, func(n *Node) {
		firstModuleName = n.Name
	})
	return firstModuleName
}

// isBaseType returns true if "name" is a golang built-in type (bool,
// string, int32, etc).
func isBaseType(name string) bool {
	switch name {
	case "float16", "float32", "float64", "int", "int16", "int32", "uint", "uint16", "uint32", "bool", "string":
		return true
	}
	return false
}

// isAnEnumType returns true if the (fully qualified) "typeName" is known as
// an enum type in "spec".
func isAnEnumType(spec *Node, typeName string) bool {
	isAnEnum := false
	spec.Traverse(func(n *Node) bool {
		return n.Kind == KindEnum && n.Name == typeName
	}, func(n *Node) {
		isAnEnum = true
	})
	return isAnEnum
}
