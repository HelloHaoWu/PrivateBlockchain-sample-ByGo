package BLC

// 该文件, 存储各种文件格式间转化的方法
import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

// bytes.buffer为一个结构, 其定义如下
//type Buffer struct {
//	buf []byte
// 	off int
// 	lastRead readOp
//}
// ↑ off标记读到的位置, lastRead用于标记上次是否是读操作, 用于读操作的回退;

// 将int64转换为字节数组↓
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	// 声明一个Buffer结构体的四种方法：
	// var b bytes.Buffer → 直接定义一个Buffer变量，不用初始化，可以直接使用
	// b := new(bytes.Buffer) → 使用New返回Buffer变量
	// b := bytes.NewBuffer(s []byte) → 从一个[]byte切片，构造一个Buffer
	// b := bytes.NewBufferString(s string) → 从一个string变量，构造一个Buffer
	err := binary.Write(buff, binary.BigEndian, num) //容器, 写入模式, 需要写入的数据
	// 在上述创建的容器(buff)中写入num这个数据
	// binary.BigEndian（大端模式）：内存的低地址存放着数据高位
	// binary.LittleEndian(小端模式)：内存的低地址存放着数据地位
	// 举个栗子：如一个 var a = 0x11223344，对于这个变量,
	// 最高字节为0x11，最低字节为0x44。假设在内存中分配地址如下(地址都是连续的)
	// 0x0001 → 0x0002 → 0x0003 → 0x0004
	// 当分别处于大小端模式下的内容存放如下
	// (1)大端模式存储（存储地址为16位）
	// 地址 数据
	// 0x0004(高地址) 0x44
	// 0x0003 0x33
	// 0x0002 0x22
	// 0x0001(低地址) 0x11
	// (2)小端模式存储（存储地址为16位）
	// 地址 数据
	// 0x0004(高地址) 0x11
	// 0x0003 0x22
	// 0x0002 0x33
	// 0x0001(低地址) 0x44
	if err != nil {
		log.Panic(err) // 带时间戳地打印"err"
	}
	// log.Panic实例↓
	// log.Print("1111")
	// 运行结果↓
	// 2018/08/20 17:49:28 1111
	return buff.Bytes()
}

// 标准json字符串转数组
func JsonToArray(jsonString string) []string {
	var sArr []string
	// ↓ err != nil → 报错了
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr // ← 返回转化完毕的数组
}

// ↓字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
