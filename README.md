# graphPartition

implementation of graph partition algorithms in bytedance summer camp 2020

### todo-list
- [x] 1.` fanout 优化 ` 在当前 bucket 内部进行sort, 选择gain最小的k个进行bucket更换。
- [x] 2.` 时间优化 ` 金涛发现不是每一个vertex的gain都需要重新计算。
    *因为每次修改单个vertex的bucket后，只会对vertex的一阶邻居(nbr_i)的bucket统计有影响,

    因此，会对vertex的二阶邻居gain有影响。所以在多次迭代之后，需要修改的vertex会很少，二阶邻居很少的情况下，
    
    可以只对此次迭代进行bucket修改的vertex的二阶邻居重新计算gain。

- [ ] 3. `时间优化` 在1的基础上，如果不选取严格前k小的vertex。可以并行在x个segment中每个选择，k/x个vertex。

- [ ] 4. `时间优化` 当前是按照vertex做的线程数据分块，如果按照边数做线程数据分块，每个线程的任务将会更加平均。
