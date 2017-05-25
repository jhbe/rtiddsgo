# Summary

The rtiddsgo package enables golang applications to use RTI DDS. It contains
RTI DDS wrappers and tools to generate go files given DDS topics in .IDL
files.

The generated types use golang built-in types; no DDS basic types are visible. A
type may contain another DDS type, but sooner or later they all resolve
to golang built-in types.

This package was designed to be ease to use and require little type conversion
at the expense of performance; the type code does not try to be efficient when
copying between the DDS representations and the golang representation.

# Quick Start!

    go get golang.org/x/tools/cmd/goyacc
    make
    ./subscriber

...and is a separate terminal:

    ./publisher

# Build

To generate the goddsgen tool:

    go get golang.org/x/tools/cmd/goyacc
    make goddsgen

# Usage

The executable ./goddsgen is used to parse IDL files and generate golang code
for the structs, unions enums and constants found in the IDL file:

    goddsgen <idl_file_name> <rti_install_dir> <c_files_dir> <c_file_name> <out_path> <package_name>

The arguments are:
- The name of the IDL file to read (with path if necessary).
- The path to the directory containing the RTI DDS installation, such as "/opt/rti\_connext\_dds-5.2.3".
- The path to the directory containing the C files generated by rtiddsgen for the same IDL.
- The file name using in the C-files for the type.
- The path to directory in which the generated golang files should go.
- The packages name to use in all generated golang files.

The generated source code depend on the C code generated for the same IDL by
rtiddsgen.

The generated souce code have the following characteristics:

- All names of types and variables found in IDL files have a leading uppercase
character in the golang representation to make them visible from outside the
package. Most have uppercase characters starting every word.
- Names use an underscore to separate modules; "Com\_This\_That\_Blah".
- Strings are unbound. The maximum length as defined in the IDL file is NOT
enforced.
- Sequences are converted to golang slices. The maximum length as defined in
the IDL file is NOT enforced, merely noted as a comment where the variable
is defined.
- A union Foo have a variable Foo_D acting as the discriminant.

Once the IDL file has been parsed and golang source code generated, an application
uses the rtiddsgo package together with the generated source code to create
and use DDS. The rtiddsgo package contains code for:

- Participants
- Publishers
- Subscribers
- Topics

while the generated source code provide:

- Type definition
- DataWriter for the type
- DataReader for the type

Do not use the DataReader and DataWriter found in the rtiddsgo package, they
are used by the generated code, not directly by an application.

# Example

The Makefile builds an example publisher and subscriber, found in the
example and main directories.

# Limitations

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

The IDL file parsed by goddsgen may not contain comments.

The IDL file parsed by goddsgen must only contain those structs
and unions that are defined in the C file with the path and name
provided as the third and fourth argument to goddsgen. As a
consequence, all #includes must be removed. Those IDLs must be
processed separately.

Sequence length statements must be integers or constants. The IDL
specification permit math expressions such as "1 + 2", but
goddsgen do not:

    sequence<string, 2> foo;     // OK
    sequence<string, 1+2> foo;   // Not OK.


# Troubleshooting

- Got

    parse/parse.go:1: running "/home/jhbe/go/bin/goyacc": fork/exec /home/jhbe/go/bin/goyacc: no such file or directory

when building goddsgen. You forgot to:

    go get golang.org/x/tools/cmd/goyacc




# Internal Design

The three parts making up this package are:

- rtiddsgo: Static golang package for participants, publishers, subscribers and topic support.
- rtiddsgo/parse: IDL parsing.
- rtiddsgo/generate: Golang source code generation.

The parse and generate packages compile into the goddsgen tool.

## Parse

The **lexer** is responsible for converting a text into **token**s. A token
is an integer describing a keyword, such as "module" or a single character,
such as ";" or a quoted string with a predefined numerical value. Some have
additional information, such as the content of the quoted string.

The tokens are fed into the **parser**, which contains the rules for how
tokens may be ordered according to the syntax of an IDL. Those rules,
and the actions taken when a rule has been met, are defined in YACC format
in parser.y. That file is used to generate parser.go (see first line of
parse.go).

The actions in the parser start from the ground up; simple rules are combined
into longer and more complex rules, which eventually end when the IDL has
been parsed completely. Each action operate on a Node within a **tree**; the
actions build the tree from the leafs towards the trunk. The YACC format
is beyond the scope of this readme, but there is lots of material on the
internet.

The result of the parser is the IDL file as a tree of Nodes, where each node
represent a major piece of the IDL. For example, a Node may represent a "module"
and its children represents structs and unions in that module and their
children in turn represents member variables.

The tree is traversed to find all **consts**, **enums** and **structs**, which
repackages the tree structure into more compact golang structs and arrays.

## Generate

The generate package takes the enums, consts and struct definitions created in
the parse package and create golang source code. A single file with all constants
is generated by **constfile**. All enumerations and collated into a single file
by **emunsfile**.

The type definitions and associated DataWriter and DataReader source code is
generated by **structfile**. It creates one golang source code file per struct
and/or union.

All three repackages the information in such a way that it can easily be used
within a golang text/template. The information is the same and sometimes
duplicated; the purpose is to keep as little logic in the actual templates as
possible.

# Acknowledgements

The rules that make up the bulk of the file parse/parser.y file was written by
Arno Puder and downloaded from http://www.corba.org/corbadownloads.htm,
specifically http://www.omg.org/attachments/yacc.yy

