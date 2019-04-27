package parse

import "testing"

func TestCNameOf(t *testing.T) {
	cases := []struct {
		path, name, expected string
	} {
		{"", "", ""},
		{"A", "", ""},
		{"", "B", "B"},
		{"_", "B", "B"},
		{"a", "b", "a_b"},
		{"A", "B", "A_B"},
		{"A_", "B", "A_B"},
		{"_a_", "B", "a_B"},
		{"a_B", "C", "a_B_C"},
		{"Ab_cD_Eef", "Gh", "Ab_cD_Eef_Gh"},
		{"", "Ab::Cd", "Ab_Cd"},
		{"Ab::cD::Eef", "gh", "Ab_cD_Eef_gh"},
	}

	for ix, c := range cases {
		out := cNameOf(c.path, c.name)
		if out != c.expected {
			t.Errorf("Case %d: Expected %s, got %s", ix, c.expected, out)
		}
	}
}

func TestGoNameOf(t *testing.T) {
	cases := []struct {
		path, name, expected string
	} {
		{"", "", ""},
		{"A", "", ""},
		{"", "B", "B"},
		{"_", "B", "B"},
		{"a", "B", "A_B"},
		{"A", "B", "A_B"},
		{"A_", "B", "A_B"},
		{"_a_", "B", "A_B"},
		{"A_B", "C", "A_B_C"},
		{"Ab_cD_Eef", "Gh", "Ab_CD_Eef_Gh"},
		{"", "Ab::Cd", "Ab_Cd"},
		{"Ab::cD::Eef", "gh", "Ab_CD_Eef_Gh"},
	}

	for ix, c := range cases {
		out := goNameOf(c.path, c.name)
		if out != c.expected {
			t.Errorf("Case %d: Expected %s, got %s", ix, c.expected, out)
		}
	}
}

func TestToTitle(t *testing.T) {
	cases := []struct {
		in, expected string
	} {
		{"", ""},
		{"A", "A"},
		{"a", "A"},
		{"abc", "Abc"},
		{"aBc", "ABc"},
		{"aBc::EfG", "ABc::EfG"},
		{"aBc::EfG::hI", "ABc::EfG::HI"},
	}

	for ix, c := range cases {
		out := toTitle(c.in)
		if out != c.expected {
			t.Errorf("Case %d: Expected %s, got %s", ix, c.expected, out)
		}
	}
}

func TestIsAQualifiedValue(t *testing.T) {
	cases := [] struct {
		in       string
		expected bool
	}{
		{"", false},
		{"A", false},
		{"(A)", false},
		{"(com::A)", true},
		{"(com:A)", false},
		{"com::A", true},
	}

	for ix, c := range cases {
		out := isAQualifiedValue(c.in)
		if out != c.expected {
			t.Errorf("Case %d: expected %v, got %v", ix, c.expected, out)
		}
	}
}

func TestGoTypeOf(t *testing.T) {
	cases := []struct{
		t, nonBasic, expected string
	}{
		{"", "", ""},
		{"boolean", "", "bool"},
		{"int16", "", "int16"},
		{"uint16", "", "uint16"},
		{"int32", "", "int32"},
		{"uint32", "", "uint32"},
		{"float32", "", "float32"},
		{"float64", "", "float64"},

		{"nonBasic", "b", "B"},
		{"nonBasic", "a::b", "A_B"},
		{"nonBasic", "a::b::c", "A_B_C"},
	}

	for ix, c := range cases {
		out := goTypeOf(c.t, c.nonBasic)
		if out != c.expected {
			t.Errorf("Case %d: Expected %s, got %s", ix, c.expected, out)
		}
	}
}

func TestDdsTypesOf(t *testing.T) {
	cases := []struct {
		t, nb, out string
	}{
		{"", "", ""},
		{"boolean", "", "DDS_Boolean"},
		{"int16", "", "DDS_Short"},
		{"uint16", "", "DDS_UnsignedShort"},
		{"int32", "", "DDS_Long"},
		{"uint32", "", "DDS_UnsignedLong"},
		{"float32", "", "DDS_Float"},
		{"float64", "", "DDS_Double"},
		{"string", "", "DDS_String"},
		{"foo", "", "foo"},
	}

	for ix, c := range cases {
		result := ddsTypeOf(c.t, c.nb)
		if result != c.out {
			t.Errorf("case %d: Expected %s, got %s", ix, c.out, result)
		}
	}
}

func TestXmlTypeOf(t *testing.T) {
	cases := []struct {
		t, nb string
		expected string
	} {
		{"", "", ""},
		{"one", "two", "one"},
		{"nonbasic", "two", "nonbasic"},
		{"nonBasic", "two", "two"},
		{"NonBasic", "two", "NonBasic"},
	}
	for ix, c := range cases {
		result := xmlTypeOf(c.t, c.nb)
		if result != c.expected {
			t.Errorf("case %d: Expected %s, got %s", ix, c.expected, result)
		}
	}
}