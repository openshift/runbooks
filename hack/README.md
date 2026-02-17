# Linting Markdown

Running a linter container

```sh
# build container
podman build -t linter -f hack/Dockerfile.lint

# run lint container
podman run -it --rm -v ${PWD}:/src:z linter
```

Running `lint.sh`

```sh
hack/lint.sh

# lint spelling
hack/lint_spelling.sh

# lint markdown
hack/lint_spelling.sh
```

Fix markdown errors

```sh
. venv/bin/activate
pymarkdown fix --recurse .
```
