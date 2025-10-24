package game

type Symboler interface {
	String() string
	ID() int32 // 返回符号编号，用于父类计算
}
