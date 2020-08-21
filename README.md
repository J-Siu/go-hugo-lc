## go-hugo-lc

Hugo site link checker written in golang. Handy for checking internal link breakage after migration.

<!-- TOC -->

- [Highlight](#highlight)
- [Limitation](#limitation)
- [Install](#install)
- [Usage](#usage)
- [Repository](#repository)
- [Contributors](#contributors)
- [Change Log](#change-log)
- [License](#license)

<!-- /TOC -->

### Highlight

- Does not use markdown parser
- Golang standard library only

### Limitation

- Only detect md links in format `[]()`
- Verify internal link only

### Install

```sh
go get github.com/J-Siu/go-hugo-lc
cd $GOPATH/src/github.com/J-Siu/go-hugo-lc
go install
```

### Usage

```sh
go-hugo-lc
```

```sh
go-hugo-lc MIT License  Copyright (c) 2020 John Siu
Support: https://github.com/J-Siu/go-hugo-lc/issues
Debug  : export _DEBUG=true
Usage  : go-hugo-lc <baseURL> <content> <public>
```

In Hugo site root:

```sh
go-hugo-lc https://example.com content public
```

Single page usage:

```sh
go-hugo-lc https://example.com content/post/post.md public
```

### Repository

- [go-hugo-lc](https://github.com/J-Siu/go-hugo-lc)

### Contributors

- [John Sing Dao Siu](https://github.com/J-Siu)

### Change Log

- 0.5
  - Initial Commit

### License

The MIT License

Copyright (c) 2020 John Siu

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.