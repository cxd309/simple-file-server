# simple-file-server

A minimal static file server for local development, designed to emulate GitHub Pages.

## Build

```bash
go build .
```

## Usage

```
./simple-file-server <directory> <port> [repo_name]
```

| Argument    | Required | Description                                                 |
| ----------- | -------- | ----------------------------------------------------------- |
| `directory` | yes      | Path to the directory to serve                              |
| `port`      | yes      | Port to listen on                                           |
| `repo_name` | no       | Repository name prefix, matching GitHub Pages URL structure |

## Examples

Serve files directly at the root:

```bash
./simple-file-server ./docs 8080
# → http://localhost:8080/
```

Emulate GitHub Pages with a repo name:

```bash
./simple-file-server ./docs 8080 my-repo
# → http://localhost:8080/my-repo/
```

This mirrors how GitHub Pages serves a repo at `https://<user>.github.io/<repo_name>/`.
