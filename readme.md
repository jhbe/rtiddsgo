# Summary

The rtiddsgo package enables golang applications to use RTI DDS. It contains
RTI DDS wrappers and tools to generate go files given DDS topics in .xml
files generated by the RTI tools from IDL files.

The generated types use golang built-in types; no DDS basic types are visible. A
type may contain another DDS type, but sooner or later they all resolve
to golang built-in types.

This package was designed to be easy to use and require little type conversion
at the expense of performance; the type code does not try to be efficient when
copying between the DDS representations and the golang representation.

# Quick Start!

Update the NDDSHOME and RTILIBDIR variables in the Makefile. Then:

    make
    ./examplesub

...and is a separate terminal:

    ./examplepub

# Build

To generate the goddsgen tool:

    make goddsgen

# Usage

Each IDL file must first be converted by the rtiddsgen tool and the -convertToXml option
to an XML represention.

The ./goddsgen executable is then used to parse the XML files and generate golang code
for the structs, unions, enums and constants found in the XML file:

    goddsgen <xml_file_name> <rti_install_dir> <rti_lib_dir> <out_path> <package_name>

The arguments are:
- The name of the XML file to read with path if necessary.
- The path to the directory containing the RTI DDS installation, such as "/opt/rti\_connext\_dds-5.3.1".
- The path to the directory containing the RTI DDS libraries, typically "/opt/rti\_connext\_dds-5.3.1/lib/x64Linux3gcc5.4.0".
- The file name using in the C-files for the type.
- The path to directory in which the generated golang files should go.
- The packages name to use in all generated golang files.

The generated source code depend on the C code generated for the same IDL by
rtiddsgen.

The generated source code have the following characteristics:

- All names of types and variables found in the IDL files have a leading uppercase
character in the golang representation to make them visible from outside the
package. Most have uppercase characters starting every word.
- Names use an underscore to separate modules; "Com\_This\_That\_Blah".
- Strings are unbound. The maximum length as defined in the IDL file is NOT
enforced.
- Sequences are converted to golang slices. The maximum length as defined in
the IDL file is NOT enforced.
- A union Foo have a variable _Discriminant acting as the discriminant.

Once the XML file has been parsed and golang source code generated, an application
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
- All-in-one participant, publisher, topic and datawriter support for the type
- All-in-one participant, subscriber, topic and datareader support for the type

Do not use the DataReader and DataWriter found in the rtiddsgo package, they
are used by the generated code, not directly by an application.

An application can use a pre-canned Participant/Publisher/Topic/DataWriter for the struct ("Com_Ex_A" is the name of the type) foud in the allinone files:

    w := example.NewWriterCom_Ex_A(33, "TheTopicName", "", "", "", "")
    defer w.Free()
    w.Dw.Write(example.Com_Ex_A{0, "Hello"})

or it can manage each separately:

    p, err := rtiddsgo.New(33, "", "")
    if err != nil {
        log.Fatal(err)
    }
    defer p.Free()

    err = example.Com_Ex_A_RegisterType(p)
    if err != nil {
        log.Fatal(err)
    }

    topic, err := p.CreateTopic("MyMessage", example.Com_Ex_A_GetTypeName(), "", "")
    if err != nil {
        log.Fatal(err)
    }
    defer topic.Free()

    pub, err := p.CreatePublisher("", "")
    if err != nil {
        log.Fatal(err)
    }
    defer pub.Free()

    dw, err := example.NewCom_Ex_ADataWriter(pub, topic, "", "")
    if err != nil {
        log.Fatal(err)
    }
    defer dw.Free()

    dw.Write(example.Com_Ex_A{0, "Hello"})

The same is true for subscribers.

Types are mapped as follows:

IDL            |   XML   |  C                | Go
---------------|---------|-------------------|--------------
boolean        | boolean | DDS_Boolean       | bool
short          | int16   | DDS_Short         | int16
unsigned short | uint16  | DDS_UnsignedShort | uint16
long           | int32   | DDS_Long          | int32
unsigned long  | uint32  | DDS_UnsignedLong  | uint32
float          | float32 | DDS_Float         | float32
double         | float64 | DDS_Double        | float64
string         | string  | DDS_Char *        | string

# Example

The Makefile builds an example publisher and subscriber, found in the
example directory.

# Limitations

The "//@key" and "//@top-level" IDL directives are NOT used by rtiddsgo.
All structs have datawriter and datareader support.

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

Sequences must be bounded.

Sequences cannot have the index based on an enum.

A union discriminant cannot be a typedef type.

# Internal Design

The three parts making up this package are:

Package           | Description
------------------|-----------------------------------------------------------------
rtiddsgo          | Static golang package for participants, publishers, subscribers and topic support.
rtiddsgo/parse    | Parses the XML file and extracts constans, enums, structs, unions and typedefs.
rtiddsgo/generate | Golang source code generation based on the result of the rtiddsgo/parse package.
rtiddsgo/main     | Main for the goddsgen tool itself.

## Parse

File              | Description
------------------|-----------------------------------------------------------------
constants         | Creates a collection of all constants found in the IDL file.
enums             | Creates a collection of all enums found in the IDL file.
structs           | Creates a collection of all structs found in the IDL file.
typedefs          | Creates a collection of all typedefs found in the IDL file.
unions            | Creates a collection of all unions found in the IDL file.
traverse          | Provides functions for traversing the content of the parsed XML file.
utils             | Various functions for converting names to other formats
xml               | Parses the XML file.

## Generate

File              | Description
------------------|-----------------------------------------------------------------
allinone          | Generates the all-in-one go source code for a DDS type.
consts            | Generates the go source code for a DDS constants.
datareader        | Generates the DDS data reader go source code for a DDS type.
datawriter        | Generates the DDS data writer go source code for a DDS type.
enums             | Generates the go source code for a DDS enums.
structs           | Generates the go source code for a DDS struct type.
typedefs          | Generates the go source code for a DDS typedef type.
unions            | Generates the go source code for a DDS union type.
flags             | Defines the CGO flags.
retreive          | Provides a function for generating the DDS type Retrieve go source code.
store             | Provides a function for generating the DDS type Store go source code.

# TODO

- Function comments.
- Register and dispose support.
- Why doesn't verificationsub always stop properly?
- 

