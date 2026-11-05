package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Kagami/go-face"
	"github.com/andre-ols/go-face-recognition/internal/entity"
	"github.com/andre-ols/go-face-recognition/internal/usecases"
)

var (
	modelsDir  = "models"
	personsDir = "persons"
	fontPath   = filepath.Join("fonts", "Roboto-Regular.ttf")
)

// This example shows the basic usage of the package: create an
// recognizer, recognize faces, classify them using few known ones.
func main() {

	initialTime := time.Now()

	defer func() {
		fmt.Printf("\x1b[34mTotal time: %s\x1b[0m\n", time.Since(initialTime))
	}()

	pathArgs := os.Args[1:]
	failFast := len(pathArgs) == 0

	jobs, err := usecases.ResolveAnalysisJobs(pathArgs)
	if err != nil {
		log.Fatalf("Can't resolve analysis jobs: %v", err)
	}

	loadPersonsUseCase := usecases.NewLoadPersonsUseCase()

	// Load persons from fs.
	persons, err := loadPersonsUseCase.Execute(personsDir)
	if err != nil {
		log.Fatalf("Can't load persons: %v", err)
	}

	initRecognizerTime := time.Now()

	// Init the recognizer.
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()

	fmt.Println("Time to init recognizer: ", time.Since(initRecognizerTime))

	// Recognize faces from a few images and assign them to persons.
	recognizePersonsUseCase := usecases.NewRecognizePersonsUseCase(rec)

	err = recognizePersonsUseCase.Execute(persons)

	if err != nil {
		log.Fatalf("Can't recognize persons: %v", err)
	}

	classifyPersonUseCase := usecases.NewClassifyPersonsUseCase(rec)

	var failed bool
	for _, job := range jobs {
		if err := processJob(job, persons, classifyPersonUseCase); err != nil {
			if failFast {
				log.Fatalf("%v", err)
			}
			log.Printf("%v", err)
			failed = true
		}
	}

	if failed {
		os.Exit(1)
	}
}

func processJob(
	job usecases.Job,
	persons map[int]*entity.Person,
	classifyPersonUseCase usecases.ClassifyPersonsUseCase,
) error {
	recognizedFaces, err := classifyPersonUseCase.Execute(job.Src, 0.3)
	if err != nil {
		return fmt.Errorf("Can't classify %s: %w", job.Src, err)
	}

	fmt.Printf("\033[0;32m%s: found %d faces\033[0m\n", job.Src, len(recognizedFaces))
	for _, recognizedFace := range recognizedFaces {
		fmt.Printf("\033[0;32mPerson: %s\033[0m\n", persons[recognizedFace.ID].Name)
	}

	drawer, err := entity.NewDrawer(job.Src, fontPath)
	if err != nil {
		return fmt.Errorf("Can't create drawer for %s: %w", job.Src, err)
	}

	for _, recognizedFace := range recognizedFaces {
		drawer.DrawFace(recognizedFace.Face.Rectangle, persons[recognizedFace.ID].Name)
	}

	if err := drawer.SaveImage(job.Dst); err != nil {
		return fmt.Errorf("Can't save %s: %w", job.Dst, err)
	}

	fmt.Printf("\033[0;32mWrote %s\033[0m\n", job.Dst)
	return nil
}
