# SemVer Tagger

A simple CLI utility to create [Semantic Version](https://semver.org) tags in Git repos.

Note that _this does not support the full semver spec_. This is a simple, one-size-fits-most tool, not a universal tool that will work on every project.

## Installation

```shell
go install github.com/markormesher/semver-tagger
```

### Usage

```shell
semver-tagger [options]
```

Options:

- Version type
  - `-M` - bump the major vesion.
  - `-m` - bump the minor vesion.
  - `-p` - bump the patch vesion.
  - `-rc` - add or increment the release candidate counter.
  - `-no-rc` - remove the release candidate counter.
  - `-init` - create the initial `v0.1.0` version.
- Other options
  - `-v` - verbose logging.
  - `-f`, `--force` - ignore the following conditions that would normally cause the tool to abort:
    - The repo work tree is not clean.
    - The repo is not on a default branch.
    - There have been no commits since the last version was tagged.
  - `-y`, `--no-confirm` - don't wait for confirmation on the new tag.
  - `-P`, `--push` - push tags to the remote after creating them.

### Examples

> For more examples, check the test cases in the code.

```
# Bump the major version
Command: semver-tagger -M
Examples:
  v1.2.3      ->  v2.0.0
  v1.2.3-rc2  ->  v2.0.0  (note that the RC counter is reset)

# Bump the major version and mark it as a release candidate
Command: semver-tagger -M -rc
Examples:
  v1.2.3      ->  v2.0.0-rc1
  v1.2.3-rc2  ->  v2.0.0-rc1  (note that the RC counter is reset)

# Bump the minor version
Command: semver-tagger -m
Examples:
  v1.2.3      ->  v1.3.0
  v1.2.3-rc2  ->  v1.3.0  (note that the RC counter is reset)

# Bump the minor version and mark it as a release candidate
Command: semver-tagger -m -rc
Examples:
  v1.2.3      ->  v1.3.0-rc1
  v1.2.3-rc2  ->  v1.3.0-rc1  (note that the RC counter is reset)

# Bump the patch version
Command: semver-tagger -p
Examples:
  v1.2.3      ->  v1.2.4
  v1.2.3-rc2  ->  v1.2.4  (note that the RC counter is reset)

# Bump the patch version and mark it as a release candidate
Command: semver-tagger -p -rc
Examples:
  v1.2.3      ->  v1.2.4-rc1
  v1.2.3-rc2  ->  v1.2.4-rc1  (note that the RC counter is reset)

# Bump the release candidate counter only
Command: semver-tagger -rc
Examples:
  v1.2.3      ->  v1.2.3-rc1  (note that you shouldn't really be tagging RC versions when a non-RC already exists)
  v1.2.3-rc2  ->  v1.2.4-rc3

# Remove the release candidate counter
Command: semver-tagger -no-rc
Examples:
  v1.2.3      ->  v1.2.3  (no change)
  v1.2.3-rc2  ->  v1.2.3
```
