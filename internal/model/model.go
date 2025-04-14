package model

type Quote struct {
	ID     int    `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type FufufafaQuote struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Datetime string `json:"datetime"`
	Doksli   string `json:"doksli"`
	ImageURL string `json:"image_url"`
}

type SPXData struct {
	Retcode int `json:"retcode"`
	Data    struct {
		FulfillmentInfo struct {
			DeliverType int `json:"deliver_type"`
		} `json:"fulfillment_info"`
		SlsTrackingInfo struct {
			SlsTn            string `json:"sls_tn"`
			ClientOrderID    string `json:"client_order_id"`
			ReceiverName     string `json:"receiver_name"`
			ReceiverTypeName string `json:"receiver_type_name"`
			Records          []struct {
				TrackingCode    string `json:"tracking_code"`
				TrackingName    string `json:"tracking_name"`
				Description     string `json:"description"`
				DisplayFlag     int    `json:"display_flag"`
				ActualTime      int    `json:"actual_time"`
				Operator        string `json:"operator"`
				OperatorPhone   string `json:"operator_phone"`
				ReasonCode      string `json:"reason_code"`
				ReasonDesc      string `json:"reason_desc"`
				Epod            string `json:"epod"`
				PinCode         string `json:"pin_code"`
				CurrentLocation struct {
					LocationName     string `json:"location_name"`
					LocationTypeName string `json:"location_type_name"`
					Lng              string `json:"lng"`
					Lat              string `json:"lat"`
					FullAddress      string `json:"full_address"`
				} `json:"current_location"`
				NextLocation struct {
					LocationName     string `json:"location_name"`
					LocationTypeName string `json:"location_type_name"`
					Lng              string `json:"lng"`
					Lat              string `json:"lat"`
					FullAddress      string `json:"full_address"`
				} `json:"next_location"`
				DisplayFlagV2     int    `json:"display_flag_v2"`
				BuyerDescription  string `json:"buyer_description"`
				SellerDescription string `json:"seller_description"`
				MilestoneCode     int    `json:"milestone_code"`
				MilestoneName     string `json:"milestone_name"`
			} `json:"records"`
		} `json:"sls_tracking_info"`
		IsInstantOrder      bool `json:"is_instant_order"`
		IsShopeeMarketOrder bool `json:"is_shopee_market_order"`
	} `json:"data"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Debug   string `json:"debug"`
}
