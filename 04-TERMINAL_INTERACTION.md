```bash
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ git branch
  00_baseline
  01_static_analysis_skill
  02_explore_testing_tools
  03_black_box_tests
* 04_white_box_tests
  main
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ git diff --stat 03_black_box_tests -- ’*.go’ ’:status_test.go’
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ git diff *.go
diff --git a/cmd/main.go b/cmd/main.go
index 53c5589..8642f93 100644
--- a/cmd/main.go
+++ b/cmd/main.go
@@ -16,24 +16,29 @@ var (
        modelsDir  = "models"
        personsDir = "persons"
        fontPath   = filepath.Join("fonts", "Roboto-Regular.ttf")
+       exitFunc   = os.Exit
 )
 
 // This example shows the basic usage of the package: create an
 // recognizer, recognize faces, classify them using few known ones.
 func main() {
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ 
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ git diff --stat 02_explore_testing_tools -- '*.go' ':!*_test.go'
 cmd/main.go                       | 24 +++++++++++++++++-------
 internal/entity/drawer.go         | 15 ++++++++++++---
 internal/usecases/load_persons.go |  7 ++++++-
 3 files changed, 35 insertions(+), 11 deletions(-)
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ git restore --source=02_explore_testing_tools -- '*.go' ':!*_test.go'
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ git diff --stat 02_explore_testing_tools -- '*.go' ':!*_test.go'
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go test ./...
# github.com/andre-ols/go-face-recognition/cmd [github.com/andre-ols/go-face-recognition/cmd.test]
cmd/main_test.go:37:10: undefined: run
cmd/main_test.go:49:10: undefined: run
cmd/main_test.go:61:10: undefined: run
cmd/main_test.go:73:10: undefined: run
cmd/main_test.go:87:10: undefined: run
cmd/main_test.go:110:10: undefined: run
cmd/main_test.go:119:22: undefined: exitFunc
cmd/main_test.go:121:2: undefined: exitFunc
cmd/main_test.go:124:3: undefined: exitFunc
cmd/main_test.go:148:10: undefined: run
cmd/main_test.go:148:10: too many errors
FAIL    github.com/andre-ols/go-face-recognition/cmd [build failed]
# github.com/andre-ols/go-face-recognition/internal/entity [github.com/andre-ols/go-face-recognition/internal/entity.test]
internal/entity/drawer_test.go:128:9: undefined: closeFile
internal/entity/drawer_test.go:129:2: undefined: closeFile
internal/entity/drawer_test.go:133:21: undefined: closeFile
internal/entity/drawer_test.go:148:9: undefined: jpegEncode
internal/entity/drawer_test.go:149:2: undefined: jpegEncode
internal/entity/drawer_test.go:152:21: undefined: jpegEncode
internal/entity/drawer_test.go:166:9: undefined: closeFile
internal/entity/drawer_test.go:167:2: undefined: closeFile
internal/entity/drawer_test.go:171:21: undefined: closeFile
FAIL    github.com/andre-ols/go-face-recognition/internal/entity [build failed]
# github.com/andre-ols/go-face-recognition/internal/usecases [github.com/andre-ols/go-face-recognition/internal/usecases.test]
internal/usecases/load_persons_test.go:90:9: undefined: closePersonsDir
internal/usecases/load_persons_test.go:91:2: undefined: closePersonsDir
internal/usecases/load_persons_test.go:95:21: undefined: closePersonsDir
FAIL    github.com/andre-ols/go-face-recognition/internal/usecases [build failed]
FAIL
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-ste
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go test ./...
ok      github.com/andre-ols/go-face-recognition/cmd19.132s
ok      github.com/andre-ols/go-face-recognition/internal/entity      0.036s
ok      github.com/andre-ols/go-face-recognition/internal/usecases    7.405s
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ STEP=04-3
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go test -coverprofile=${STEP}-coverage.out ./...
ok      github.com/andre-ols/go-face-recognition/cmd18.745s       coverage: 100.0% of statements
ok      github.com/andre-ols/go-face-recognition/internal/entity      (cached)        coverage: 97.1% of statements
ok      github.com/andre-ols/go-face-recognition/internal/usecases    7.417s  coverage: 99.1% of statements
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go tool cover -func=${STEP}-coverage.out &> ${STEP}-coverage.txt
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go tool cover -html=${STEP}-coverage.out -o ${STEP}-covera
ge.html
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ STEP=04-4
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go test -coverprofile=${STEP}-coverage.out ./...
ok      github.com/andre-ols/go-face-recognition/cmd   19.527s  coverage: 100.0% of statements
ok      github.com/andre-ols/go-face-recognition/internal/entity        (cached)        coverage: 100.0% of statements
ok      github.com/andre-ols/go-face-recognition/internal/usecases      (cached)        coverage: 100.0% of statements
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go tool cover -func=${STEP}-coverage.out &> ${
STEP}-coverage.txt
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go tool cover -html=${STEP}-coverage.out -o ${STEP}-coverage.html
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ 
```