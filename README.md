[![Build Status](https://travis-ci.com/robinmitra/forest.svg?branch=master)](https://travis-ci.com/robinmitra/forest)

# forest
For the forest on your computer

## Features

### Analyse files

The `analyse` command analyses files and directories at a given path, and summarises the following
metrics:
* Total number of files and directories
* Total disk space usage
* Top 5 file types (by occurrence and disk usage)
* Ability to ignore certain files and/or directories

#### General usage

```bash
# Analyse the forest at the given path.
forest analyse .
```

#### Options

* `--include-dot-files`: Include hidden dot files in the analysis. These are excluded by default.
* `--format`: The output format of the summary. Options include `normal` (default) and `rainbow`.

## TODO

In no particular order:

- [ ] Dockerise
- [ ] Output formats (e.g. rainbow)
- [x] Installable via Homebrew
- [x] Verbose and version flags
