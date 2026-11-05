```bash
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go test ./... -cover
        github.com/andre-ols/go-face-recognition/cmd   coverage: 0.0% of statements
        github.com/andre-ols/go-face-recognition/internal/entity                coverage: 0.0% of statements
ok      github.com/andre-ols/go-face-recognition/internal/usecases      0.016s  coverage: 41.6% of statements
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ STEP=03
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go test -coverprofile=${STEP}-coverage.out ./...
        github.com/andre-ols/go-face-recognition/cmd   coverage: 0.0% of statements
        github.com/andre-ols/go-face-recognition/internal/entity                coverage: 0.0% of statements
ok      github.com/andre-ols/go-face-recognition/internal/usecases      0.015s  coverage: 41.6% of statements
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go tool cover -func=${STEP}-coverage.out &> ${STEP}-coverage.txt
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go tool cover -html=${STEP}-coverage.out -o ${STEP}-coverage.html
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ 
```