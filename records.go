package main

import "github.com/pocketbase/pocketbase/tools/types"

type ResultRecord struct {
	Id         string        `db:"id" json:"id"`
	Semester   string        `db:"semester" json:"semester"`
	Data       types.JSONRaw `db:"data" json:"data"`
	Pdf        string        `db:"pdf" json:"pdf"`
	LeadPoints float64       `db:"lead_points" json:"lead_points"`
	Untis      string        `db:"untis" json:"untis"`
}

type PdfsRecord struct {
	Id       string `db:"id" json:"id"`
	Semester string `db:"semester" json:"semester"`
	User     string `db:"user" json:"user"`
	Pdf      string `db:"pdf" json:"pdf"`
}

type UserRecord struct {
	Id    string `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
	Name  string `db:"name" json:"name"`
	Short string `db:"short" json:"short"`
}

type YearsRecordState string

const (
	YearsRecordStateUploaded YearsRecordState = "uploaded"
	YearsRecordStateOpen     YearsRecordState = "open"
	YearsRecordStateClosed   YearsRecordState = "closed"
)

type YearsRecord struct {
	ID          string           `db:"id" json:"id"`
	StartYear   int              `db:"start_year" json:"start_year"`
	Semester    int              `db:"semester" json:"semester"`
	State       YearsRecordState `db:"state" json:"state"`
	Split       types.DateTime   `db:"split" json:"split"`
	BaseMul     float64          `db:"base_mul" json:"base_mul"`         // Base multiplier for the year, used to calculate the points
	LeadPoints  float64          `db:"lead_points" json:"lead_points"`   // Points for each class lead
	TotalPoints float64          `db:"total_points" json:"total_points"` // Total hours for the year
}

type TeacherDataRecord struct {
	ID          string  `db:"id" json:"id"`                 // ID of the record
	Semester    string  `db:"semester" json:"semester"`     // ID of the years record
	UserID      string  `db:"user" json:"user"`             // ID of the user
	AverageTime float64 `db:"avg_time" json:"avg_time"`     // Average time worked by the teacher
	ClassLead   float64 `db:"class_lead" json:"class_lead"` // Class lead percentage
	AddPoints   float64 `db:"add_points" json:"add_points"` // Additional points for the teacher (e.g. for "Mündliches Abitur")
	Short       string  `db:"short" json:"short"`           // Short name of the teacher (Not in the TeacherData Table, needs to be JOINed with users table)
	Name        string  `db:"name" json:"name"`             // Full name of the teacher (Not in the TeacherData Table, needs to be JOINed with users table)
}

type FilesRecordType string

const (
	FilesRecordTypeExam   FilesRecordType = "exam"
	FilesRecordTypeCourse FilesRecordType = "course"
	FilesRecordTypeHours  FilesRecordType = "hours"
)

type FilesRecord struct {
	ID   string          `db:"id" json:"id"`             // ID of the file record
	Year int             `db:"year" json:"year"`         // Year of the record
	Sem  int             `db:"semester" json:"semester"` // Semester of the record
	Type FilesRecordType `db:"type" json:"type"`         // Type of the file (exam, course, hours)
	File string          `db:"file" json:"file"`         // Name of the file
}

type PartPointsRecord struct {
	ID     string  `db:"id" json:"id"`         // ID of the record
	Class  string  `db:"class" json:"class"`   // Class (D, M, E, etc.)
	Grade  string  `db:"grade" json:"grade"`   // Grade (05, 06, Q2, etc.)
	Points float64 `db:"points" json:"points"` // Points earned in the subject
}

type ResultsRecord struct {
	ID       string        `db:"id" json:"id"`             // ID of the results
	Semester string        `db:"semester" json:"semester"` // ID of the years record
	Data     types.JSONRaw `db:"data" json:"data"`         // Data of the results, contains the scores and other information
	Pdf      string        `db:"pdf" json:"pdf"`           // Name of the PDF file containing the results
	Untis    string        `db:"untis" json:"untis"`       // Name of the text file containing the untis data
}
