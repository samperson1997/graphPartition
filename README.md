# graphPartition

### todo-list
- [ ] 1.` fanout 优化 ` 在当前 bucket 内部进行sort, 选择gain最小的k个进行bucket更换。
- [ ] 2.` 时间优化 ` 金涛发现不是每一个vertex的gain都需要重新计算。
    *因为每次修改单个vertex的bucket后，只会对vertex的一阶邻居(nbr_i)的bucket统计有影响,

    因此，会对vertex的二阶邻居gain有影响。所以在多次迭代之后，需要修改的vertex会很少，二阶邻居很少的情况下，
    
    可以只对此次迭代进行bucket修改的vertex的二阶邻居重新计算gain。

- [ ] 3. `时间优化` 在1的基础上，如果不选取严格前k小的vertex。可以并行在x个segment中每个选择，k/x个vertex。

- [ ] 4. `时间优化` 当前是按照vertex做的线程数据分块，如果按照边数做线程数据分块，每个线程的任务将会更加平均。


算法主要分为三步

1.ComputMoveGain

* 需要收集每个vertex的邻居的bucket信息。
* 可以进行vertex级别的并行

2.ComputMoveProb
* 为了保证每个节点中的vertex数量差距不大，计算move概率
* 

3.SetNew
* 进行每个vertex的bucket更新
* 可以进行vertex级别的并行


#### 并没有什么用的优化idea

* 由算法可知，需要三步才能够进行bucket更新，其实是异步的，然而如果可以在步骤1的过程中确定bucket是否需要更新，那么步骤三其实是可以省略一部分的。
比如说在步骤1的过程中，发现gain值比较合适并且迁移过程中可以保证分块中vertex数量平衡，那么直接在步骤1中进行bucket更新。

* 第一步计算vertex的moveGain，可以选择不计算全部vertex，在一定进行计算。

* 如果进行随机选取vertex计算，是否可以给与每一个顶点不同的选择概率，比如当前iteration中**bucket不需要进行替换**的vertex，在下一个iteration中被选择的概率变低。

* 每一步计算需要替换bucket的次数应该是随着iteration的增长逐渐减少的。是否能够根据iteration需要替换的bucket个数做不同的计算方法。

### 单机多线程版本结构
Vertex
* (Net)告知别的节点自己属于哪个bucket
* (Local)计算对每个bucket的move gain
* (Net)发送自己的target给master
* (Net)接受master发送回来的转移概率

Master
* (Net)接受节点发送来的target
* (Local)汇总成矩阵S
* (Local)计算每个节点的转移概率
* (Net)发送转移概率给每个节点

Bucket
* (Net)进行节点转移
