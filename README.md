### QuickSearch

#### Getting Started

TODO

#### Run or build from source

To run the `quicksearch` from source, you should clone the repo firstly.

```sh
git clone git@github.com:feimingxliu/quicksearch.git
# or use 'git clone https://github.com/feimingxliu/quicksearch.git' if you don't set SSH key.
```

The you should download the dependencies.

```sh
cd quicksearch && go mod tidy -compat=1.17 # go version >= 1.17
```

You can just run the following command to start the `quicksearch`.

```sh
go run github.com/feimingxliu/quicksearch/cmd/quicksearch -c configs/config.yaml
```

Or you can run build the binary like this:

```sh
go build -o bin/quicksearch github.com/feimingxliu/quicksearch/cmd/quicksearch
```

Run binary:

```sh
bin/quicksearch -c configs/config.yaml
```

#### Tests

If you want to run the tests, it will use some testdata which stores with [git-lfs](https://git-lfs.github.com/). After
you install the [git-lfs](https://git-lfs.github.com/), you can run

```sh
git lfs pull
```

in the project root to fetch the large test file.

Then you should run

```sh
go run github.com/feimingxliu/quicksearch/cmd/xmltojson
```

The above command will generate `test/testdata/zhwiki-20220601-abstract.json`, you can open it to see the content.

In the end, you can run all the tests by:

```sh
go test -timeout 0 ./...
```

If everything works well, you will see a `ok` in the end of output.