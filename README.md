# delta

A Difference Algorithm based on "An O(ND)Difference Algorithm and Its Variations" by Eugene W. Myers.

## Install

`$ go install github.com/stuartmscott/delta/cmd/delta`

## Usage

```
$ echo "Hello World" > file1
$ echo "Hi Earth" > file2
$ delta file1 file2
0   72  H
1 - 101 e
2 - 108 l
3 - 108 l
4 - 111 o
1 + 105 i
2   32   
3 - 87  W
4 - 111 o
3 + 69  E
4 + 97  a
5   114 r
6 - 108 l
7 - 100 d
6 + 116 t
7 + 104 h
```
