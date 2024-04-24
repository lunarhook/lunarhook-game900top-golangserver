package Utils

var (
	Success_200  = 200
	Error_token_201 = 201
	Error_user_password_202 = 202
	Error_update_userpassword_203 = 203
	Error_invalid_mobile_204  = 204
	Error_invalid_token_205 = 205
	Error_user_create_206 = 206
	Error_invalid_sms_207 = 207
)
var CodeNum map[int]string

func CodeNumInit()  {
	CodeNum = make(map[int]string)
	CodeNum[200] = "Success_200"
	CodeNum[201] = "Error_token_201"
	CodeNum[202] = "Error_user_password_202"
	CodeNum[203] = "Error_update_userpassword_203"
	CodeNum[204] = "Error_invalid_mobile_204"
	CodeNum[205] = "Error_invalid_token_205"
	CodeNum[206] = "Error_user_create_206"
	CodeNum[207] = "Error_invalid_sms_207" //短信验证码失败

	CodeNum[2000] = "userfileinfo_2000"
	//CodeNum[2002] = "Error_Databaseerror_2002"
	CodeNum[2001] = "Error_userfileinfo_2001"
	CodeNum[2002] = "Error_Databaseerror_2002"
	CodeNum[2003] = "Error_Database_Fileerror_2003"
}

