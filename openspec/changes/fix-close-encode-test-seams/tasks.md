## 1. Drawer seams

- [x] 1.1 Add package-level `closeFile` and `jpegEncode` seams in `internal/entity/drawer.go` (default to real `Close` / `jpeg.Encode`)
- [x] 1.2 Wire `loadImage` and `SaveImage` to use those seams (preserve existing error preference / encode-then-close behavior)

## 2. Load-persons seam

- [x] 2.1 Add package-level `closePersonsDir` seam in `internal/usecases/load_persons.go` (default to real `Close`)
- [x] 2.2 Wire the persons-directory defer close through the seam

## 3. Tests

- [x] 3.1 Add/restore drawer tests for close-after-decode failure, encode failure, and (if still needed) close-after-encode failure
- [x] 3.2 Add/restore load-persons test that forces close seam error and still returns persons

## 4. Verify

- [x] 4.1 Run `go test ./... -coverprofile=coverage.out -coverpkg=./...` and confirm total statement coverage is 100%
- [x] 4.2 Regenerate `coverage.html` and confirm the three previously blocked regions are covered
