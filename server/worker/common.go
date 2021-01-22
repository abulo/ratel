package worker

const (
	MIN_DATA_SIZE = 12
	UINT32_SIZE   = 4
	QUEUE_SIZE    = 64
	BUFFER_SIZE   = 512
	PARAMS_SCOPE  = 0x3A

	//package data type
	PDT_OK            = 1
	PDT_ERROR         = 2
	PDT_CAN_DO        = 3
	PDT_CANT_DO       = 4
	PDT_NO_JOB        = 5
	PDT_HAVE_JOB      = 6
	PDT_TOSLEEP       = 7
	PDT_WAKEUP        = 8
	PDT_WAKEUPED      = 9
	PDT_S_GET_DATA    = 10
	PDT_S_RETURN_DATA = 11
	PDT_W_GRAB_JOB    = 12
	PDT_W_ADD_FUNC    = 13
	PDT_W_DEL_FUNC    = 14
	PDT_W_RETURN_DATA = 15
	PDT_C_DO_JOB      = 16
	PDT_C_GET_DATA    = 17
)

const (
	CONN_TYPE_INIT   = 0
	CONN_TYPE_SERVER = 1
	CONN_TYPE_WORKER = 2
	CONN_TYPE_CLIENT = 3
	PARAMS_TYPE_ONE  = 4
	PARAMS_TYPE_MUL  = 5
	JOB_STATUS_INIT  = 6
	JOB_STATUS_DOING = 7
	JOB_STATUS_DONE  = 8
)

type RetStruct struct {
	Code int
	Msg  string
	Data []byte
}
