package requests

type MakeUserPassive struct {
	UserID uint `locals:"userid" json:"-" validate:"required"`
}

type GetUser struct {
	UserID uint `locals:"userid" json:"-" validate:"required"`
}

type SetUserSettings struct {
	UserID uint `locals:"userid" json:"-" validate:"required"`

	Address                string `json:"address"`
	TaxOffice              string `json:"tax_office"`
	TaxNo                  string `json:"tax_no"`
	TaxTitle               string `json:"tax_title"`
	CityID                 uint   `json:"city_id"`
	DistrictID             uint   `json:"district_id"`
	SchoolType             string `json:"school_type"`
	MessageType            string `json:"message_type"`
	InfoMessage            bool   `json:"info_message"`
	DefaultSenderID        uint   `json:"default_sender_id"`
	ShowRepresentativeInfo bool   `json:"show_representative_info"`
	Signature              string `json:"signature"`
	AttendanceNotification bool   `json:"attendance_notification"`
	AttendanceMessage      string `json:"attendance_message"`
}

type GetDashboardData struct {
	UserID   uint   `locals:"userid" json:"-" validate:"required"`
	Timespan string `query:"timespan" validate:"required,oneof=today week month year" example:"today" description:"today, week, month, year"`
}
