package entity

type Chapter struct {
	Order         int               `bson:"order"`
	Timestamp     string            `json:"timestamp"`
	Duration      string            `json:"duration" `
	Title         map[string]string `json:"title" `
	StotraIds     []string          `json:"stotra_ids"`
	DurationInSec int               `json:"-"`
}

type Variant struct {
	Duration  string    `json:"duration"`
	Chapters  []Chapter `json:"chapters"`
	IsDefault bool      `json:"is_default"`
}

type Prarthana struct {
	TmpId              string              `bson:"TmpId"`
	Id                 string              `bson:"_id"`
	Title              map[string]string   `bson:"title"`
	FestivalIds        []string            `bson:"festival_ids"`
	AudioInfo          AudioInfo           `bson:"audio_info"`
	Days               []int               `bson:"days" `
	Description        map[string]string   `bson:"description" `
	Importance         map[string]string   `bson:"importance"`
	Variants           []Variant           `bson:"variants" `
	Instruction        map[string]string   `bson:"instruction" `
	ItemsRequired      map[string][]string `bson:"items_required" `
	DeityIds           []string            `bson:"deity_ids"`
	UiInfo             PrarthanaUIInfo     `bson:"ui_info"`
	AvailableLanguages []KeyValue          `bson:"available_languages"`
}

type AudioInfo struct {
	IsAudioAvailable bool   `json:"is_audio_available" bson:"is_audio_available"`
	AudioUrl         string `json:"audio_url" bson:"audio_url"`
	IsStudioRecorded bool   `json:"is_studio_recorded" bson:"is_studio_recorded"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PrarthanaUIInfo struct {
	AlbumArt        string `json:"album_art" bson:"album_art"`
	DefaultImageUrl string `json:"default_image_url" bson:"default_image_url"`
	TemplateNumber  string `json:"template_number"`
}
