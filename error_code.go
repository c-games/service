package service

// app、service
const Code_Success int = 0

const (
	Code_No_Such_Command  int = 3000
	Code_Unmarshal_Failed int = 3001
	Code_Marshal_Failed   int = 3002
	Code_Unexpected_Data  int = 3003
	Code_Not_Homo_Data    int = 3004
	Code_Empty_Request_Data int = 3005
	Code_Request_Has_Illegal_Zero_Value int = 3006

	// user or admin
	Code_User_Not_Exist             int = 3100
	Code_User_Login_Failed          int = 3101
	Code_User_Already_Exist         int = 3102
	Code_User_Authentication_Failed int = 3103
	Code_User_Validate_Failed       int = 3104
	Code_User_Id_Not_Exist          int = 3105
	Code_Wrong_Password             int = 3106
	Code_Permission_Denied          int = 3107
	Code_Wrong_Account_Type         int = 3108
	Code_Unexpected_Logout          int = 3109

	Code_Wrong_Token           int = 3200
	Code_No_Token              int = 3201
	Code_Multiple_Query_Failed int = 3202

	Code_Order_Not_Exist               int = 3300
	Code_Insufficient_Credit           int = 3301
	Code_Wallet_No_Update_Reason       int = 3302
	Code_Query_Wallet_Failed           int = 3303
	Code_Create_New_Wallet_Failed      int = 3304
	Code_Insert_Wallet_Exchange_Failed int = 3305
	Code_Update_Wallet_Failed          int = 3306
)

// DB
const (
	Code_DB_Error           int = 4000
	Code_Query_No_Result    int = 4001
	Code_DB_No_Row_Affected int = 4002
)

// Calculate 5 開頭
const ()

//draw result or with_draw_result
const(
	Code_Repeat_Draw_Result int = 6000
	Code_Wrong_Game_ID int = 6001
)
