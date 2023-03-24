package booking

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/travisbale/birdies-up/internal/clubhouse"
)

type apiClient interface {
	Login(username string, password string) (*clubhouse.User, error)
	GetGolfCourse(courseName string) (clubhouse.GolfCourse, error)
	GetTeeTimes(date time.Time, courses ...clubhouse.GolfCourse) ([]clubhouse.TeeTime, error)
	LookupMembers() ([]clubhouse.Member, error)
	BookTeeTime(teeTime clubhouse.TeeTime, bookingRequest *clubhouse.BookingRequest) (clubhouse.BookingResponse, error)
}

type logger interface {
	Debug(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
}

type Config struct {
	Client         apiClient
	Username       string
	Password       string
	ConfigFileName string
	Logger         logger
}

type Service struct {
	config Config
	client apiClient
	logger logger
}

func NewService(config Config) *Service {
	return &Service{
		config: config,
		client: config.Client,
		logger: config.Logger,
	}
}

type TeeTime struct {
	WeekDay time.Weekday
	TeeTime time.Time
	MemberIDs []int
}

func (s *Service) Run() error {
	return nil
}

func (s *Service) BookTeeTime(courseName string, targetTime time.Time, playerNames []string) error {
	user, err := s.client.Login(s.config.Username, s.config.Password)
	if err != nil {
		return fmt.Errorf("BookTeeTimes: %w", err)
	}

	course, err := s.client.GetGolfCourse(courseName)
	if err != nil {
		return fmt.Errorf("BookTeeTimes: %w", err)
	}

	teeTimes, err := s.client.GetTeeTimes(targetTime, course)
	if err != nil {
		return fmt.Errorf("BookTeeTimes: %w", err)
	}

	teeTime, err := s.getTargetTeeTime(teeTimes, targetTime)
	if err != nil {
		return fmt.Errorf("BookTeeTimes: %w", err)
	}

	members, err := s.client.LookupMembers()
	if err != nil {
		return fmt.Errorf("BookTeeTimes: %w", err)
	}

	bookingRequest := s.createBookingRequest(user, members, playerNames...)

	_, err = s.client.BookTeeTime(teeTime, bookingRequest)
	if err != nil {
		return fmt.Errorf("BookTeeTimes: %w", err)
	}

	data, _ := json.MarshalIndent(bookingRequest, "", "    ")

	fmt.Printf(string(data))

	return nil
}

func (s *Service) getTargetTeeTime(teeTimes []clubhouse.TeeTime, targetTime time.Time) (clubhouse.TeeTime, error) {
	targetString := targetTime.Format("15:04:05")
	targetTeeTime := teeTimes[0]

	for _, teeTime := range teeTimes {
		if teeTime.TeeTime == targetString {
			if teeTime.IsBookable && teeTime.AvailPlayers == 5 {
				return teeTime, nil
			}

			s.logger.Error("the tee time has already been booked", "teetime", targetTime.Format("Jan 2, 2006, 15:04:05"))
			return clubhouse.TeeTime{}, fmt.Errorf("the %s tee time has already been booked", targetTime.Format("Jan 2, 2006 15:04:05"))
		}
	}

	s.logger.Error("the tee time does not exist", "teetime", targetTime.Format("Jan 2, 2006, 15:04:05"))
	return clubhouse.TeeTime{}, fmt.Errorf("the %s tee time does not exist", targetTime.Format("Jan 2, 2006 15:04:05"))
}

func (s *Service) createBookingRequest(user *clubhouse.User, members []clubhouse.Member, playerNames ...string) *clubhouse.BookingRequest {
	playerReservations := make([]clubhouse.PlayerReservation, len(playerNames)+1)

	for index, name := range playerNames {
		for _, member := range members {
			if member.Name == name {
				s.logger.Debug(fmt.Sprintf("adding %s (%d)", member.Name, member.ID))
				playerReservations[index] = clubhouse.PlayerReservation{
					MemberID: member.ID,
				}
				break
			}
		}
	}

	playerReservations[len(playerNames)] = clubhouse.PlayerReservation{
		MemberID: user.ID,
	}

	return &clubhouse.BookingRequest{
		Mode:         "Booking",
		OwnerID:      user.ID,
		Reservations: playerReservations,
	}
}
