package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Store struct {
	ID          string    `sql:"id,primary"`
	Name        string    `sql:"name"`
	Description *string   `sql:"description"`
	Lat         float64   `sql:"lat"`
	Lng         float64   `sql:"lng"`
	Slug        string    `sql:"slug"`
	UserID      string    `sql:"user_id"`
	CreatedAt   time.Time `sql:"created_at"`
	UpdatedAt   time.Time `sql:"updated_at"`

	Entity
}

func NewStore(id, name string, description *string, lat, lng float64, userID string) (*Store, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	if len(name) > 100 {
		return nil, fmt.Errorf("name cannot exceed 100 characters")
	}

	if description != nil && len(*description) > 500 {
		return nil, fmt.Errorf("description cannot exceed 500 characters")
	}

	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}

	if lng < -180 || lng > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}

	now := time.Now().UTC()
	slug := generateSlug(name, now.Unix())

	store := &Store{
		ID:          id,
		Name:        name,
		Description: description,
		Lat:         lat,
		Lng:         lng,
		Slug:        slug,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	data, _ := json.Marshal(StoreEventData{
		ID:          store.GetID(),
		Name:        store.GetName(),
		Description: store.GetDescription(),
		Location:    store.Location(),
		Slug:        store.GetSlug(),
		UserID:      store.GetUserID(),
		CreatedAt:   store.GetCreatedAt(),
		UpdatedAt:   store.GetUpdatedAt(),
	})
	event := Event{
		ID:        fmt.Sprintf("%s_%v", store.GetID(), store.GetCreatedAt().Unix()),
		Type:      StoreCreated,
		Data:      data,
		Timestamp: store.GetCreatedAt(),
	}
	store.events = append(store.events, event)

	return store, nil
}

func (s *Store) Update(name string, description *string, lat, lng float64, userID string) error {
	if userID != s.UserID {
		return fmt.Errorf("user id does not match")
	}

	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("name cannot exceed 100 characters")
	}

	if description != nil && len(*description) > 500 {
		return fmt.Errorf("description cannot exceed 500 characters")
	}

	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}

	if lng < -180 || lng > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}

	if s.Name != name {
		s.Slug = generateSlug(name, time.Now().Unix())
	}

	s.Name = name
	s.Description = description
	s.Lat = lat
	s.Lng = lng
	s.UpdatedAt = time.Now().UTC()

	data, _ := json.Marshal(s.createEventData())
	event := Event{
		ID:        fmt.Sprintf("%s_%v", s.GetID(), s.GetUpdatedAt().Unix()),
		Type:      StoreUpdated,
		Data:      data,
		Timestamp: s.GetUpdatedAt(),
	}
	s.events = append(s.events, event)

	return nil
}

func (s *Store) PrepareDelete() {
	data, _ := json.Marshal(s.createEventData())

	event := Event{
		ID:        fmt.Sprintf("%s_%v_delete", s.GetID(), s.GetUpdatedAt().Unix()),
		Type:      StoreDeleted,
		Data:      data,
		Timestamp: s.GetUpdatedAt(),
	}

	s.events = append(s.events, event)
}

func (s *Store) createEventData() StoreEventData {
	return StoreEventData{
		ID:          s.GetID(),
		Name:        s.GetName(),
		Description: s.GetDescription(),
		Location:    s.Location(),
		Slug:        s.GetSlug(),
		UserID:      s.GetUserID(),
		CreatedAt:   s.GetCreatedAt(),
		UpdatedAt:   s.GetUpdatedAt(),
	}
}

// Location returns a Location value object
func (s *Store) Location() Location {
	return Location{Lat: s.Lat, Lng: s.Lng}
}

func generateSlug(name string, timestamp int64) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	validChars := ""
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			validChars += string(char)
		}
	}

	return fmt.Sprintf("%s-%d", validChars, timestamp)
}

func (s *Store) TableName() string {
	return "stores"
}

// Location is a value object
type Location struct {
	Lat float64
	Lng float64
}

// Getters
func (s *Store) GetID() string           { return s.ID }
func (s *Store) GetName() string         { return s.Name }
func (s *Store) GetDescription() *string { return s.Description }
func (s *Store) GetLat() float64         { return s.Lat }
func (s *Store) GetLng() float64         { return s.Lng }
func (s *Store) GetSlug() string         { return s.Slug }
func (s *Store) GetUserID() string       { return s.UserID }
func (s *Store) GetCreatedAt() time.Time { return s.CreatedAt }
func (s *Store) GetUpdatedAt() time.Time { return s.UpdatedAt }
