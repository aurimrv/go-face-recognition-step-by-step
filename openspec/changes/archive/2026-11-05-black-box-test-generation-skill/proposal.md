## Why

The teaching fork needs black-box Go tests driven by equivalence partitioning (EP) and boundary value analysis (BVA), with explicit partition documentation in the suite. A project skill standardizes how agents derive partitions, emit annotated tests, and—when recognition behavior is under test—create fixture face/unknown images that exercise different match strengths.

## What Changes

- Add project skill `.cursor/skills/black-box-test-generation/SKILL.md` that instructs the agent to generate Go tests using EP and BVA.
- Require the generated test set to document partition definitions in comments and to label each test case with the partition or boundary it covers.
- Document when and how to create face images and unknown images with similar faces so classification/recognition tests can cover different degrees of matching (strong match, weak/near-threshold, non-match).
- Scope the skill to black-box behavior of package APIs (inputs → observable outputs/errors); it does not redesign the recognition pipeline, models, or Docker native dependencies.

## Capabilities

### New Capabilities

- `black-box-test-generation`: Agent workflow and policies for EP/BVA-based Go test generation, including partition/boundary comments and face-image fixture guidance for this project.

### Modified Capabilities

- (none)

## Impact

- New files under `.cursor/skills/black-box-test-generation/`.
- No application runtime changes as part of *this* change; later skill runs may add `*_test.go` files and test fixtures under e.g. `testdata/` or `images/` test trees.
- Complements existing coverage/mutation tooling (`go test -cover`, `gomutants`) by improving the *quality* of the suite design, not replacing those tools.
