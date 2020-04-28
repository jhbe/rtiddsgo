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

	// Create the output file. To keep the number of files down we stick all golang code in a single file.
	goFile, err := os.Create(outPathName + "/" + name + ".go")
	if err != nil {
		log.Fatal(err)
	}
	defer goFile.Close()

	unsafe := false
	for _, sd := range t.ModuleElements.GetStructsDef() {
		for _, m := range sd.Members {
			if m.GoType == "string" {
				unsafe = true
			}
		}
	}
	for _, ud := range t.ModuleElements.GetUnionsDef() {
		for _, m := range ud.Members {
			if m.GoType == "string" {
				unsafe = true
			}
		}
	}
	for _, t := range t.ModuleElements.GetTypeDefs() {
		if t.GoType == "string" {
			unsafe = true
		}
	}

	if err := generate.HeaderFile(packageName, rtiInstallDir, rtiLibDir, cFileName, unsafe, goFile); err != nil {
		log.Fatal(err)
	}
	if err := generate.ConstsFile(t.ModuleElements.GetConstsDef(), packageName, goFile); err != nil {
		log.Fatal(err)
	}
	if err := generate.EnumsFile(t.ModuleElements.GetEnumsDef(), packageName, rtiInstallDir, rtiLibDir, cFileName, goFile); err != nil {
		log.Fatal(err)
	}
	if err := generate.TypeDefsFile(t.ModuleElements.GetTypeDefs(), packageName, rtiInstallDir, rtiLibDir, cFileName, goFile); err != nil {
		log.Fatal(err)
	}

	// Iterate over all structs and create a type file, datawriter file, datareader file and all-in-one file for each.
	for _, sd := range t.ModuleElements.GetStructsDef() {
		if err = generate.StructFile(sd, packageName, rtiInstallDir, rtiLibDir, cFileName, goFile); err != nil {
			log.Fatal(err)
		}

		if !sd.Nested {
			if err = generate.DataWriterFile(sd, packageName, rtiInstallDir, rtiLibDir, cFileName, goFile); err != nil {
				log.Fatal(err)
			}

			if err = generate.DataReaderFile(sd, packageName, rtiInstallDir, rtiLibDir, cFileName, goFile); err != nil {
				log.Fatal(err)
			}

			if err = generate.AllInOneFile(sd, packageName, goFile); err != nil {
				log.Fatal(err)
			}
		}
	}

	// Iterate over all unions and create a type file.
	for _, ud := range t.ModuleElements.GetUnionsDef() {
		if err = generate.UnionFile(ud, packageName, rtiInstallDir, rtiLibDir, cFileName, goFile); err != nil {
			log.Fatal(err)
		}
	}
}
