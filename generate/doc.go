/*
Package generate creates golang source code files implementing the
type-specific functionality given the IDL file parsed in package parse.

Three types of files are created. The "enum.go" contains all enumerations.
The "constants.go" file contains all constants. The third type of file
represent one struct or union and is named after the type-name of the
IDL struct/union. There are as many such files as there were structs and
unions in the IDL file.

The files are created with the functions CreateEnumsFile(), CreateConstsFile()
and CreateStructFile().
 */
package generate
