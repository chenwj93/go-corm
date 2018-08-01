package utils

//time format
const (
	TIME_FORMAT_1 = "2006-01-02 15:04:05"
	TIME_FORMAT_2 = "2006/01/02 15:04:05"
	TIME_FORMAT_3 = "2006-01-02"
	TIME_FORMAT_4 = "2006/01/02"

	TIME_FORMAT_M    = "2006-01-02 15:04"
	TIME_FORMAT_H_24 = "2006-01-02 15"
	TIME_FORMAT_H_12 = "2006-01-02 03"

	TIME_REG_1 = "^[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-4]):([0-5][0-9]):([0-5][0-9])$"
	TIME_REG_2 = "^[0-9]{4}/(0[1-9]|1[0-2])/(0[1-9]|[1-2][0-9]|3[0-1]) ([0-1][0-9]|2[0-4]):([0-5][0-9]):([0-5][0-9])$"
	TIME_REG_3 = "^[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])"
	TIME_REG_4 = "^[0-9]{4}/(0[1-9]|1[0-2])/(0[1-9]|[1-2][0-9]|3[0-1])"
)

//
const (
	OPERATE_SUCCESS   = "操作成功"
	OPERATE_FAILED    = "操作失败"
	DATA_FORMAT_FALSE = "数据格式不正确！"
	DATA_MISTAKE      = "数据不正确！"
	DATA_DUPLICATE    = "关键数据重复"
	NET_FAILED        = "网络错误"
	OPERATE_DENY      = "非法访问"
	OPERATE_BUSY      = "操作频繁"
	CHECK_FAILED      = "验证不通过"
	NOT_FOUND         = "method not found"
)

const (
	EMPTY_STRING  = ""
	COMMA         = ","
	SEMICOLON     = ";"
	COLON         = ":"
	URL_SEPARATOR = "/"
)

const (
	WECHAT    = "微信"
	BACKSTAGE = "后台"
	INTERFACE = "接口创建"
)

const (
	PROTOCOL_DATA    = "protocol"
	DATA_SOURCE_DATA = "data.source"
	EXTENTION_DATA   = "extention"
)

const (
	LOCATION = "Asia/Shanghai"
)

const (
	ZERO int = iota
	ONE
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
)

const (
	ZERO_I int8 = iota
	ONE_I
	TWO_I
	THREE_I
	FOUR_I
	FIVE_I
	SIX_I
	SEVEN_I
	EIGHT_I
	NINE_I
	TEN_I
)

const (
	ZERO_B int64 = iota
	ONE_B
	TWO_B
	THREE_B
	FOUR_B
	FIVE_B
	SIX_B
	SEVEN_B
	EIGHT_B
	NINE_B
	TEN_B
)
