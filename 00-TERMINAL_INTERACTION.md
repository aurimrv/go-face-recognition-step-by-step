# Terminal interaction: markdownlint sample output

Verbatim capture of a class terminal session (lint report dump).
Long lines inside the fenced transcript are intentional.

```bash
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ markdownlint-cli2 "**/*.md" "#openspec" "#.cursor" "#node_modules"
markdownlint-cli2 v0.24.0 (markdownlint v0.42.0)
Finding: **/*.md !openspec !.cursor !node_modules
Linting: 2 files
Summary: 84 issues in 2 files
DEVCONTAINER_SETUP.md:10:1 error MD060/table-column-style Table column style [Table pipe is missing space to the right for style "compact"]
DEVCONTAINER_SETUP.md:10:5 error MD060/table-column-style Table column style [Table pipe is missing space to the left for style "compact"]
DEVCONTAINER_SETUP.md:10:5 error MD060/table-column-style Table column style [Table pipe is missing space to the right for style "compact"]
DEVCONTAINER_SETUP.md:10:9 error MD060/table-column-style Table column style [Table pipe is missing space to the left for style "compact"]
DEVCONTAINER_SETUP.md:10:9 error MD060/table-column-style Table column style [Table pipe is missing space to the right for style "compact"]
DEVCONTAINER_SETUP.md:10:13 error MD060/table-column-style Table column style [Table pipe is missing space to the left for style "compact"]
DEVCONTAINER_SETUP.md:13:81 error MD013/line-length Line length [Expected: 80; Actual: 101]
DEVCONTAINER_SETUP.md:16:81 error MD013/line-length Line length [Expected: 80; Actual: 81]
DEVCONTAINER_SETUP.md:31:81 error MD013/line-length Line length [Expected: 80; Actual: 89]
DEVCONTAINER_SETUP.md:41:81 error MD013/line-length Line length [Expected: 80; Actual: 89]
DEVCONTAINER_SETUP.md:44 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```bash"]
DEVCONTAINER_SETUP.md:46 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```"]
DEVCONTAINER_SETUP.md:58:81 error MD013/line-length Line length [Expected: 80; Actual: 87]
DEVCONTAINER_SETUP.md:61 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```ini"]
DEVCONTAINER_SETUP.md:65 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```"]
DEVCONTAINER_SETUP.md:69:81 error MD013/line-length Line length [Expected: 80; Actual: 92]
DEVCONTAINER_SETUP.md:80 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```bash"]
DEVCONTAINER_SETUP.md:83 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```"]
DEVCONTAINER_SETUP.md:85:81 error MD013/line-length Line length [Expected: 80; Actual: 87]
DEVCONTAINER_SETUP.md:92:81 error MD013/line-length Line length [Expected: 80; Actual: 87]
DEVCONTAINER_SETUP.md:108 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```bash"]
DEVCONTAINER_SETUP.md:111 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```"]
DEVCONTAINER_SETUP.md:115:81 error MD013/line-length Line length [Expected: 80; Actual: 94]
DEVCONTAINER_SETUP.md:117 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```bash"]
DEVCONTAINER_SETUP.md:124:81 error MD013/line-length Line length [Expected: 80; Actual: 88]
DEVCONTAINER_SETUP.md:125:81 error MD013/line-length Line length [Expected: 80; Actual: 102]
DEVCONTAINER_SETUP.md:129:81 error MD013/line-length Line length [Expected: 80; Actual: 166]
DEVCONTAINER_SETUP.md:134:81 error MD013/line-length Line length [Expected: 80; Actual: 93]
DEVCONTAINER_SETUP.md:153:81 error MD013/line-length Line length [Expected: 80; Actual: 101]
DEVCONTAINER_SETUP.md:158:81 error MD013/line-length Line length [Expected: 80; Actual: 124]
DEVCONTAINER_SETUP.md:159:81 error MD013/line-length Line length [Expected: 80; Actual: 100]
DEVCONTAINER_SETUP.md:160:81 error MD013/line-length Line length [Expected: 80; Actual: 100]
DEVCONTAINER_SETUP.md:161:81 error MD013/line-length Line length [Expected: 80; Actual: 89]
DEVCONTAINER_SETUP.md:163:81 error MD013/line-length Line length [Expected: 80; Actual: 100]
DEVCONTAINER_SETUP.md:169:81 error MD013/line-length Line length [Expected: 80; Actual: 108]
DEVCONTAINER_SETUP.md:174:81 error MD013/line-length Line length [Expected: 80; Actual: 86]
DEVCONTAINER_SETUP.md:180:81 error MD013/line-length Line length [Expected: 80; Actual: 91]
DEVCONTAINER_SETUP.md:183:81 error MD013/line-length Line length [Expected: 80; Actual: 111]
DEVCONTAINER_SETUP.md:184:81 error MD013/line-length Line length [Expected: 80; Actual: 100]
DEVCONTAINER_SETUP.md:185:81 error MD013/line-length Line length [Expected: 80; Actual: 93]
DEVCONTAINER_SETUP.md:186:81 error MD013/line-length Line length [Expected: 80; Actual: 97]
DEVCONTAINER_SETUP.md:194:1 error MD060/table-column-style Table column style [Table pipe is missing space to the right for style "compact"]
DEVCONTAINER_SETUP.md:194:5 error MD060/table-column-style Table column style [Table pipe is missing space to the left for style "compact"]
DEVCONTAINER_SETUP.md:194:5 error MD060/table-column-style Table column style [Table pipe is missing space to the right for style "compact"]
DEVCONTAINER_SETUP.md:194:9 error MD060/table-column-style Table column style [Table pipe is missing space to the left for style "compact"]
DEVCONTAINER_SETUP.md:195:81 error MD013/line-length Line length [Expected: 80; Actual: 219]
DEVCONTAINER_SETUP.md:196:81 error MD013/line-length Line length [Expected: 80; Actual: 409]
DEVCONTAINER_SETUP.md:197:81 error MD013/line-length Line length [Expected: 80; Actual: 170]
DEVCONTAINER_SETUP.md:198:81 error MD013/line-length Line length [Expected: 80; Actual: 203]
DEVCONTAINER_SETUP.md:199:81 error MD013/line-length Line length [Expected: 80; Actual: 613]
DEVCONTAINER_SETUP.md:200:81 error MD013/line-length Line length [Expected: 80; Actual: 336]
DEVCONTAINER_SETUP.md:201:81 error MD013/line-length Line length [Expected: 80; Actual: 473]
DEVCONTAINER_SETUP.md:202:81 error MD013/line-length Line length [Expected: 80; Actual: 170]
DEVCONTAINER_SETUP.md:203:81 error MD013/line-length Line length [Expected: 80; Actual: 135]
DEVCONTAINER_SETUP.md:204:81 error MD013/line-length Line length [Expected: 80; Actual: 162]
DEVCONTAINER_SETUP.md:205:81 error MD013/line-length Line length [Expected: 80; Actual: 139]
DEVCONTAINER_SETUP.md:217 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```json"]
DEVCONTAINER_SETUP.md:219 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```"]
DEVCONTAINER_SETUP.md:220:81 error MD013/line-length Line length [Expected: 80; Actual: 89]
DEVCONTAINER_SETUP.md:223 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```bash"]
DEVCONTAINER_SETUP.md:228 error MD031/blanks-around-fences Fenced code blocks should be surrounded by blank lines [Context: "```"]
DEVCONTAINER_SETUP.md:231:81 error MD013/line-length Line length [Expected: 80; Actual: 89]
DEVCONTAINER_SETUP.md:234:81 error MD013/line-length Line length [Expected: 80; Actual: 97]
DEVCONTAINER_SETUP.md:236:81 error MD013/line-length Line length [Expected: 80; Actual: 139]
README.md:3:81 error MD013/line-length Line length [Expected: 80; Actual: 459]
README.md:7:81 error MD013/line-length Line length [Expected: 80; Actual: 285]
README.md:11:3 error MD051/link-fragments Link fragments should be valid [Context: "[Overview](#overview)"]
README.md:27:81 error MD013/line-length Line length [Expected: 80; Actual: 768]
README.md:31:81 error MD013/line-length Line length [Expected: 80; Actual: 326]
README.md:35:81 error MD013/line-length Line length [Expected: 80; Actual: 134]
README.md:37:81 error MD013/line-length Line length [Expected: 80; Actual: 161]
README.md:39:81 error MD013/line-length Line length [Expected: 80; Actual: 221]
README.md:45:81 error MD013/line-length Line length [Expected: 80; Actual: 382]
README.md:49:81 error MD013/line-length Line length [Expected: 80; Actual: 310]
README.md:53:81 error MD013/line-length Line length [Expected: 80; Actual: 190]
README.md:57:81 error MD013/line-length Line length [Expected: 80; Actual: 191]
README.md:67:1 error MD029/ol-prefix Ordered list item prefix [Expected: 1; Actual: 2; Style: 1/1/1]
README.md:73:1 error MD029/ol-prefix Ordered list item prefix [Expected: 1; Actual: 3; Style: 1/1/1]
README.md:79:1 error MD029/ol-prefix Ordered list item prefix [Expected: 1; Actual: 4; Style: 1/1/1]
README.md:85:81 error MD013/line-length Line length [Expected: 80; Actual: 93]
README.md:89 error MD040/fenced-code-language Fenced code blocks should have a language specified [Context: "```"]
README.md:152:81 error MD013/line-length Line length [Expected: 80; Actual: 814]
README.md:156:81 error MD013/line-length Line length [Expected: 80; Actual: 158]
README.md:160:81 error MD013/line-length Line length [Expected: 80; Actual: 92]
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go vet ./...
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ gofmt -l .
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go test ./... -cover
        github.com/andre-ols/go-face-recognition/cmd       coverage: 0.0% of statements
        github.com/andre-ols/go-face-recognition/internal/entity           coverage: 0.0% of statements
        github.com/andre-ols/go-face-recognition/internal/usecases         coverage: 0.0% of statements
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ openspec init --tools cursor
✔ Setup complete for Cursor

OpenSpec Setup Complete

Refreshed: Cursor
4 skills and 4 commands in .cursor/
Config: skipped (non-interactive mode)

Getting started:
  Start your first change: /opsx:propose "your idea"

Learn more: https://github.com/Fission-AI/OpenSpec
Feedback:   https://github.com/Fission-AI/OpenSpec/issues

Restart your IDE for slash commands to take effect.

vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go run ./cmd
Time to load persons:  2.301098ms
Time to init recognizer:  960.429122ms
Time to recognize known faces:  1.362008171s
Time to recognize unknown faces:  202.536522ms
images/unknown.jpg: found 1 faces
Person: Denzel-Washington
Wrote images/result.jpg
Total time: 2.615837894s
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go run ./cmd images/
Denzel-Washington/ Tom-Cruise/
Kate-Winslet/      result.jpg
Megan-Fox/         unknown.jpg
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go run ./cmd images/Megan-Fox/unknown-01.jpg images/Tom-Cruise/unknown-01.jpg 
Time to load persons:  2.123359ms
Time to init recognizer:  976.703082ms
Time to recognize known faces:  1.384318419s
Time to recognize unknown faces:  612.357445ms
images/Megan-Fox/unknown-01.jpg: found 1 faces
Person: Megan-Fox
Wrote images/Megan-Fox/result_unknown-01.jpg
Time to recognize unknown faces:  952.953695ms
images/Tom-Cruise/unknown-01.jpg: found 1 faces
Person: Tom-Cruise
Wrote images/Tom-Cruise/result_unknown-01.jpg
Total time: 4.11317073s
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ go run ./cmd images/
Time to load persons:  1.1371ms
Time to init recognizer:  939.288993ms
Time to recognize known faces:  1.411861788s
Time to recognize unknown faces:  568.956646ms
images/Denzel-Washington/unknown-01.jpg: found 1 faces
Person: Denzel-Washington
Wrote images/Denzel-Washington/result_unknown-01.jpg
Time to recognize unknown faces:  449.117609ms
images/Denzel-Washington/unknown-02.jpg: found 1 faces
Person: Denzel-Washington
Wrote images/Denzel-Washington/result_unknown-02.jpg
Time to recognize unknown faces:  198.26307ms
images/Denzel-Washington/unknown-03.jpg: found 1 faces
Person: Denzel-Washington
Wrote images/Denzel-Washington/result_unknown-03.jpg
Time to recognize unknown faces:  402.594063ms
images/Kate-Winslet/unknown-01.jpg: found 1 faces
Person: Kate-Winslet
Wrote images/Kate-Winslet/result_unknown-01.jpg
Time to recognize unknown faces:  789.35611ms
images/Kate-Winslet/unknown-02.jpg: found 0 faces
Wrote images/Kate-Winslet/result_unknown-02.jpg
Time to recognize unknown faces:  589.544085ms
images/Megan-Fox/unknown-01.jpg: found 1 faces
Person: Megan-Fox
Wrote images/Megan-Fox/result_unknown-01.jpg
Time to recognize unknown faces:  924.154992ms
images/Tom-Cruise/unknown-01.jpg: found 1 faces
Person: Tom-Cruise
Wrote images/Tom-Cruise/result_unknown-01.jpg
Time to recognize unknown faces:  1.53913229s
images/Tom-Cruise/unknown-02.jpg: found 0 faces
Wrote images/Tom-Cruise/result_unknown-02.jpg
Time to recognize unknown faces:  197.0545ms
images/unknown.jpg: found 1 faces
Person: Denzel-Washington
Wrote images/result_unknown.jpg
Total time: 9.31258472s
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ git status --short
 M README.md
 M cmd/main.go
 M internal/entity/drawer.go
?? .cursor/
?? images/Denzel-Washington/result_unknown-01.jpg
?? images/Denzel-Washington/result_unknown-02.jpg
?? images/Denzel-Washington/result_unknown-03.jpg
?? images/Kate-Winslet/result_unknown-01.jpg
?? images/Kate-Winslet/result_unknown-02.jpg
?? images/Megan-Fox/result_unknown-01.jpg
?? images/Tom-Cruise/result_unknown-01.jpg
?? images/Tom-Cruise/result_unknown-02.jpg
?? images/result.jpg
?? images/result_unknown.jpg
?? internal/usecases/resolve_analysis_jobs.go
?? internal/usecases/resolve_analysis_jobs_test.go
?? openspec/
vscode@cbcc797d2928:/workspaces/go-face-recognition-step-by-step$ git add -A;git commit --allow-empty -m "step 00: baseline";git switch -c 01_static_analysis_skill
```
