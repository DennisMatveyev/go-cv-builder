package media

import "errors"

var (
	ErrFileTooLarge    = errors.New("File too large (max 5MB)")
	ErrInvalidFileType = errors.New("Invalid file type (only JPEG are allowed)")
	ErrSaveFile        = errors.New("Failed to save file")
	ErrRemoveOldFile   = errors.New("Failed to remove old file, new file not saved")
)
