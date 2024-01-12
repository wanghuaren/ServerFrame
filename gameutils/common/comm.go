package common

const MICRO_ERROR int32 = -1
const MICRO_ERROR_TIMEOUT int32 = -2
const MICRO_ERROR_PROTOBUF_NIL int32 = -3
const MICRO_ERROR_PROTOBUF_DATA_NIL int32 = -4

const CLIENT_TRANS_ERROR int32 = -5
const CLIENT_TRANS_ERROR_PARSE int32 = -6
const CLIENT_TRANS_ERROR_MICRO int32 = -7
const CLIENT_TRANS_ERROR_DECODE int32 = -8

const CLIENT_TRANS_TOKEN_ERROR int32 = 401
const CLIENT_TRANS_ERROR_NUM_ERROR int32 = 402
const CLIENT_TRANS_MSG_NUM_ERROR int32 = 403

const DB_ORDER_ADD string = "add"
const DB_ORDER_DEL string = "del"
const DB_ORDER_EDIT string = "edit"
const DB_ORDER_FIND string = "find"
const DB_ORDER_FIND_ADD string = "find_add"
const DB_ORDER_FIND_TABLE string = "find_table"

const DB_FIND_KEY string = "find_key"

const DB_USERAPI_GetUserData string = "GetUserData"
const DB_USERAPI_GetUserDataFromToken string = "GetUserDataFromToken"
const DB_USERAPI_GetUserDataRank string = "GetUserDataRank"
const DB_USERAPI_SetUserDataFromToken string = "SetUserDataFromToken"
const DB_USERAPI_GetUserTokenFixed string = "GetUserTokenFixed"

// const DB_USERAPI_GetUserTokenFixedFromToken string = "GetUserTokenFixedFromToken"
const DB_USERAPI_GetUserToken string = "GetUserToken"
const DB_USERAPI_GetUserTokenFromKey string = "GetUserTokenFromKey"

// const DB_USERAPI_SetUserTokenFromKey string = "SetUserTokenFromKey"
const DB_USERAPI_CleanUserTokenFromKey string = "CleanUserTokenFromKey"

// const DB_USERAPI_CleanUserIndex string = "CleanUserIndex"
const DB_USERAPI_HeartJump string = "HeartJump"
const DB_USERAPI_AddUserItem string = "AddUserItem"
const DB_USERAPI_UseupUserItem string = "UseupUserItem"
const DB_USERAPI_FinishTask string = "FinishTask"
const DB_USERAPI_Login string = "Login"

const DB_TABLEAPI_InitStaticTables string = "InitStaticTables"

var HTTP_TIMEOUT_MILLISECOND = 10000
var MICRO_TIMEOUT_MILLISECOND = 8000

const MaxBagCount int = 18

const DefaultWeaponID int32 = 100001
