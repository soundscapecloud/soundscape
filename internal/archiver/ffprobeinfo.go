package archiver

type ffprobeInfo struct {
	Format struct {
		BitRate        string `json:"bit_rate"`
		Duration       string `json:"duration"`
		Filename       string `json:"filename"`
		FormatLongName string `json:"format_long_name"`
		FormatName     string `json:"format_name"`
		NbPrograms     int    `json:"nb_programs"`
		NbStreams      int    `json:"nb_streams"`
		ProbeScore     int    `json:"probe_score"`
		Size           string `json:"size"`
		StartTime      string `json:"start_time"`
		Tags           struct {
			CompatibleBrands string `json:"compatible_brands"`
			CreationTime     string `json:"creation_time"`
			MajorBrand       string `json:"major_brand"`
			MinorVersion     string `json:"minor_version"`
		} `json:"tags"`
	} `json:"format"`
	Streams []struct {
		AvgFrameRate       string `json:"avg_frame_rate"`
		BitRate            string `json:"bit_rate"`
		BitsPerRawSample   string `json:"bits_per_raw_sample"`
		ChromaLocation     string `json:"chroma_location"`
		CodecLongName      string `json:"codec_long_name"`
		CodecName          string `json:"codec_name"`
		CodecTag           string `json:"codec_tag"`
		CodecTagString     string `json:"codec_tag_string"`
		CodecTimeBase      string `json:"codec_time_base"`
		CodecType          string `json:"codec_type"`
		CodedHeight        int    `json:"coded_height"`
		CodedWidth         int    `json:"coded_width"`
		ColorPrimaries     string `json:"color_primaries"`
		ColorRange         string `json:"color_range"`
		ColorSpace         string `json:"color_space"`
		ColorTransfer      string `json:"color_transfer"`
		DisplayAspectRatio string `json:"display_aspect_ratio"`
		Disposition        struct {
			AttachedPic     int `json:"attached_pic"`
			CleanEffects    int `json:"clean_effects"`
			Comment         int `json:"comment"`
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			Karaoke         int `json:"karaoke"`
			Lyrics          int `json:"lyrics"`
			Original        int `json:"original"`
			TimedThumbnails int `json:"timed_thumbnails"`
			VisualImpaired  int `json:"visual_impaired"`
		} `json:"disposition"`
		Duration          string `json:"duration"`
		DurationTs        int    `json:"duration_ts"`
		HasBFrames        int    `json:"has_b_frames"`
		Height            int    `json:"height"`
		Index             int    `json:"index"`
		IsAvc             string `json:"is_avc"`
		Level             int    `json:"level"`
		NalLengthSize     string `json:"nal_length_size"`
		NbFrames          string `json:"nb_frames"`
		PixFmt            string `json:"pix_fmt"`
		Profile           string `json:"profile"`
		RFrameRate        string `json:"r_frame_rate"`
		Refs              int    `json:"refs"`
		SampleAspectRatio string `json:"sample_aspect_ratio"`
		StartPts          int    `json:"start_pts"`
		StartTime         string `json:"start_time"`
		Tags              struct {
			CreationTime string `json:"creation_time"`
			HandlerName  string `json:"handler_name"`
			Language     string `json:"language"`
		} `json:"tags"`
		TimeBase string `json:"time_base"`
		Width    int    `json:"width"`
	} `json:"streams"`
}
