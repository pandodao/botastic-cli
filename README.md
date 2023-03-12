# botastic-cli

A command line interface for [Botastic](https://developers.pando.im/references/botastic/api.html).

## Installation

```bash
go install github.com/pandodao/botastic-cli@latest
```

## Set Environment Variables

```bash
export BOTASTIC_APP_ID=YOUR_APP_ID
export BOTASTIC_SECRET=YOUR_APP_SECRET
```

## Usage

### Build Indices from Markdown files

```bash
botastic-cli scan --dir ./docs --type md --mode paragraph
```

in which,

- `--dir` is the directory of your markdown files
- `--type` is the type of your markdown files, currently only `md` is supported
- `--mode` is the mode of building indices, currently supports `paragraph` and `line`

### Create Indices 

```bash
botastic-cli index --act create --file ./indices.json
```

### Query

```bash
botastic-cli index --act query --query "hello world"
```

