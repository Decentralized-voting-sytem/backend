package models

import (
    "time"
    "gorm.io/gorm"
)

type Voter struct {
    ID        uint           `gorm:"primaryKey"`
    VoterID   string         `gorm:"unique;not null"` 
    Name      string         `gorm:"not null"`
    DOB       time.Time      `gorm:"not null;type:date"`        
    Password  string         `gorm:"not null"`       
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Candidate struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Vote struct {
    ID          uint           `gorm:"primaryKey"`
    VoterID     string           `gorm:"not null"`       
    CandidateID uint           `gorm:"not null"`     
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Admin struct {
    ID        uint           `gorm:"primaryKey"`
    AdminID   string         `gorm:"unique;not null"` 
    Name      string         `gorm:"not null"`
    Password  string         `gorm:"not null"`        
    CreatedAt time.Time
    UpdatedAt time.Time
}