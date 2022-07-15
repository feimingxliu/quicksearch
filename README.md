# QuickSearch

### Introduction

Quicksearch is a lightweight search engine which deploys and runs as a single binary. It's inspired
by [Zinc](https://github.com/zinclabs/zinc) but it uses [Bleve](https://github.com/blevesearch/bleve) as underlying
indexing library and supports chinese by default. It just supports local storage now.

### Getting Started

If you have installed [golang](https://golang.google.cn/learn/), just run following:

 ```bash
// if go version >= 1.16
go install github.com/feimingxliu/quicksearch/cmd/quicksearch@latest
// or `go get -u github.com/feimingxliu/quicksearch/cmd/quicksearch` for go < 1.16
 ```

This will generate binary in your `$GOPATH/bin`. Or you can get the prebuilt binary
from [Releases](https://github.com/feimingxliu/quicksearch/releases).

To run the quicksearch. Copy the example config.

```bash
wget https://raw.githubusercontent.com/feimingxliu/quicksearch/master/configs/config.yaml
// run quicksearch
quicksearch -c config.yaml
```

### API Reference

#### Index API

+ *Create Index*

```
POST /<index>
{
    "settings": <Index Settings>,
    "mappings": <Index Mappings>
}
```

The request body can be ignored which use default.
`<Index Settings>` is an object which contains index's setting

```
{
    "number_of_shards": int 
}
```

`<Index Mappings>` is an object which defines index's mapping

```

{
	"types": {
		"Document Type": <Documnet Mapping>, 
		....
		"Document Type": <Documnet Mapping>
	}
	"default_mapping": <Documnet Mapping>,
	"type_field": string,
	"default_type": string,
	"default_analyzer": string 
}

```

`<Documnet Mapping>`  is an object which defines document's mapping

```

{
	"disabled": bool,	# disable this documnet mapping
	"properties": {
			"name": <Documnet Mapping>,	# this enables nested json
			...
			"name": <Documnet Mapping>
		},
	"fields": <Field Mapping>,
	"default_analyzer": string 
}
```

`<Field Mapping>` is an array which defines field level mapping

```

[
	# you can define one more field mapping for one field, that's why its array {
	"type": string,	# support "keyword", "text", "datetime", "number", "boolean", "geopoint", "IP"
	"analyzer": string,	# specifies the name of the analyzer to use for this field
	"store": bool,	# indicates whether to store field values in the index
	"index": bool	# indicates whether to analyze the field }
]

```

+ *Update Index Mapping*

```
PUT /<index>/_mapping
<Index Mapping>
```

+ *Get Index Detail*

```
GET /<index>
```

+ *Open Index*

```
POST /<index>/_open
```

+ *Close Index*

```
POST /<index>/_close
```

+ *Clone Index*

```
POST /<index>/_clone/<cloned index>
```

+ *List Indices*

```
GET /_all
```

+ *Delete Index*

```
DELETE /<index>
```

#### Document API

+ *Index Document*

```
POST /<index>/_doc
<document json object>
# or with custom documnet id
POST /<index>/_doc/<docID>
<document json object>
```

  If index a document with same docID, the newer one will cover old fully.

+ *Bulk*

```
POST /_bulk 
POST /<index>/_bulk
<Action Line>
<optional document json object>
......
<Action Line>
<optional document json object>
```

`<Action Line>` is an object defines which operation to execute.

```
{ 
	# <Action> can be `create`, `delete`, `index`, `update`
	<Action>: {
		"_index": string,
		"_id": string 
	} 
}
```

+ *Update Document*

```
PUT /<index>/_doc/<docID>
{
	"fieldName": any 
}
```

  This can update part fields of document.

+ *Get Document*

```
GET /<index>/_doc/<docID>
```

+ *Delete Document*

```
DELETE /<index>/_doc/<docID>
```

#### Search API

```
POST /<index>/_search
GET /<index>/_search
POST /_search 
GET /_search 
{ 
	"query": <Query>,
	"size": int,
	"from": int,
	"highlight": bool,
	"fields": []string,
	"facets": {
		<facet name>: {
            "size": int,
            "field": string,
            "numeric_ranges": [
                {
                    "name": string,
                    "min": float64,
                    "max": float64
                }
            ],
            "date_ranges": [
                {
                    "name": string,
                    "start": datetime, # RFC3339
                    "end": datetime # RFC3339 
                }
            ]
         } 
     },
     "explain": bool,
     "sort": []sting,
     "includeLocations": bool,
     "search_after": []sting,
     "search_before": []string 
}
```

`<Query>` indicates different query, see [Queries](http://blevesearch.com/docs/Query/).

+ *Query String Query*

```
{
	"query": string,
	"boost": int 
}
```

This is the simplest query for search,
see  [full query language specification](http://blevesearch.com/docs/Query-String-Query/).

### Run or build from source

To run the `quicksearch` from source, clone the repo firstly.

```sh
git clone git@github.com:feimingxliu/quicksearch.git
# or use 'git clone https://github.com/feimingxliu/quicksearch.git' if you don't set SSH key.
```

Then download the dependencies.

```sh
cd quicksearch && go mod tidy -compat=1.17 # go version >= 1.17
```

Run the following command to start the `quicksearch`.

```sh
go run github.com/feimingxliu/quicksearch/cmd/quicksearch -c configs/config.yaml
```

Or build the binary like this:

```sh
go build -o bin/quicksearch github.com/feimingxliu/quicksearch/cmd/quicksearch
```

Run binary:

```sh
bin/quicksearch -c configs/config.yaml
```

### Tests

The tests use some testdata which stores with [git-lfs](https://git-lfs.github.com/). After you have installed
the [git-lfs](https://git-lfs.github.com/), you can run

```sh
git lfs pull
```

in the project root to fetch the large test file.

Then run

```sh
go run github.com/feimingxliu/quicksearch/cmd/xmltojson
```

The above command will generate `test/testdata/zhwiki-20220601-abstract.json`, you can open it to see the content.

In the end, just run all the tests by

```sh
go test -timeout 0 ./...
```

If everything works well, an `ok` will appear at the end of output.