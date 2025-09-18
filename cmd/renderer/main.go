package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"docgen-service/internal/docgen"
)

func main() {
	// Define command-line flags
	var (
		shellPath      = flag.String("shell", "", "Path to the shell DOCX file")
		componentsDir  = flag.String("components", "", "Directory containing component XML files")
		planPath       = flag.String("plan", "", "Path to the JSON plan file")
		outputPath     = flag.String("output", "", "Path where the generated DOCX should be saved")
	)
	flag.Parse()

	// Validate required arguments
	if *shellPath == "" || *componentsDir == "" || *planPath == "" || *outputPath == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -shell <path> -components <dir> -plan <path> -output <path>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Printf("Starting DocGen CLI renderer...")
	log.Printf("Shell: %s", *shellPath)
	log.Printf("Components: %s", *componentsDir)
	log.Printf("Plan: %s", *planPath)
	log.Printf("Output: %s", *outputPath)

	// Initialize the engine
	log.Printf("Initializing DocGen engine...")
	engine, err := docgen.NewEngine(*shellPath, *componentsDir)
	if err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	log.Printf("Loaded components: %v", engine.GetLoadedComponents())

	// Read and parse the plan file
	log.Printf("Reading plan file...")
	planFile, err := os.Open(*planPath)
	if err != nil {
		log.Fatalf("Failed to open plan file: %v", err)
	}
	defer planFile.Close()

	planData, err := io.ReadAll(planFile)
	if err != nil {
		log.Fatalf("Failed to read plan file: %v", err)
	}

	var plan docgen.DocumentPlan
	if err := json.Unmarshal(planData, &plan); err != nil {
		log.Fatalf("Failed to parse plan JSON: %v", err)
	}

	log.Printf("Plan loaded: %d components to render", len(plan.Body))

	// Generate the document
	log.Printf("Assembling document...")
	result, err := engine.Assemble(plan)
	if err != nil {
		log.Fatalf("Failed to assemble document: %v", err)
	}

	log.Printf("Document assembled successfully, size: %d bytes", len(result))

	// Write the result to the output file
	log.Printf("Writing output file...")
	if err := os.WriteFile(*outputPath, result, 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	log.Printf("Document generated successfully: %s", *outputPath)
}