## Getting started

```console
$ go install github.com/koron/funddb
```

```console
# Initiation
$ funddb database initschema
$ funddb fund import list.tsv

$ funddb price fetchlatest
```

list.tsv has fund information.
The format is

```
{Association ID}\t{Fund Name}\t{Fund URL}[\t{Fetch ID}]
```

The format of Fetch ID is `{scheme}:{id}`
Currently `{scheme}` support `ammufg` and `fidelity` only.
