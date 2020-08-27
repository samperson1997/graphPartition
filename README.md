# graphPartition



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

由算法可知，需要三步才能够进行bucket更新，其实是异步的，然而如果可以在步骤1的过程中确定bucket是否需要更新那么步骤三其实是可以省略一部分的。
比如说在步骤1的过程中，发现gain值比较合适并且迁移过程中可以保证分块中vertex数量平衡，那么直接在步骤1中进行bucket更新。

### todo-list
- [ ] 实现串行算法
- [ ] 实现单机线程级并行算法
- [ ] 算法优化
- [ ] 多机并行，vertex数据迁移