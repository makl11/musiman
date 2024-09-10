package audio

var MUSIC_FILE_TYPES = map[string]bool{
	// https://en.wikipedia.org/wiki/MP3
	"mp3": true,
	// https://en.wikipedia.org/wiki/Ogg
	"ogg": true,
	"oga": true,
	// https://en.wikipedia.org/wiki/Windows_Media_Audio
	"wma": true,
	// https://en.wikipedia.org/wiki/Free_Lossless_Audio_Codec
	"flac": true,
	// https://en.wikipedia.org/wiki/Waveform_Audio_File_Format
	"wav": true,
	// https://en.wikipedia.org/wiki/Audio_Interchange_File_Format
	"aiff": true,
	"aif":  true,
	"aifc": true,
	"snd":  true,
	"iff":  true,
}
