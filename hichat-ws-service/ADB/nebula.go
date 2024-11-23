package adb

import (
	"errors"
	"fmt"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"go-websocket-server/config"
	"log"
	"strings"
	"time"
)

const (
	useHTTP2 = false
)

type nebulaStruct struct {
	conn   *nebula.SessionPool
	Status int
}

var NebulaInstance = &nebulaStruct{
	conn:   nil,
	Status: 0,
}

func (n *nebulaStruct) createNebula() {

	hostAddress := nebula.HostAddress{Host: config.NebulaAddress, Port: config.NebulaPort}

	// Create configs for session pool
	confnebula, err := nebula.NewSessionPoolConf(
		config.NebulaUserName,
		config.NebulaPassWord,
		[]nebula.HostAddress{hostAddress},
		"HiChat",
		nebula.WithHTTP2(useHTTP2),
	)
	if err != nil {
		log.Panic(fmt.Sprintf("创建 Nebula 配置文件失败, %s\n", err.Error()))
	}

	sessionPool, err := nebula.NewSessionPool(*confnebula, nebula.DefaultLogger{})
	if err != nil {
		log.Panic(fmt.Sprintf("初始化 Nebula失败, %s\n", err.Error()))
	}

	n.conn = sessionPool

}

func (n *nebulaStruct) GetNebulaSession() *nebula.SessionPool {
	if n.conn == nil {
		n.createNebula()
	}
	return n.conn
}

func (n *nebulaStruct) CloseNebula() {
	n.conn.Close()
	n.conn = nil
}

func (n *nebulaStruct) InsertEdge(edgeType, srcVid, dstVid string, propKey []string, propValue []any) error {
	// 检查 propKey 和 propValue 数量是否一致
	if len(propKey) != len(propValue) {
		return fmt.Errorf("propKey 和 propValue 的长度不一致")
	}

	// 构建属性部分 "(key1, key2, ...)" 和 "(val1, val2, ...)"
	keys := "(" + joinStrings(propKey, ", ") + ")"
	values := "(" + joinValues(propValue) + ")"

	// 构建 nGQL 插入语句
	query := fmt.Sprintf(
		`INSERT EDGE %s %s VALUES "%s"->"%s": %s;`,
		edgeType, keys, srcVid, dstVid, values,
	)

	// 执行 nGQL 查询
	result, err := NebulaInstance.GetNebulaSession().Execute(query)
	if err != nil {
		return fmt.Errorf("执行 nGQL 失败: %w", err)
	}

	if ok := checkResultSet(query, result); !ok {
		return errors.New("执行nebula操作失败")
	}

	return nil
}

// DeleteEdge 删除指定边类型的边
func (n *nebulaStruct) DeleteEdge(edgeType, srcVid, dstVid string, IgnoreDirections bool) error {
	// 构建删除边的查询语句
	// 执行查询
	// 检查结果

	query := fmt.Sprintf("DELETE EDGE %s '%s' -> '%s' ", edgeType, srcVid, dstVid)
	fmt.Println(query)
	result, err := NebulaInstance.GetNebulaSession().Execute(query)
	if err != nil {
		return fmt.Errorf("执行nebula操作失败: %v", err)
	}
	if ok := checkResultSet(query, result); !ok {
		return errors.New("执行nebula操作失败")
	}

	//当忽略方向时,执行反向删除
	if IgnoreDirections {
		query2 := fmt.Sprintf("DELETE EDGE %s '%s' -> '%s' ", edgeType, dstVid, srcVid)
		result2, err := NebulaInstance.GetNebulaSession().Execute(query2)
		if err != nil {
			return fmt.Errorf("执行nebula操作失败: %v", err)
		}
		if ok := checkResultSet(query, result2); !ok {
			return errors.New("执行nebula操作失败")
		}
	}

	return nil
}
func checkResultSet(prefix string, res *nebula.ResultSet) bool {
	if !res.IsSucceed() {
		log.Println(fmt.Sprintf("%s, ErrorCode: %v, ErrorMsg: %s", prefix, res.GetErrorCode(), res.GetErrorMsg()))
		return false
	}
	return true
}

// 辅助函数：将字符串切片拼接成 "key1, key2, key3" 的格式
func joinStrings(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// 辅助函数：格式化不同类型的值，并拼接成 "val1, val2, val3" 格式
func joinValues(values []any) string {
	var result []string
	for _, val := range values {
		result = append(result, formatValue(val))
	}
	return strings.Join(result, ", ")
}

// 辅助函数：格式化值以符合 nGQL 语法
func formatValue(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v) // 字符串需要用单引号括起来
	case int, int64, float64:
		return fmt.Sprintf("%v", v) // 数值直接使用其值
	case bool:
		return fmt.Sprintf("%t", v) // 布尔值转换为 "true"/"false"
	case time.Time:
		return fmt.Sprintf("datetime('%s')", v.Format("2006-01-02 15:04:05"))
	default:
		return "NULL" // 其他未知类型处理为 NULL
	}
}
