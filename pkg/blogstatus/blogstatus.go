package blogstatus

type BlogStatus string

const (
	Draft     BlogStatus = "draft"
	Published BlogStatus = "published"
	Archived  BlogStatus = "archived"
)

func (s BlogStatus) IsValid() bool {
	switch s {
	case Draft, Published, Archived:
		return true
	}
	return false
}
