package main

import (
	"rtiddsgo/generate"
	"rtiddsgo/parse"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// This is the main entrypoint for the goddsgen tool. It reads DDS XML files
// and emits golang source code implementing the type-specific functionality
// required for a DDS type. It does for golang what rtiddsgen does for other
// languages.
//
// The command line arguments are:
// - The name of the XML file to read with path if necessary.
// - The path to the directory containing the RTI DDS installation, such as "/opt/rti_connext_dds-5.3.1".
// - The path to the directory containing the RTI DDS libraries, typically "/opt/rti_connext_dds-5.3.1/lib/x64Linux3gcc5.4.0".
// - The file name using in the C-files for the type.
// - The path to directory in which the generated golang files should go.
// - The packages name to use in all generated golang files.
//
func main() {
	if len(os.Args) != 6 {
		log.Fatal("Usage: goddsgen <xml_file_name> <rti_install_dir> <rti_lib_dir> <out_path> <package_name>")
	}

	xmlFileName := os.Args[1]
	rtiInstallDir := strings.TrimRight(os.Args[2], "/")
	rtiLibDir := strings.TrimRight(os.Args[3], "/")
	name := strings.TrimSuffix(filepath.Base(xmlFileName), ".xml")
	cFileName := name
	outPathName := os.Args[4]
	packageName := os.Args[5]


	// Read the XML files.
	xmlFile, err := os.Open(xmlFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	t, err := parse.ReadXml(xmlFile)
	if err != nil {
		log.Fatal(err)
	}
	//t.Dump(os.Stdout, 0)

	// Create the consts file. A single file holds all constants.
	constsFile, err := os.Create(outPathName + "/" + name + "_constants.go")
	if err != nil {
		log.Fatal(err)
	}
	defer constsFile.Close()

	if generate.ConstsFile(t.ModuleElements.GetConstsDef(), packageName, constsFile) != nil {
		log.Fatal(err)
	}

	// Create the enums file. A single file holds all enums.
	enumsFile, err := os.Create(outPathName + "/" + name + "_enums.go")
	if err != nil {
		log.Fatal(err)
	}
	defer enumsFile.Close()
	if generate.EnumsFile(t.ModuleElements.GetEnumsDef(), packageName, rtiInstallDir, rtiLibDir, cFileName, enumsFile) != nil {
		log.Fatal(err)
	}

	// Create the typedefs file. A single file holds all typedefs.
	typeDefsFile, err := os.Create(outPathName + "/" + name + "_typedefs.go")
	if err != nil {
		log.Fatal(err)
	}
	defer typeDefsFile.Close()
	if generate.TypeDefsFile(t.ModuleElements.GetTypeDefs(), packageName, rtiInstallDir, rtiLibDir, cFileName, typeDefsFile) != nil {
		log.Fatal(err)
	}

	// Iterate over all structs and create a type file, datawriter file, datareader file and all-in-one file for each.
	for _, sd := range t.ModuleElements.GetStructsDef() {
		fileName := outPathName + "/" + strings.ToLower(sd.GoName)

		structFile, err := os.Create(fileName + ".go")
		if err != nil {
			log.Fatal(err)
		}
		defer structFile.Close()
		if err = generate.StructFile(sd, packageName, rtiInstallDir, rtiLibDir, cFileName, structFile); err != nil {
			log.Fatal(err)
		}

		dwFile, err := os.Create(fileName + "_datawriter.go")
		if err != nil {
			log.Fatal(err)
		}
		defer structFile.Close()
		if err = generate.DataWriterFile(sd, packageName, rtiInstallDir, rtiLibDir, cFileName, dwFile); err != nil {
			log.Fatal(err)
		}

		drFile, err := os.Create(fileName + "_datareader.go")
		if err != nil {
			log.Fatal(err)
		}
		defer structFile.Close()
		if err = generate.DataReaderFile(sd, packageName, rtiInstallDir, rtiLibDir, cFileName, drFile); err != nil {
			log.Fatal(err)
		}

		allInOneFile, err := os.Create(fileName + "_allinone.go")
		if err != nil {
			log.Fatal(err)
		}
		defer structFile.Close()
		if err = generate.AllInOneFile(sd, packageName, allInOneFile); err != nil {
			log.Fatal(err)
		}
	}

	// Iterate over all structs and create a type file.
	for _, ud := range t.ModuleElements.GetUnionsDef() {
		fileName := outPathName + "/" + strings.ToLower(ud.GoName)

		unionFile, err := os.Create(fileName + ".go")
		if err != nil {
			log.Fatal(err)
		}
		defer unionFile.Close()
		if err = generate.UnionFile(ud, packageName, rtiInstallDir, rtiLibDir, cFileName, unionFile); err != nil {
			log.Fatal(err)
		}
	}
}
