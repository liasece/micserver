package conf

type DBTableConfig struct {
	// 数据库实例索引
	DBIndex uint32 `json:"dbindex"`
	// 该表对应的表名字
	TableName string `json:"name"`
}

type DBConfig struct {
	// 数据库实例
	// 	key 		---  value
	// 	数据库实例索引    连接字符串
	Dbs map[uint32]string `json:"dbs"`
	// 数据库 表 实例
	// 	key 					---  value
	//  表索引，用于哈希Openid        表配置对象
	Tables map[uint32]*DBTableConfig `json:"tables"`
}
