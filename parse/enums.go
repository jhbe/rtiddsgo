package parse

import "strconv"

// EnumsDef represents a collection of enums.
type EnumsDef []EnumDef

// EnumDef represents a single enum.
type EnumDef struct {
	GoName   string            // Fully qualified go name ("Module1_Module2_Name") of the enum.
	CName   string            // Fully qualified  name ("module1_Module2_name") of the enum.
	Enums  []string          // Fully qualified enum members ("Module1_Module_2_Enum1") in order. This array is the true order of the enum (the map below does not preserve order-of-addition).
	Values map[string]string // Enum values given enums.
}

// GetEnumsDef traverses the result of the parsing of the XML file ("me")
// and extracts all enum definitions. All enum names will be fully qualified.
// Enum members, i.e. the actual values, are left unchanged. They must be
// unique already or the C DDS implementation won't work.
func (me ModuleElements) GetEnumsDef() EnumsDef {
	enums := EnumsDef{}
	nextValue := 0

	me.TraverseEnums("", "", func(cPath, goPath string, e EnumDecl) {
		// Initiate a new EnumDef for this enum.
		enumDef := EnumDef{
			GoName: goNameOf(goPath, e.Name),
			CName: cNameOf(cPath, e.Name),
		}
		enumDef.Values = make(map[string]string, len(e.Enumerators))

		// Loop over the members of this enum...
		for _, s := range e.Enumerators {
			fullEnumName := goNameOf(goPath, s.Name)
			enumDef.Enums = append(enumDef.Enums, fullEnumName)

			// The enum member may or may not have a value, for example:
			//
			// enum C {
			//  C_One,
			//  C_Two,
			//  C_Three = 34,
			//  C_Four
			// };
			//
			// C_One has no value and therefore defaults to zero.
			// C_Two has no value either and is therefore assigned one.
			// C_Three does have a value "34".
			// C_Four does not and gets the last value plus one => 35.
			//
			if len(s.Value) == 0 {
				enumDef.Values[fullEnumName] = strconv.Itoa(nextValue)
				nextValue = nextValue + 1
			} else {
				enumDef.Values[fullEnumName] = s.Value
				// Calculate the next value by parsing this one and adding one.
				value, err := strconv.Atoi(s.Value)
				if err != nil {
					panic("")
				}
				nextValue = value + 1
			}
		}
		enums = append(enums, enumDef)
	})

	return enums
}
