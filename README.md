# SemVer Tagger

A simple CLI utility to create [Semantic Version](https://semver.org) tags in Git repos.

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
  - `-a` - auto-detect the version type based on commit messages since the last version was tagger (default behaviour).
  - `-rc` - make this a release candidate version.
- Other options
  - `-v` - verbose logging.
  - `-f`, `--force` - ignore the following conditions that would normally cause the tool to abort:
    - The repo work tree is not clean.
    - The repo is not on a default branch.
    - There have been no commits since the last version was tagged.
  - `-y`, `--no-confirm` - don't wait for confirmation on the new tag.
  - `-P`, `--push` - push tags to the remote after creating them.

### Verison Notes

- **This does not support the full semver spec!** This is a simple, one-size-fits-most tool, not a universal tool that will work on every project.
- Major, minor and patch bumps will reset the RC counter, unless the `-rc` flag is passed.
  - See the tests file for more details on the expected behaviour.
