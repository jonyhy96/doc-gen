# doc-gen 

doc-gen helps you to generate apidoc like comment

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

```
$ go get github.com/jonyhy96/doc-gen
$ doc-gen -h
$ doc-gen -f examples/models/user.go

$ apidoc -f ".*\\.doc$" -o apidoc // you need have your own apidoc.json first
```

### Prerequisites

For develop

 - go ^1.10

### Installing

```
$ go build -o doc-gen .
$ sudo mv doc-gen /usr/local/bin
```

### Params

| param | required | e.g. |
| :--------: | :-----: | :----: |
| f     | true | -f examples/models/user.go |

### Coding style

[CODEFMT](https://github.com/golang/go/wiki/CodeReviewComments)

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://gitlab.domain.com/golang/containerM/tags). 

## Authors

* **HAO YUN** - *Initial work* - [haoyun](https://github.com/jonyhy96)

See also the list of [contributors](CONTRIBUTORS.md) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

* nothing