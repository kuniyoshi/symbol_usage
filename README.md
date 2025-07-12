# Symbol Usage

Symbol Usage は SCIP を使ってコードベースを読むための
プログラムです。

# Usage

```
$ sy /tmp/index.scip Foo.bar
 Baz.main
 Baz.qux
*    Foo.bar
         Qux.sum
$
```

# Specification

SCIP は Protocol Buffers でフォーマットされた indexer です。

https://4.4.sourcegraph.com/code_navigation/explanations/writing_an_indexer

# TODO

1. SCIP を読み込む
2. シンボルをリストする
3. パラメーターでシンボルを受け取る。シンボルはユーザーが指定しやすい形式を優先し、SCIP の完全なフォーマット指定ではない
4. 受け取ったシンボルを SCIP 用のフォーマットに変換する
5. シンボルの caller/caller を列挙する
