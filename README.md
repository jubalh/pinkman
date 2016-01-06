# pinkman
**the totally kafkaesque chess game**

pinkman is a command line based chess interface to UCI compatible chess engines, particularly Stockfish; using Unicode written in Go.

# Installation

For 64 bit Linux you can find the pinkman and stockfish binaries zipped into `pinkman-linux-64.tar.gz` in [releases](https://github.com/jubalh/pinkman/releases).
If you want to compile them yourself or use another operating system you need to:

```
go get github.com/jubalh/pinkman
```

Download [stockfish](https://stockfishchess.org/download/) or install it via your package manager.
Either place the stockfish binary in the same directory as the pinkman binary or specify a path to the stockfish engine:

```
$GOPATH/bin/pinkman --path /usr/bin/stockfish
```

Of course if you have added `$GOPATH/bin` to your `PATH` variable you can skip that prefix.

Use the `--help` option to learn about pinkmans arguments and use the in game `help` command to learn about its commands.
