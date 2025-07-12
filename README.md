# Symbol Usage

Symbol Usage は SCIP を使ってコードベースを読むための
プログラムです。

# Usage

## Basic Usage

```bash
# Find symbol usage by short name
$ sy find index.scip fetchCmd
*    fukumimi/cmd/fetchCmd
*    fukumimi/cmd/fetchCmd:cmd

$ sy find index.scip args
*    fukumimi/cmd/fetchCmd:args
*    fukumimi/cmd/mergeCmd:args
*    fukumimi/cmd/versionCmd:args

$ sy find index.scip Execute
*    fukumimi/cmd/Execute
*    cobra/Command#Execute

# List all symbols (now with simplified names)
$ sy list index.scip | grep fukumimi/cmd
fukumimi/cmd/Execute
fukumimi/cmd/fetchCmd
fukumimi/cmd/fetchCmd:args
fukumimi/cmd/fetchCmd:cmd
fukumimi/cmd/fetchCmd:error
fukumimi/cmd/init
fukumimi/cmd/mergeCmd
```

## Symbol Name Format

シンボル名は短い名前で指定できます（大幅に改善されました）：

1. **単純名** (推奨):
   - `args` - 変数名で検索
   - `fetchCmd` - 関数名で検索 
   - `Execute` - メソッド名で検索
   - `Command` - 型名で検索

2. **パッケージ付き**:
   - `cmd/fetchCmd` - パッケージ名付き
   - `cobra/Command` - 外部パッケージの型

3. **詳細指定**:
   - `fetchCmd:args` - 特定の関数の引数
   - `fetchCmd:cmd` - 特定の関数内の変数

**改善点**:
- 長いSCIP名 (`github.com/kuniyoshi/fukumimi.v0.3.0...`) を短縮
- シンボル名の部分マッチングを改善
- より直感的な検索が可能

## Verbose Mode

実際のSCIPシンボル名を確認したい場合は `-v` フラグを使用：

```bash
# List with SCIP symbol names
$ sy list index.scip -v | head -5
github.com/kuniyoshi/fukumimi/cmd/Execute          => scip-go gomod github.com/kuniyoshi/fukumimi v0.3.0 `github.com/kuniyoshi/fukumimi/cmd`/Execute().
github.com/kuniyoshi/fukumimi/cmd/fetchCmd         => scip-go gomod github.com/kuniyoshi/fukumimi v0.3.0 `github.com/kuniyoshi/fukumimi/cmd`/fetchCmd.
github.com/kuniyoshi/fukumimi/cmd/fetchCmd:args    => scip-go gomod github.com/kuniyoshi/fukumimi v0.3.0 `github.com/kuniyoshi/fukumimi/cmd`/fetchCmd:args.

# Find with SCIP symbol names
$ sy find index.scip Execute -v
Searching for symbol: Execute
(Will also match SCIP patterns containing this symbol)

*    github.com/kuniyoshi/fukumimi/cmd/Execute          [scip-go gomod github.com/kuniyoshi/fukumimi v0.3.0 `github.com/kuniyoshi/fukumimi/cmd`/Execute().]
```

# Specification

SCIP は Protocol Buffers でフォーマットされた indexer です。

https://4.4.sourcegraph.com/code_navigation/explanations/writing_an_indexer

