package snippet

import "math/rand"

// Subset
// https://sre.google/sre-book/load-balancing-datacenter/
func Subset(backends []string, clientID, subsetSize int) []string {

	subsetCount := len(backends) / subsetSize

	// Group clients into rounds; each round uses the same shuffled list:
	round := clientID / subsetCount

	r := rand.New(rand.NewSource(int64(round)))
	r.Shuffle(len(backends), func(i, j int) { backends[i], backends[j] = backends[j], backends[i] })

	// The subset id corresponding to the current client:
	subsetID := clientID % subsetCount

	start := subsetID * subsetSize
	return backends[start : start+subsetSize]
}

/*
为什么是均匀的
- shuffle 算法保证在 round 一致的情况下，backend 的排列一定是一致的
- client 上下线
  - client 上下线用滚动更新的方式，并不会影响其它 client 的连接分布，所以每个 client 下线时，只是对应的后端少了一些连接，暂时会导致某些 backend 的连接比其它 backend 少 1。
  - 上线 client 从尾部开始，client id 依然是递增的，按照该算法，这些 client 会继续排在其它 client 后面，一个 round 一个 round 地将连接分布在后端服务上，也必然是均匀的。
- server 上下线
  - 与 client 上下线类似，server 的滚动升级和上下线也是不会有大影响的，因为每个 server 会随机地分布在不同 client 的子集中，不会因为该 server 上下线，导致计算结果有大变化。

算法的问题
- 在 client 或者 server 端的实例数量发生变化的时候，会导致大量的连接迁移和重建
- 每个服务都能被分配从 0 到 N 的连续唯一 id，这一点在没有外部依赖的情况下比较难做到。绑定了外部基础设施的方案又可能比较难推广。比如 k8s 的 statefulset，也没办法强制所有服务都使用
- 服务下线时，并不一定能保证下线的服务的 client id 是连续的，这样就总是可以构造出一些极端情况，在拿到一些 client 之后，让某台 backend 的连接数变为 0。
- 现在大规模的服务节点很多，有些批量发布一次性发布几百个节点，Google 的这个算法说一般 100 条连接(We typically use a subset size of 20 to 100 backend tasks)就够了？如果正好批量发布的后端都被同一个 client 选中了，那这个 client 就废掉了。
- client 服务是需要知道 backends 的 id 的，否则当 backend 发生下线时，会导致 client 端的连接重新排部。

Reinventing Backend Subsetting at Google  https://queue.acm.org/detail.cfm?id=3570937
- consistent subsetting
  -算法并不会有较好的连接平衡度(server 端连接分布不均匀)和连接区分度(可能会把连接的实例分配给同一个 client)
- Ringsteady subsetting
  -后端节点的分布要尽可能让后端变化时，前端和后端建立的连接也成相同比例变化，所以 Google 工程师选用了一个特殊的分布序列：binary van der Corput sequence
  -算法能保证较好的 subset diversity，subset spread，并且能均匀地将后端和前端节点分布在圆环上
  - 问题
   - ringsteady subsetting 的连接平衡度较差
   - 前文的 binary van der Corput sequence 导致 client 和 server 在环上离得较近，而并不是按照距离来均匀分布
- 灵活组合 Ringsteady 来解决所有问题的算法
  - 将前后端所有实例进行分组
  - 然后将 server 端的分组内的实例进行 shuffle，这里需要保证每个 client 组内看到的 shuffle 结果是相同的，所以可以用 frontend 的 lot id 来作为伪随机算法的输入种子
  -然后将 server 的 LOT 也就是组按照 ringsteady 的排序方式进行打乱

*/
// binary van der Corput sequence
double corput(int n, int base) {
	double q=0, bk=(double)1/base;

	while (n > 0) {
		q += (n % base)*bk;
		n /= base;
		bk /= base;
	}

	return q;
}
