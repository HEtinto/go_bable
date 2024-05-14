package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type DsNode struct {
	dsId uint32
}

func (ds *DsNode) DsGetSegments() (sids []uint64, err error) {
	// 暂时删除全部
	// segment ID生成算法: DS ID 左移 32位
	var (
		dsCapacity uint64 = 257
	)
	if dsCapacity < 256 {
		sids = append(sids, uint64((*ds).dsId)<<32)
	} else if dsCapacity%256 != 0 {
		sids = append(sids, uint64((*ds).dsId)<<32+dsCapacity/256)
	}
	for sid, i := uint64((*ds).dsId)<<32, uint64(0); i < dsCapacity/256; i++ {
		sids = append(sids, sid+i)
	}
	return sids, err
}

func main() {
	// 1.创建客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("create etcd client failed:%v.\n", err)
		return
	}
	defer cli.Close()
	// 2.创建操作上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3.插入键值对
	_, err = cli.Put(ctx, "my-key", "my-value")
	if err != nil {
		// handle error!
		fmt.Printf("insert key-value failed:%v.\n", err)
		return
	}
	_, err = cli.Put(ctx, "my-key-1", "my-value-1")
	if err != nil {
		// handle error!
		fmt.Printf("insert key-value failed:%v.\n", err)
		return
	}

	// 4.获取键值对-指定键
	rsp, err := cli.Get(ctx, "my-key")
	if err != nil {
		// handle error!
		fmt.Printf("get key-value failed:%v.\n", err)
		return
	}

	// 4.1输出查询结果
	fmt.Printf("Key value pair:\n")
	fmt.Printf("指定键查询原始结果:\n%v\n", *rsp)
	// {cluster_id:14841639068965178418 member_id:10276657743932975437 revision:6 raft_term:5  [key:"my-key" create_revision:2 mod_revision:6 version:5 value:"my-value" ] false 1 {} [] 0}
	for _, ev := range rsp.Kvs {
		fmt.Printf("[key]:%s\n[value]:%s\n", ev.Key, ev.Value)
	}

	// 5.获取键值对-模糊查询-指定前缀
	rsp, err = cli.Get(ctx, "my-key", clientv3.WithPrefix())
	if err != nil {
		// handle error!
		fmt.Printf("get key-value failed:%v.\n", err)
		return
	}

	// 5.1输出查询结果
	fmt.Printf("Key value pair:\n")
	fmt.Printf("模糊查询原始结果:\n%v\n", *rsp)
	fmt.Printf("模糊查询结果拆分:\n")
	for _, ev := range rsp.Kvs {
		fmt.Printf("[key]:%s\n[value]:%s\n", ev.Key, ev.Value)
	}

	// 6.删除键值对-指定值
	fmt.Printf("Try to delete:%s\n", "my-key")
	rspDelete, err := cli.Delete(ctx, "my-key")
	if err != nil {
		// handle error!
		fmt.Printf("delete key-value failed:%v.\n", err)
		return
	}
	fmt.Printf("删除指定键原始结果:\n%v.\n", rspDelete)

	// 7.批量操作
	// 7.1定义事务操作
	ops := []clientv3.Op{
		clientv3.OpPut("my-key-2", "my-value-2"),
		clientv3.OpPut("my-key-3", "my-value-3"),
		clientv3.OpPut("my-key-4", "my-value-4"),
		clientv3.OpPut("my-key-5", "my-value-5"),
		clientv3.OpPut("my-key-6", "my-value-6"),
		clientv3.OpPut("my-key-7", "my-value-7"),
		clientv3.OpPut("my-key-8", "my-value-8"),
		clientv3.OpPut("my-key-9", "my-value-9"),
		clientv3.OpGet("my-key-2"),
		clientv3.OpDelete("my-key-1"),
		// 不能在一次事务中提交对相同键值的删除操作,例如提交A和删除A在同一次事务中
		// batch operation failed:etcdserver: duplicate key given in txn request.
	}
	// 7.2提交事务
	txnRsp, err := cli.Txn(context.Background()).Then(ops...).Commit()
	if err != nil {
		// handle error!
		fmt.Printf("batch operation failed:%v.\n", err)
		return
	}
	// 7.3检查事务是否成功
	if txnRsp.Succeeded {
		fmt.Printf("batch operation succeeded.\n")
	} else {
		fmt.Printf("batch operation failed.\n")
	}
	// 7.4打印响应详情
	for i, rspTxn := range txnRsp.Responses {
		fmt.Printf("第%v次操作,结果:%v\n", i+1, rspTxn)
	}
}
