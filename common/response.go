package common

// ErrorRes api 에서 사용될 공통 사용 에러 구조체
type ErrorRes struct {
	ServiceId int    `json:"service_id,omitempty"` // used in websocket handler
	TaskId    int    `json:"task_id,omitempty"`    // used in websocket handler
	Packet    string `json:"packet,omitempty"`     // used in websocket handler
	Code      int    `json:"code"`                 // 0 = success , handler error(<0) code , Fault, Failure , (>0) Informational Message, Exception Handling

	// 아래는 code != 0 일 때
	Key         string `json:"key,omitempty"`         // key for user friendly message
	Message     string `json:"message,omitempty"`     // Error Code Short Text ( explanation for Code )
	ErrorText   string `json:"error,omitempty"`       // err.Error()
	Description string `json:"description,omitempty"` // message to client programmer ^^
	//Params interface{} `json:"params,omitempty"` // parameters for user friendly message
}

// AccessRes 부로 노출용으로 사용됨, 혹시 몰라서
type AccessRes struct {
	LangCode         string `json:"lang_code,omitempty"`          // 언어코드
	RegionCode       string `json:"region_code,omitempty"`        // 지역(국가) 코드 (로그인상태가 아닌경우는 AccessRegionCode 처럼 ip 기반, 로그인 상태라면 유저의 국가 코드, 즉 미국 고른 유저가 한국 ip 로 접근하면 US 인거임)
	CurrencyCode     string `json:"currency_code,omitempty"`      // 통화 코드
	AccessRegionCode string `json:"access_region_code,omitempty"` // 접속한 국가 코드 (ip 기반 값이다, 즉 미국 고른 유저가 한국 ip 에서의 요청이면 KR 인거임)
	TimeZone         string `json:"time_zone,omitempty"`          // 접속한 타임존
}

// Response api 에서 사용될 기본 response
type Response struct {
	Req    interface{} `json:"req,omitempty"`    // request 값을 저장해서 남기기 위해 또는 받아가는 쪽에서 자기가 보낸 걸 알게 하기 위해, 좀 의미 없을 것 같음
	Res    interface{} `json:"res,omitempty"`    // response
	Error  ErrorRes    `json:"error"`            // 에러
	Access *AccessRes  `json:"access,omitempty"` // 필요시 사용하면 됨.
}
