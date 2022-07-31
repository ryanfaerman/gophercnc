package tool

type Tool struct {
	Material     Material `json:"BMC"`
	Description  string   `json:"description"`
	GUID         string   `json:"guid"`
	LastModified int      `json:"last_modified"`
	PostProcess  struct {
		BreakControl     bool   `json:"break-control"`
		Comment          string `json:"comment"`
		DiameterOffset   int    `json:"diameter-offset"`
		LengthOffset     int    `json:"length-offset"`
		Live             bool   `json:"live"`
		ManualToolChange bool   `json:"manual-tool-change"`
		Number           int    `json:"number"`
		Turret           int    `json:"turret"`
	} `json:"post-process"`
	ProductID     string `json:"product-id"`
	ProductLink   string `json:"product-link"`
	ReferenceGUID string `json:"reference_guid"`
	StartValues   struct {
		Presets []CuttingData `json:"presets"`
	} `json:"start-values"`
	Vendor   string  `json:"vendor"`
	Type     Endmill `json:"type"`
	Unit     string  `json:"unit"`
	Geometry struct {
		Diameter                 float64 `json:"DC"`
		ShaftDiameter            float64 `json:"SFDM"`
		OverallLength            float64 `json:"OAL"`
		LengthBelowHolder        float64 `json:"LB"`
		ShoulderLength           float64 `json:"shoulder-length"`
		FluteLength              float64 `json:"LCF"`
		FluteCount               int     `json:"NOF"`
		CounterClockwiseRotation bool    `json:"HAND"`
	} `json:"geometry"`
	Shaft struct {
		Segments []struct {
			Height        float64 `json:"height"`
			UpperDiamater float64 `json:"upper-diamater"`
			LowerDiameter float64 `json:"lower-diameter"`
		} `json:"segments"`
		Type string `json:"type"`
	} `json:"shaft"`
}

func (t Tool) Number() int { return t.PostProcess.Number }

type CuttingData struct {
	GUID              string
	Name              string
	FeedPerRevolution float64 `json:"f_n"`
	FeedPerTooth      float64 `json:"f_z"`
	SpindleSpeed      float64 `json:"n"`
	RampSpindleSpeed  float64 `json:"n_ramp"`
	ToolCoolant       string  `json:"tool-coolant"`
	UseStepdown       bool    `json:"use-stepdown"`
	UseStepover       bool    `json:"use-stepover"`
	SurfaceSpeed      float64 `json:"v_c"`
	CuttingFeedrate   float64 `json:"v_f"`
	LeadInFeedrate    float64 `json:"v_f_leadIn"`
	LeadOutFeedrate   float64 `json:"v_f_leadOut"`
	PlungeFeedrate    float64 `json:"v_f_plunge"`
	RampFeedrate      float64 `json:"v_f_ramp"`
}
