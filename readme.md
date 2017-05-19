# Summary

The rtiddsgo package enables golang applications to use RTI DDS. It contains
RTI DDS wrappers and tools to generate go files given DDS topics in .IDL
files.

The generated types use golang built-in types; no DDS basic types are visible. A
type may contain another golang type, but sooner or later they all resolve
to golang built-in types.

# Quick Start!

    go get golang.org/x/tools/cmd/goyacc
    make
    ./subscriber

...and is a separate terminal:

    ./publisher

# Build

    go get golang.org/x/tools/cmd/goyacc
    make

# Usage

The directory rtiddsgo/eb contains all generated code.

All names of types and variables found in IDL files have a leading uppercase
character in the golang representation to make them visible from outside the
package. Most have uppercase characters starting every word.

Names use an underscore to separate modules; "Com_This_That_Blah".

Strings are unbound. The maximum length as defined in the IDL file is NOT
enforced.

Sequences are converted to golang slices. The maximum length as defined in
the IDL file is NOT enforced, merely noted as a comment where the variable
is defined.

A union Foo have a variable Foo_D acting as the discriminant.

The "//@key" and "//@top-level" IDL directives are NOT used by rtiddsgo.
All structs and types have datawriter and datareader support.

A type is expected to be either fully qualified ("com::this::that") or be
defined in the same module it is used. Example:

    module A {
      struct B {
        string C;
      };
      
      struct D {
        B myB;  // This is OK.
      };
      
      module E {
        struct F {
          B myB;  // This is NOT ok. Must be A::B, not just B.
        };
      };
    };



# Internal Design



# Acknowledgements

The rules that make up the bulk of the file parse/parser.y file was written by
Arno Puder and downloaded from http://www.corba.org/corbadownloads.htm,
specifically http://www.omg.org/attachments/yacc.yy
