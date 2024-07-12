package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/aliml92/go-typesense/typesense"
	"github.com/jackc/pgx/v5/pgxpool"
	jsoniter "github.com/json-iterator/go"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/aliml92/anor/postgres/repository/category"
	"github.com/aliml92/anor/postgres/repository/product"
	"github.com/aliml92/anor/postgres/repository/store"
	"github.com/aliml92/anor/postgres/repository/user"
	ts "github.com/aliml92/anor/typesense"
)

var cfg config

type ProductJSON struct {
	Categories []string            `json:"categories"`
	Name       string              `json:"name"`
	Price      string              `json:"price"`
	Available  string              `json:"available"`
	Sold       string              `json:"sold"`
	ImageUrls  []string            `json:"image_links"`
	Specs      map[string]string   `json:"specs"`
	Attributes map[string][]string `json:"attributes"`

	Brand          string `json:"-"`
	RootCategoryID int32  `json:"-"`
}

func init() {
	cfg = config{}
	flag.StringVar(&cfg.datasetDir, "dataset-dir", "", "path to the dataset folder")
	flag.StringVar(&cfg.database, "database", "", "PostgresSQL database connection string")
	flag.StringVar(&cfg.typesense, "typesense", "", "Typesense server url")
	flag.StringVar(&cfg.imported, "import", "./cmd/ebay-data-importer/imported_files.txt",
		"A file that contains a list of import files into database and typesense server")
	flag.BoolVar(&cfg.ignoreImported, "ignore-imported", false, "import all json files forcefully")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if cfg.datasetDir == "" {
		return fmt.Errorf("dataset-dir path is required")
	}

	if cfg.database == "" {
		return fmt.Errorf("PostgresSQL database connection string is required")
	}

	if cfg.typesense == "" {
		return fmt.Errorf("typesense server url string is required")
	}

	db, err := pgxpool.New(ctx, cfg.database)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}
	defer db.Close()

	tsClient, _ := typesense.NewClient(nil, cfg.typesense, "xyz")

	userRepository := user.NewRepository(db)
	productRepository := product.NewRepository(db)
	categoryRepository := category.NewRepository(db)
	storeRepository := store.NewRepository(db)
	searcher := ts.NewSearcher(tsClient)

	d := &DataImporter{
		cfg:       cfg,
		database:  db,
		typesense: tsClient,

		userRepository:     userRepository,
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
		storeRepository:    storeRepository,
		searcher:           searcher,
	}

	datasetList, err := getNewJSONLFiles(d.cfg.datasetDir, d.cfg.imported, d.cfg.ignoreImported)
	if err != nil {
		return fmt.Errorf("failed to get new jsonl files: %w", err)
	}
	if len(datasetList) == 0 {
		return fmt.Errorf("no new jsonl files found")
	}

	fmt.Printf("dataset list: %v\n", datasetList)
	d.datasetList = datasetList

	err = d.safeImportData(ctx)
	if err != nil {
		return err
	}

	return nil
}

func getNewJSONLFiles(datasetDir string, importedFilePath string, ignoreImported bool) ([]string, error) {
	var newJSONLFiles []string
	importedFiles := make(map[string]bool)

	// Read the imported files if not ignoring
	if !ignoreImported {
		file, err := os.Open(importedFilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			importedFiles[filepath.Base(scanner.Text())] = true
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	// Read the directory
	entries, err := os.ReadDir(datasetDir)
	if err != nil {
		return nil, err
	}

	// Iterate over the entries
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".jsonl") {
			if !ignoreImported && importedFiles[entry.Name()] {
				continue
			}
			newJSONLFiles = append(newJSONLFiles, filepath.Join(datasetDir, entry.Name()))
		}
	}

	return newJSONLFiles, nil
}

func parse(filepath string) ([]ProductJSON, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	var products []ProductJSON
	reader := bufio.NewReaderSize(f, 4*1024*1024) // 4MB buffer
	json := jsoniter.ConfigFastest

	lineNumber := 0
	totalSize := 0
	skippedLines := 0

	for {
		lineNumber++
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading line %d: %w", lineNumber, err)
		}

		totalSize += len(line)
		if len(strings.TrimSpace(string(line))) == 0 {
			if err == io.EOF {
				break
			}
			continue
		}

		var productJSON ProductJSON
		if err := json.Unmarshal(line, &productJSON); err != nil {
			// Skip lines that can't be parsed as JSON
			skippedLines++
			if lineNumber%1000 == 0 {
				fmt.Printf("Skipped %d lines so far due to parsing errors\n", skippedLines)
			}
			if err == io.EOF {
				break
			}
			continue
		}

		products = append(products, productJSON)

		if lineNumber%10000 == 0 {
			fmt.Printf("Processed %d lines (%.2f MB), skipped %d lines\n", lineNumber, float64(totalSize)/1024/1024, skippedLines)
		}

		if err == io.EOF {
			break
		}
	}

	return products, nil
}

func updateImported(fileName string, newEntries []string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range newEntries {
		if _, err := writer.WriteString(entry + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %v", err)
	}

	return nil
}
