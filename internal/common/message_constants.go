package common

const (
	VerifyCommandName  string = "/verify"
	VolumeCommandName         = "/volume"
	StartCommandName          = "/start"
	HelpCommandName           = "/help"
	StatusCommandName         = "/status"
	JoinCommandName           = "/join"
	AccountCommandName        = "/account"
)
const (
	WelcomeMessage string = `🦀≡≡≡≡≡≡≡≡▷►◈◄◁≡≡≡≡≡≡≡≡🦀

🔝歡迎各位加入🪣比奇堡海之霸VIP🔇
因我個人的策略關係，(槓桿高，盈虧比高)📈
有些訂單接受度較低，在大群發布容易被噴
因此簡單設立一個門檻，避免路人粉瞎操作爆倉
加上之後會有獎勵活動避免註冊後白嫖
每個月最低10000u交易量(含槓桿)非常低的標準
每月1號核對，若交易額不符合要求將會移除VIP群及交流群
直到交易額再次達到10000u或一個月後點即可/rejoin重新加回
如果不知道是否符合要求，機器人也有交易額查詢功能可以使用

內容：
🅰️最高20%合約20%+20%現貨手續費減免✨
🅱️專屬團隊bitget跟單服務
🆎蟹老闆🦀親自帶單
⚡️VIP群高盈虧比策略分享
🧽交流群加入資格
加入步驟：
1️⃣點擊鏈接註冊賬號⭐️
https://partner.bitget.fit/bg/MrKrabs
2️⃣發送UID給機器人確認♥️
https://t.me/wedjatbtcVIP_bot
⚠️邀請鏈接為一次性使用⚠️
⚠️註冊後記得點擊下方加入在退出⚠️
⚠️否則無法再點擊⚠️

☢️以下是不同會員入群指令以及交易額查詢☢️

/start    	  - 開始使用機器人
/help     	  - 了解所有指令說明 請輸入此指令
/verify <uid>   	  - 驗證uid指令 請輸入你的數字UID
/volume <uid>   - 交易總額查詢 請輸入此指令
/account <uid>  - 更改電報帳號綁定

有任何疑問請直接私訊本人謝謝🕳
https://t.me/wedjatbtc

🦀≡≡≡≡≡≡≡≡▷►◈◄◁≡≡≡≡≡≡≡≡🦀`

	HelpMessage = "```\n" +
		`/start          - 開始使用機器人
/help           - 了解所有指令說明 請輸入此指令
/status <uid>   - 查詢目前電報帳號狀態
/verify <uid>   - 驗證uid指令，請輸入你的數字UID
/volume <uid>   - 交易總額查詢，請輸入此指令
/account <uid>  - 更改電報帳號綁定` +
		"\n```"

	ProcessingMessage          = "正在驗證 UID，請稍候..."
	ServerErrorMessage         = "驗證服務暫時無法使用，請稍後重試❌"
	InternalServerErrorMessage = `伺服器處理過程中發生錯誤，請稍後重試`
)

const (
	InvalidCommandFormatMessage string = "❌请使用正确的格式：%s <UID>\n範例：%s 123456"
	InvalidUIDFormatMessage            = "❌無效的UID格式\n範例：%s 123456"
)

const (
	SuccessVerifyReplyMessage            string = `🦀您已驗證成功!感謝關注!✅\n以下是交流群以及VIP群的鏈接`
	InvalidUidVerifyReplyMessage                = `🦀您輸入的UID不存在 驗證失敗❌ 請查詢正確後再次輸入`
	ExistsUidVerifyReplyMessage                 = `🦀您要驗證的uid已存在,無須再次驗證!祝您交易順利!✅`
	ExistsSocialUserIdVerifyReplyMessage        = `🦀您已綁定過電報帳號 請使用/account變更您的綁定電報帳號❌`
	InvalidUidStatusMessage                     = "🦀此UID所綁定的社交帳號狀態為非活躍 請使用/volume %s 檢查您的交易額度是否達標 或聯絡群組主❌"
	DuplicatedUserReplyMessage                  = "🦀此UID所綁定的社交帳號無需更改"
)

const (
	SuccessVolumeReplyMessage string = `🔎查詢成功,距離本月1號到今日,您的交易額為: USDT$%.2f`
	FailureVolumeReplyMessage        = `❌查詢失敗請重試`
)

const (
	MemberStatusReplyMessage string = "⚠️ 您目前使用該uid: %s 查詢的電報用戶群組狀態為： %s"
	MemberInfoUpdatedMessage        = "🦀您目前的社交帳號資訊已更新成功✅"
)

const (
	UserWarningMessage string = "⚠️ @%s 請不要在群組中發送任何与指令 電報链接 網頁連結 UID...等等敏感訊息 謝謝合作"
)
