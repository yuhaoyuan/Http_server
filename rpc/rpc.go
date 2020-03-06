package rpc

import (
	"github.com/yuhaoyuan/Http_server/config"
	"github.com/yuhaoyuan/RPC_server/corn"
	"log"
	"net"

	"sync"
	//"reflect"
)

// specialRpcClient rpc客户端
var specialRPCClient *corn.Client

// Mut rpcClient的锁
var Mut sync.Mutex

// GetSingleton 获得rpcClient
func GetSingleton() *corn.Client {
	if specialRPCClient == nil {
		SpecialRPClientInit()
	}
	connIsOk := specialRPCClient.CheckConn()
	if !connIsOk {
		SpecialRPClientInit()
	}
	return specialRPCClient
}

// GetNewClient 获得rpcClient
func GetNewClient() *corn.Client {
	var specialRPCClient *corn.Client
	conn, err := net.Dial("tcp", config.BaseConf.Addr)
	if err != nil {
		log.Printf("client-dial failed!, err = ", err)
	}
	log.Println("----------------make RpcClient------------------")
	specialRPCClient = corn.NewClient(conn)
	return specialRPCClient
}

// SpecialRPClientInit RpcClient 构造方法
func SpecialRPClientInit() {
	conn, err := net.Dial("tcp", config.BaseConf.Addr)
	if err != nil {
		log.Printf("client-dial failed!, err = ", err)
	}
	log.Println("----------------make RpcClient------------------")
	specialRPCClient = corn.NewClient(conn)
}

//// 下面是客户端连接池的时候，  用来优化单例
//var RpcClientPool Pool
//
//func InitRpc(){
//	rpcFactory := func() (*corn.Client, error) {
//		conn, err := net.Dial("tcp", config.BaseConf.Addr)
//		if err != nil {
//			log.Printf("client-dial failed!, err = ", err)
//		}
//		rpcClient := corn.NewClient(conn)
//		return rpcClient, nil
//	}
//
//	//close 关闭连接的方法
//	rpcClose :=  func(v *corn.Client) error {
//		v.Close()
//		return nil
//	}
//
//	//创建一个连接池： 初始化2，最大连接10，空闲连接数是4
//	poolConfig := &rpcConfig{
//		InitialCap: 2,
//		MaxIdle:    4,
//		MaxCap:     10,
//		Factory:    rpcFactory,
//		Close:      rpcClose,
//		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
//		IdleTimeout: 15 * time.Second,
//	}
//	RpcClientPool, _ = NewChannelPool(poolConfig)
//
//}
//
//// rpcConfig 连接池相关配置
//type rpcConfig struct {
//	//连接池中拥有的最小连接数
//	InitialCap int
//	//最大并发存活连接数
//	MaxCap int
//	//最大空闲连接
//	MaxIdle int
//	//生成连接的方法
//	Factory func() (*corn.Client, error)
//	//关闭连接的方法
//	Close func(*corn.Client) error
//	//检查连接是否有效的方法
//	Ping func(*corn.Client) error
//	//连接最大空闲时间，超过该事件则将失效
//	IdleTimeout time.Duration
//}
//
//type idleConn struct {  // 连接的基类
//	conn *corn.Client
//	t    time.Time
//}
//
//// channelPool 存放连接信息
//type channelPool struct {
//	mu                       sync.RWMutex
//	conns                    chan *idleConn     // 存放连接的chan
//	factory                  func() (*corn.Client, error)
//	close                    func(*corn.Client) error
//	ping                     func(*corn.Client) error
//	idleTimeout, waitTimeOut time.Duration
//	maxActive                int
//	openingConns             int
//}
//
//
//// NewChannelPool 初始化连接
//func NewChannelPool(poolConfig *rpcConfig) (Pool, error) {
//	if ! (poolConfig.InitialCap <= poolConfig.MaxIdle && poolConfig.MaxCap >= poolConfig.MaxIdle && poolConfig.InitialCap >= 0 ){
//		return nil, errors.New("invalid capacity settings")
//	}
//	if poolConfig.Factory == nil {
//		return nil, errors.New("invalid factory func settings")
//	}
//	if poolConfig.Close == nil {
//		return nil, errors.New("invalid close func settings")
//	}
//
//	c := &channelPool{
//		conns:        make(chan *idleConn, poolConfig.MaxIdle),   //带缓冲区的chan
//		factory:      poolConfig.Factory,
//		close:        poolConfig.Close,
//		idleTimeout:  poolConfig.IdleTimeout,
//		maxActive:    poolConfig.MaxCap,
//		openingConns: poolConfig.InitialCap,
//	}
//
//	if poolConfig.Ping != nil {
//		c.ping = poolConfig.Ping
//	}
//
//	for i := 0; i < poolConfig.InitialCap; i++ {   // 初始放入 InitialCap 个建立好的连接
//		conn, err := c.factory()
//		if err != nil {
//			c.Release()
//			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
//		}
//		c.conns <- &idleConn{conn: conn, t: time.Now()}
//	}
//
//	return c, nil
//}
//
//// getConns 获取所有连接
//func (c *channelPool) getConns() chan *idleConn {
//	c.mu.Lock()
//	conns := c.conns
//	c.mu.Unlock()
//	return conns
//}
//
//// Get 从pool中取一个连接
//func (c *channelPool) Get() (*corn.Client, error) {
//	conns := c.getConns()
//	if conns == nil {
//		return nil, ErrClosed
//	}
//	for {
//		select {
//		case wrapConn := <-conns:   // 如果chan中存在建立好的连接则直接用。
//			if wrapConn == nil {
//				return nil, ErrClosed
//			}
//			//判断是否超时，超时则丢弃
//			if timeout := c.idleTimeout; timeout > 0 {
//				if wrapConn.t.Add(timeout).Before(time.Now()) {
//					//丢弃并关闭该连接
//					c.Close(wrapConn.conn)
//					continue
//				}
//			}
//			//判断是否失效，失效则丢弃，如果用户没有设定 ping 方法，就不检查
//			if c.ping != nil {
//				if err := c.Ping(wrapConn.conn); err != nil {
//					c.Close(wrapConn.conn)
//					continue
//				}
//			}
//			return wrapConn.conn, nil
//		default:  // 否则 排队
//			c.mu.Lock()
//
//			defer c.mu.Unlock()
//			if c.openingConns >= c.maxActive {  // 如果超过上限，报错
//				return nil, errors.New("MaxActiveConnReached")
//			}
//			if c.factory == nil {  // 没设置创建连接的函数,报错
//				return nil, ErrClosed
//			}
//			conn, err := c.factory()  // 创建一个新的连接
//			if err != nil {
//				return nil, err
//			}
//			c.openingConns++
//			return conn, nil
//		}
//	}
//}
//
//// Put 将连接放回pool中
//func (c *channelPool) Put(client *corn.Client) error {
//	if client == nil {
//		return errors.New("connection is nil. rejecting")
//	}
//
//	c.mu.Lock()
//
//	if c.conns == nil {
//		c.mu.Unlock()
//		return c.Close(client)
//	}
//
//	select {
//	case c.conns <- &idleConn{conn: client, t: time.Now()}:
//		c.mu.Unlock()
//		return nil
//	default:
//		c.mu.Unlock()
//		//连接池已满，直接关闭该连接
//		return c.Close(client)
//	}
//}
//
//// Close 关闭单条连接
//func (c *channelPool) Close(client *corn.Client) error {
//	if client == nil {
//		return errors.New("connection is nil. rejecting")
//	}
//	c.mu.Lock()
//	defer c.mu.Unlock()
//	if c.close == nil {
//		return nil
//	}
//	c.openingConns--
//	return c.close(client)
//}
//
//// Ping 检查单条连接是否有效
//func (c *channelPool) Ping(client *corn.Client) error {
//	if client == nil {
//		return errors.New("connection is nil. rejecting")
//	}
//	return c.ping(client)
//}
//
//// Release 释放连接池中所有连接
//func (c *channelPool) Release() {
//	c.mu.Lock()
//	conns := c.conns
//	c.conns = nil
//	c.factory = nil
//	c.ping = nil
//	closeFun := c.close
//	c.close = nil
//	c.mu.Unlock()
//
//	if conns == nil {
//		return
//	}
//
//	close(conns)
//	for wrapConn := range conns {
//		//log.Printf("Type %v\n",reflect.TypeOf(wrapConn.conn))
//		closeFun(wrapConn.conn)
//	}
//}
//
//// Len 连接池中已有的连接
//func (c *channelPool) Len() int {
//	return len(c.getConns())
//}
