# pso-optimizer 样例（example/）

`pso-optimizer` 是纯计算工具，不读取输入文件，靠命令行 flag 驱动；给定相同 `-seed` 结果可复现。
本目录提供一份可复现的样例调用及其期望输出（`expected.txt`，由 `seed=7` 固定生成）。

## 运行

```bash
go run . -problem sphere -dims 5 -iters 100 -swarm 20 -seed 7
```

## 期望输出（见 expected.txt）

```
problem:   sphere (dims=5)
best fit:  7.427837148588364e-08
best pos:  [ 0.0002,  0.0001, -0.0001, -0.0001, -0.0001]
pareto front size: 1
front (objective vectors):
  [ 0.0000]
```

> 说明：浮点数值在不同平台/Go 版本下末位可能略有差异，`best fit` 接近 `0`（sphere 问题全局最优在原点）即视为正确。
> 其它可选问题：`-problem zdt1`（多目标，会输出 pareto front）、`-problem constrained`（带约束二次问题）。
