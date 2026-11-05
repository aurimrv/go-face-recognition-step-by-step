package usecases

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Job is a single image analysis unit: source path and destination output path.
type Job struct {
	Src string
	Dst string
}

const (
	defaultImagesDir = "images"
	defaultUnknown   = "unknown.jpg"
	defaultResult    = "result.jpg"
)

// ResolveAnalysisJobs maps CLI path arguments to analysis jobs.
// With no args, it returns the default unknown.jpg → result.jpg job.
// File args and JPEG files under directory args become beside-source result_<stem>.jpg jobs.
// Previously generated result images are skipped. Non-JPEG files are omitted from directory
// walks; an explicit non-JPEG file still becomes a job so the batch loop can fail that job.
func ResolveAnalysisJobs(args []string) ([]Job, error) {
	if len(args) == 0 {
		src := filepath.Join(defaultImagesDir, defaultUnknown)
		dst := filepath.Join(defaultImagesDir, defaultResult)
		return []Job{{Src: src, Dst: dst}}, nil
	}

	var jobs []Job
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, fmt.Errorf("stat %q: %w", arg, err)
		}
		if info.IsDir() {
			dirJobs, err := jobsFromDir(arg)
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, dirJobs...)
			continue
		}
		if isResultOutput(filepath.Base(arg)) {
			continue
		}
		jobs = append(jobs, Job{Src: arg, Dst: resultPathFor(arg)})
	}

	if len(jobs) == 0 {
		return nil, fmt.Errorf("no analysis jobs resolved from arguments")
	}
	return jobs, nil
}

func jobsFromDir(root string) ([]Job, error) {
	var jobs []Job
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		base := d.Name()
		if isResultOutput(base) {
			return nil
		}
		if !isJPEG(path) {
			return nil
		}
		jobs = append(jobs, Job{Src: path, Dst: resultPathFor(path)})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk %q: %w", root, err)
	}
	return jobs, nil
}

func resultPathFor(src string) string {
	dir := filepath.Dir(src)
	base := filepath.Base(src)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	return filepath.Join(dir, "result_"+stem+".jpg")
}

func isJPEG(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg"
}

func isResultOutput(base string) bool {
	if base == "result.jpg" {
		return true
	}
	if !strings.HasPrefix(base, "result_") {
		return false
	}
	ext := strings.ToLower(filepath.Ext(base))
	return ext == ".jpg" || ext == ".jpeg"
}
