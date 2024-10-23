package model

import (
	"one-api/common/config"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	TokensPriceType    = "tokens"
	TimesPriceType     = "times"
	DefaultPrice       = 30.0
	DollarRate         = 0.002
	RMBRate            = 0.014
	DefaultCacheRatios = 0.5
	DefaultAudioRatio  = 40
)

type Price struct {
	Model       string  `json:"model" gorm:"type:varchar(100)" binding:"required"`
	Type        string  `json:"type"  gorm:"default:'tokens'" binding:"required"`
	ChannelType int     `json:"channel_type" gorm:"default:0" binding:"gte=0"`
	Input       float64 `json:"input" gorm:"default:0" binding:"gte=0"`
	Output      float64 `json:"output" gorm:"default:0" binding:"gte=0"`

	ExtraRatios map[string]float64 `json:"extra_ratios,omitempty" gorm:"-"`
}

func GetAllPrices() ([]*Price, error) {
	var prices []*Price
	if err := DB.Find(&prices).Error; err != nil {
		return nil, err
	}

	for _, price := range prices {
		price.ExtraRatios = getExtraRatioMap(price.Model)
	}

	return prices, nil
}

func getExtraRatioMap(modelName string) map[string]float64 {
	if !strings.HasPrefix(modelName, "gpt-4o-realtime") && !strings.HasPrefix(modelName, "gpt-4o-audio") {
		return nil
	}

	extraRatios := make(map[string]float64)
	if strings.HasPrefix(modelName, "gpt-4o-realtime") {
		extraRatios["input_audio_tokens_ratio"] = 20
		extraRatios["output_audio_tokens_ratio"] = 10
	} else {
		extraRatios["input_audio_tokens_ratio"] = 40
		extraRatios["output_audio_tokens_ratio"] = 20
	}

	return extraRatios
}

func (price *Price) Update(modelName string) error {
	if err := DB.Model(price).Select("*").Where("model = ?", modelName).Updates(price).Error; err != nil {
		return err
	}

	return nil
}

func (price *Price) Insert() error {
	if err := DB.Create(price).Error; err != nil {
		return err
	}

	return nil
}

func (price *Price) GetInput() float64 {
	if price.Input <= 0 {
		return 0
	}
	return price.Input
}

func (price *Price) GetOutput() float64 {
	if price.Output <= 0 || price.Type == TimesPriceType {
		return 0
	}

	return price.Output
}

func (price *Price) GetExtraRatio(key string) float64 {
	if key == "cached_tokens_ratio" {
		return DefaultCacheRatios
	}

	// 目前只有 音频，如果为空说明有问题，返回最大的一个倍率
	if price.ExtraRatios == nil {
		return DefaultAudioRatio
	}

	if ratio, ok := price.ExtraRatios[key]; ok {
		return ratio
	}

	return DefaultAudioRatio
}

func (price *Price) FetchInputCurrencyPrice(rate float64) string {
	r := decimal.NewFromFloat(price.GetInput()).Mul(decimal.NewFromFloat(rate))
	return r.String()
}

func (price *Price) FetchOutputCurrencyPrice(rate float64) string {
	r := decimal.NewFromFloat(price.GetOutput()).Mul(decimal.NewFromFloat(rate))
	return r.String()
}

func UpdatePrices(tx *gorm.DB, models []string, prices *Price) error {
	err := tx.Model(Price{}).Where("model IN (?)", models).Select("*").Omit("model").Updates(
		Price{
			Type:        prices.Type,
			ChannelType: prices.ChannelType,
			Input:       prices.Input,
			Output:      prices.Output,
		}).Error

	return err
}

func DeletePrices(tx *gorm.DB, models []string) error {
	err := tx.Where("model IN (?)", models).Delete(&Price{}).Error

	return err
}

func InsertPrices(tx *gorm.DB, prices []*Price) error {
	err := tx.CreateInBatches(prices, 100).Error
	return err
}

func DeleteAllPrices(tx *gorm.DB) error {
	err := tx.Where("1=1").Delete(&Price{}).Error
	return err
}

func (price *Price) Delete() error {
	return DB.Where("model = ?", price.Model).Delete(&Price{}).Error
}

type ModelType struct {
	Ratio []float64
	Type  int
}

// 1 === $0.002 / 1K tokens
// 1 === ￥0.014 / 1k tokens
func GetDefaultPrice() []*Price {
	ModelTypes := map[string]ModelType{
		// 	$0.03 / 1K tokens	$0.06 / 1K tokens
		"gpt-4":      {[]float64{15, 30}, config.ChannelTypeOpenAI},
		"gpt-4-0314": {[]float64{15, 30}, config.ChannelTypeOpenAI},
		"gpt-4-0613": {[]float64{15, 30}, config.ChannelTypeOpenAI},
		// 	$0.06 / 1K tokens	$0.12 / 1K tokens
		"gpt-4-32k":      {[]float64{30, 60}, config.ChannelTypeOpenAI},
		"gpt-4-32k-0314": {[]float64{30, 60}, config.ChannelTypeOpenAI},
		"gpt-4-32k-0613": {[]float64{30, 60}, config.ChannelTypeOpenAI},
		// 	$0.01 / 1K tokens	$0.03 / 1K tokens
		"gpt-4-preview":          {[]float64{5, 15}, config.ChannelTypeOpenAI},
		"gpt-4-turbo":            {[]float64{5, 15}, config.ChannelTypeOpenAI},
		"gpt-4-turbo-2024-04-09": {[]float64{5, 15}, config.ChannelTypeOpenAI},
		"gpt-4-1106-preview":     {[]float64{5, 15}, config.ChannelTypeOpenAI},
		"gpt-4-0125-preview":     {[]float64{5, 15}, config.ChannelTypeOpenAI},
		"gpt-4-turbo-preview":    {[]float64{5, 15}, config.ChannelTypeOpenAI},
		"gpt-4-vision-preview":   {[]float64{5, 15}, config.ChannelTypeOpenAI},
		// $0.005 / 1K tokens	$0.015 / 1K tokens
		"gpt-4o": {[]float64{2.5, 7.5}, config.ChannelTypeOpenAI},
		// 	$0.0005 / 1K tokens	$0.0015 / 1K tokens
		"gpt-3.5-turbo":      {[]float64{0.25, 0.75}, config.ChannelTypeOpenAI},
		"gpt-3.5-turbo-0125": {[]float64{0.25, 0.75}, config.ChannelTypeOpenAI},
		// 	$0.0015 / 1K tokens	$0.002 / 1K tokens
		"gpt-3.5-turbo-0301":     {[]float64{0.75, 1}, config.ChannelTypeOpenAI},
		"gpt-3.5-turbo-0613":     {[]float64{0.75, 1}, config.ChannelTypeOpenAI},
		"gpt-3.5-turbo-instruct": {[]float64{0.75, 1}, config.ChannelTypeOpenAI},
		// 	$0.003 / 1K tokens	$0.004 / 1K tokens
		"gpt-3.5-turbo-16k":      {[]float64{1.5, 2}, config.ChannelTypeOpenAI},
		"gpt-3.5-turbo-16k-0613": {[]float64{1.5, 2}, config.ChannelTypeOpenAI},
		// 	$0.001 / 1K tokens	$0.002 / 1K tokens
		"gpt-3.5-turbo-1106": {[]float64{0.5, 1}, config.ChannelTypeOpenAI},
		// 	$0.0020 / 1K tokens
		"davinci-002": {[]float64{1, 1}, config.ChannelTypeOpenAI},
		// 	$0.0004 / 1K tokens
		"babbage-002": {[]float64{0.2, 0.2}, config.ChannelTypeOpenAI},
		// $0.006 / minute -> $0.006 / 150 words -> $0.006 / 200 tokens -> $0.03 / 1k tokens
		"whisper-1": {[]float64{15, 15}, config.ChannelTypeOpenAI},
		// $0.015 / 1K characters
		"tts-1":      {[]float64{7.5, 7.5}, config.ChannelTypeOpenAI},
		"tts-1-1106": {[]float64{7.5, 7.5}, config.ChannelTypeOpenAI},
		// $0.030 / 1K characters
		"tts-1-hd":               {[]float64{15, 15}, config.ChannelTypeOpenAI},
		"tts-1-hd-1106":          {[]float64{15, 15}, config.ChannelTypeOpenAI},
		"text-embedding-ada-002": {[]float64{0.05, 0.05}, config.ChannelTypeOpenAI},
		// 	$0.00002 / 1K tokens
		"text-embedding-3-small": {[]float64{0.01, 0.01}, config.ChannelTypeOpenAI},
		// 	$0.00013 / 1K tokens
		"text-embedding-3-large": {[]float64{0.065, 0.065}, config.ChannelTypeOpenAI},
		"text-moderation-stable": {[]float64{0.1, 0.1}, config.ChannelTypeOpenAI},
		"text-moderation-latest": {[]float64{0.1, 0.1}, config.ChannelTypeOpenAI},
		// $0.016 - $0.020 / image
		"dall-e-2": {[]float64{8, 8}, config.ChannelTypeOpenAI},
		// $0.040 - $0.120 / image
		"dall-e-3": {[]float64{20, 20}, config.ChannelTypeOpenAI},

		// $0.80/million tokens $2.40/million tokens
		"claude-instant-1.2": {[]float64{0.4, 1.2}, config.ChannelTypeAnthropic},
		// $8.00/million tokens $24.00/million tokens
		"claude-2.0": {[]float64{4, 12}, config.ChannelTypeAnthropic},
		"claude-2.1": {[]float64{4, 12}, config.ChannelTypeAnthropic},
		// $15 / M $75 / M
		"claude-3-opus-20240229": {[]float64{7.5, 22.5}, config.ChannelTypeAnthropic},
		//  $3 / M $15 / M
		"claude-3-sonnet-20240229": {[]float64{1.3, 3.9}, config.ChannelTypeAnthropic},
		//  $0.25 / M $1.25 / M  0.00025$ / 1k tokens 0.00125$ / 1k tokens
		"claude-3-haiku-20240307": {[]float64{0.125, 0.625}, config.ChannelTypeAnthropic},

		// 0.0005$ / 1k tokens 0.0015$ / 1k tokens
		"gemini-pro":        {[]float64{0.25, 0.75}, config.ChannelTypeGemini},
		"gemini-pro-vision": {[]float64{0.25, 0.75}, config.ChannelTypeGemini},
		"gemini-1.0-pro":    {[]float64{0.25, 0.75}, config.ChannelTypeGemini},
		// $7 / 1 million tokens  $21 / 1 million tokens
		"gemini-1.5-pro":          {[]float64{1.75, 5.25}, config.ChannelTypeGemini},
		"gemini-1.5-pro-latest":   {[]float64{1.75, 5.25}, config.ChannelTypeGemini},
		"gemini-1.5-flash":        {[]float64{0.175, 0.265}, config.ChannelTypeGemini},
		"gemini-1.5-flash-latest": {[]float64{0.175, 0.265}, config.ChannelTypeGemini},
		"gemini-ultra":            {[]float64{1, 1}, config.ChannelTypeGemini},

		"open-mistral-7b":       {[]float64{0.125, 0.125}, config.ChannelTypeMistral}, // 0.25$ / 1M tokens	0.25$ / 1M tokens  0.00025$ / 1k tokens
		"open-mixtral-8x7b":     {[]float64{0.35, 0.35}, config.ChannelTypeMistral},   // 0.7$ / 1M tokens	0.7$ / 1M tokens  0.0007$ / 1k tokens
		"mistral-small-latest":  {[]float64{1, 3}, config.ChannelTypeMistral},         // 2$ / 1M tokens	6$ / 1M tokens  0.002$ / 1k tokens
		"mistral-medium-latest": {[]float64{1.35, 4.05}, config.ChannelTypeMistral},   // 2.7$ / 1M tokens	8.1$ / 1M tokens  0.0027$ / 1k tokens
		"mistral-large-latest":  {[]float64{4, 12}, config.ChannelTypeMistral},        // 8$ / 1M tokens	24$ / 1M tokens  0.008$ / 1k tokens
		"mistral-embed":         {[]float64{0.05, 0.05}, config.ChannelTypeMistral},   // 0.1$ / 1M tokens 0.1$ / 1M tokens  0.0001$ / 1k tokens
	}

	var prices []*Price

	for model, modelType := range ModelTypes {
		prices = append(prices, &Price{
			Model:       model,
			Type:        TokensPriceType,
			ChannelType: modelType.Type,
			Input:       modelType.Ratio[0],
			Output:      modelType.Ratio[1],
		})
	}

	return prices
}
