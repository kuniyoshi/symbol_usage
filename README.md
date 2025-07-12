# Symbol Usage

Symbol Usage は SCIP を使ってコードベースを読むための
プログラムです。

# Usage

```
$ sy find /tmp/index.scip Foo.bar
 Baz.main
 Baz.qux
*    Foo.bar
         Qux.sum

$ sy list /tmp/index.scip
Baz.main
Baz.qux
Foo.bar
Qux.sum

$ sy --help
Symbol Usage は SCIP を使ってコードベースを読むためのプログラムです。
```

# Specification

SCIP は Protocol Buffers でフォーマットされた indexer です。

https://4.4.sourcegraph.com/code_navigation/explanations/writing_an_indexer

