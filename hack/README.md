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
```
