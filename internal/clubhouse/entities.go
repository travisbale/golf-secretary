package clubhouse

import "net/http"

type ApiClient struct {
	baseUrl    string
	httpClient *http.Client
	logger     logger
}

type BaseResponse struct {
	RetCode          int    `json:"retCode"`
	Title            string `json:"title"`
	InfoMsg          string `json:"infoMsg"`
	ErrorMessage     string `json:"errorMessage"`
	DisplayMessage   string `json:"displayMessage"`
	ServerStackTrace string `json:"serverStackTrace"`
	Result           bool   `json:"result"`
}

type ErrorResponse struct {
	BaseResponse
}

func (e ErrorResponse) Error() string {
	return e.ErrorMessage
}

type GolfCourseResponse struct {
	BaseResponse
	GolfCourses []GolfCourse `json:"data"`
}

type GolfCourse struct {
	Name       string `json:"name"`
	Id         int    `json:"id"`
	ExternalId string `json:"externalId"`
}

type UserResponse struct {
	BaseResponse
	User User `json:"data"`
}

type User struct {
	FullName          string `json:"fullName"`
	FirstName         string `json:"firstName"`
	ID                int    `json:"id"`
	LastName          string `json:"lastName"`
	MemberID          int    `json:"memberId"`
	MemberNumber      string `json:"memberNumber"`
	MemberNumberSaved string `json:"memberNumberSaved"`
	Email             string `json:"email"`
	Photo             string `json:"photo"`
	IsEventAdmin      bool   `json:"isEventAdmin"`
	IsPublicUser      bool   `json:"isPublicUser"`
}

type TeeTimesResponse struct {
	BaseResponse
	Data struct {
		Availability []interface{} `json:"availability"`
		TeeSheet     []TeeTime     `json:"teeSheet"`
	} `json:"data"`
}

type TeeTime struct {
	TeeSheetTimeID           int    `json:"teeSheetTimeId"`
	AvailableToClubNetwork   bool   `json:"availableToClubNetwork"`
	AvailableToGlobalNetwork bool   `json:"availableToGlobalNetwork"`
	AvailableToMembers       bool   `json:"availableToMembers"`
	AvailableToPublic        bool   `json:"availableToPublic"`
	EighteenAllowed          bool   `json:"eighteenAllowed"`
	NineAllowed              bool   `json:"nineAllowed"`
	Hole                     string `json:"hole"`
	IsBookable               bool   `json:"isBookable"`
	CanQuickBook             bool   `json:"canQuickBook"`
	AllowReservations        bool   `json:"allowReservations"`
	BookingTimeType          int    `json:"bookingTimeType"`
	BookingTimeTypeTxt       string `json:"bookingTimeTypeTxt"`
	Players                  []struct {
		PlayerType    int    `json:"playerType"`
		PlayerTypeTxt string `json:"playerTypeTxt"`
		PlayerLabel   string `json:"playerLabel"`
	} `json:"players"`
	TeeTime           string      `json:"teeTime"`
	TimeOfDay         int         `json:"timeOfDay"`
	TimeOfDayTxt      string      `json:"timeOfDayTxt"`
	Processed         bool        `json:"processed"`
	LotteryRequests   int         `json:"lotteryRequests"`
	IsTournament      bool        `json:"isTournament"`
	TournamentDetails interface{} `json:"tournamentDetails"`
	AvailPlayers      int         `json:"availPlayers"`
	LockedBy          int         `json:"lockedBy"`
	LockedUntil       interface{} `json:"lockedUntil"`
	LotteryEndTime    string      `json:"lotteryEndTime"`
	CartCost          float64     `json:"cartCost"`
	PlayerCost        float64     `json:"playerCost"`
	TeeSheetBank      struct {
		BankID        int `json:"bankId"`
		BankNumber    int `json:"bankNumber"`
		NumberOfHoles int `json:"numberOfHoles"`
		TeeSheetKey   struct {
			Course                  string      `json:"course"`
			CourseID                int         `json:"courseId"`
			DateOfAvailability      string      `json:"dateOfAvailability"`
			DateOfAvailabilityJS    string      `json:"dateOfAvailabilityJS"`
			DateNoLongerAvailable   interface{} `json:"dateNoLongerAvailable"`
			DateNoLongerAvailableJS interface{} `json:"dateNoLongerAvailableJS"`
			MaxPlayers              int         `json:"maxPlayers"`
			TbdExpiration           interface{} `json:"tbdExpiration"`
		} `json:"teeSheetKey"`
	} `json:"teeSheetBank"`
}

type LookupMembersResponse struct {
	BaseResponse
	Members []Member `json:"data"`
}

type Member struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IsBuddy bool   `json:"isBuddy"`
	Suffix  string `json:"suffix"`
}

type CommitBookingResponse struct {
	BaseResponse
	Data struct {
		LatestTime             interface{} `json:"latestTime"`
		EarliestTime           interface{} `json:"earliestTime"`
		BookingType            string      `json:"bookingType"`
		BookingID              int         `json:"bookingId"`
		ConfirmationNumber     string      `json:"confirmationNumber"`
		Course                 string      `json:"course"`
		CourseID               int         `json:"courseId"`
		NumberOfHoles          int         `json:"numberOfHoles"`
		Date                   string      `json:"date"`
		Time                   string      `json:"time"`
		StartingHole           string      `json:"startingHole"`
		Notes                  string      `json:"notes"`
		ConfirmationCustomText interface{} `json:"confirmationCustomText"`
		Owner                  struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Type      string `json:"type"`
		} `json:"owner"`
		IsOwner      bool `json:"isOwner"`
		AllowCancel  bool `json:"allowCancel"`
		IsActive     bool `json:"isActive"`
		Reservations []struct {
			FirstName     string      `json:"firstName"`
			LastName      string      `json:"lastName"`
			Caddy         bool        `json:"caddy"`
			HasTBD        bool        `json:"hasTBD"`
			FullName      string      `json:"fullName"`
			Transport     string      `json:"transport"`
			TransportNum  int         `json:"transportNum"`
			ClubRental    string      `json:"clubRental"`
			ClubRentalNum interface{} `json:"clubRentalNum"`
			Type          string      `json:"type"`
		} `json:"reservations"`
		TbdeXpirationMsg  interface{} `json:"tbdeXpirationMsg"`
		IsWaitingList     bool        `json:"isWaitingList"`
		PlayersAllowedAdd bool        `json:"playersAllowedAdd"`
		CanEdit           bool        `json:"canEdit"`
	} `json:"data"`
}

type BookingRequest struct {
	Mode             string              `json:"Mode"`
	BookingID        int                 `json:"BookingId"`
	OwnerID          int                 `json:"OwnerId"`
	EditingBookingID interface{}         `json:"editingBookingId"`
	Reservations     []PlayerReservation `json:"Reservations"`
	Holes            int                 `json:"Holes"`
	StartingHole     string              `json:"StartingHole"`
	Wait             bool                `json:"wait"`
	Allowed          interface{}         `json:"Allowed"`
	Enabled          bool                `json:"enabled"`
	StartTime        interface{}         `json:"startTime"`
	EndTime          interface{}         `json:"endTime"`
	Notes            string              `json:"Notes"`
}

type PlayerReservation struct {
	ReservationID   int    `json:"ReservationId"`
	ReservationType int    `json:"ReservationType"`
	FullName        string `json:"FullName"`
	Transport       string `json:"Transport"`
	Caddy           string `json:"Caddy"`
	Rentals         string `json:"Rentals"`
	MemberID        int    `json:"MemberId,omitempty"`
	FirstName       string `json:"FirstName,omitempty"`
	LastName        string `json:"LastName,omitempty"`
	GuestID         int    `json:"GuestId,omitempty"`
}

type BookingResponse struct {
	BaseResponse
}
