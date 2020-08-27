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

