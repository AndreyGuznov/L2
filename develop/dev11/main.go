package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

/*
=== HTTP server ===
Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.
В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

const dateLayout = "2006-01-02"

type Logger struct {
	handler http.Handler
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handler: handlerToWrap}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

type Event struct {
	UserID      int       `json:"user_id"`
	EventID     int       `json:"event_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func (e *Event) Decode(r io.Reader) error {
	err := json.NewDecoder(r).Decode(&e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Event) Validate() error {
	if e.UserID <= 0 {
		return fmt.Errorf("Invalid user id: %v;", e.UserID)
	}
	if e.EventID <= 0 {
		return fmt.Errorf("Invalid event id: %v;", e.EventID)
	}
	if e.Title == "" {
		return fmt.Errorf("Title is empty")
	}
	return nil
}

type Storage struct {
	mu     *sync.Mutex
	events map[int][]Event
}

func (s *Storage) Create(e *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if events, ok := s.events[e.UserID]; ok {
		for _, event := range events {
			if event.EventID == e.EventID {
				return fmt.Errorf("Event already exist")
			}
		}
	}
	s.events[e.UserID] = append(s.events[e.UserID], *e)
	return nil
}

func (s *Storage) Update(e *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ind := -1
	events := make([]Event, 0)
	ok := false
	if events, ok = s.events[e.UserID]; !ok {
		return fmt.Errorf("User does not find")
	}
	for i, event := range events {
		if event.EventID == e.EventID {
			ind = i
			break
		}
	}
	if ind == -1 {
		return fmt.Errorf("There is no event")
	}
	s.events[e.UserID][ind] = *e

	return nil
}

func (s *Storage) Delete(e *Event) (*Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ind := -1
	events := make([]Event, 0)
	var ok bool
	if events, ok = s.events[e.UserID]; !ok {
		return nil, fmt.Errorf("User does not find")
	}
	for i, event := range events {
		if event.EventID == e.EventID {
			ind = i
			break
		}
	}
	if ind == -1 {
		return nil, fmt.Errorf("There is no event")
	}
	eventsLength := len(s.events[e.UserID])
	deletedEvent := s.events[e.UserID][ind]
	s.events[e.UserID][ind] = s.events[e.UserID][eventsLength-1]
	s.events[e.UserID] = s.events[e.UserID][:eventsLength-1]

	return &deletedEvent, nil
}

func (s *Storage) GetEventsForDay(userID int, date time.Time) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []Event

	events := make([]Event, 0)
	ok := false

	if events, ok = s.events[userID]; !ok {
		return nil, fmt.Errorf("User does not find")
	}
	for _, event := range events {
		if event.Date.Year() == date.Year() && event.Date.Month() == date.Month() && event.Date.Day() == date.Day() {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForWeek(userID int, date time.Time) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []Event
	events := make([]Event, 0)
	var ok bool
	if events, ok = s.events[userID]; !ok {
		return nil, fmt.Errorf("User does not find")
	}
	for _, event := range events {
		y1, w1 := event.Date.ISOWeek()
		y2, w2 := date.ISOWeek()
		if y1 == y2 && w1 == w2 {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) GetEventsForMonth(userID int, date time.Time) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []Event
	events := make([]Event, 0)
	var ok bool
	if events, ok = s.events[userID]; !ok {
		return nil, fmt.Errorf("User does not find")
	}
	for _, event := range events {
		if event.Date.Year() == date.Year() && event.Date.Month() == date.Month() {
			result = append(result, event)
		}
	}
	return result, nil
}

var storage Storage = Storage{
	events: make(map[int][]Event),
	mu:     &sync.Mutex{},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", CreateEventHandler)
	mux.HandleFunc("/update_event", UpdateEventHandler)
	mux.HandleFunc("/delete_event", DeleteEventHandler)

	mux.HandleFunc("/events_for_day", EventsForDayHandler)
	mux.HandleFunc("/events_for_week", EventsForWeekHandler)
	mux.HandleFunc("/events_for_month", EventsForMonthHandler)

	wrappedMux := NewLogger(mux)

	port := ":1234"
	func() {
		temp := os.Getenv("PORT")
		if temp != "" {
			port = temp
		}
	}()

	log.Println("Server is listening")
	log.Fatalln(http.ListenAndServe(port, wrappedMux))
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var e Event
	if err := e.Decode(r.Body); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := e.Validate(); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := storage.Create(&e); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	resultResponse(w, "Event has been created", []Event{e}, http.StatusCreated)
	fmt.Println(storage.events)
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	var e Event
	if err := e.Decode(r.Body); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := e.Validate(); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := storage.Update(&e); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	resultResponse(w, "Event has been updated successfully!", []Event{e}, http.StatusOK)
	fmt.Println(storage.events)
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	var e Event
	if err := e.Decode(r.Body); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	var deletedEvent *Event
	var err error
	if deletedEvent, err = storage.Delete(&e); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	resultResponse(w, "Event has been deleted", []Event{*deletedEvent}, http.StatusOK)
}

func EventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	date, err := time.Parse(dateLayout, r.URL.Query().Get("date"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	var events []Event
	if events, err = storage.GetEventsForDay(userID, date); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	resultResponse(w, "Here: ", events, http.StatusOK)
}

func EventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	date, err := time.Parse(dateLayout, r.URL.Query().Get("date"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	var events []Event
	if events, err = storage.GetEventsForWeek(userID, date); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	resultResponse(w, "Here: ", events, http.StatusOK)
}

func EventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	date, err := time.Parse(dateLayout, r.URL.Query().Get("date"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	var events []Event
	if events, err = storage.GetEventsForMonth(userID, date); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	resultResponse(w, "Here: ", events, http.StatusOK)
}

func errorResponse(w http.ResponseWriter, e string, status int) {
	errorResponse := struct {
		Error string `json:"error"`
	}{Error: e}

	js, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func resultResponse(w http.ResponseWriter, r string, e []Event, status int) {
	resultResponse := struct {
		Result string  `json:"result"`
		Events []Event `json:"events"`
	}{Result: r, Events: e}

	js, err := json.Marshal(resultResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
